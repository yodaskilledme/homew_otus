package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	outFileName := "out.txt"
	defer os.Remove(outFileName)

	t.Run("src file is not regular", func(t *testing.T) {
		err := Copy("/dev/urandom", outFileName, 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("src file not exists err", func(t *testing.T) {
		err := Copy("non_existent_file.txt", outFileName, 0, 0)
		require.Error(t, err)
	})

	t.Run("offset exceeds src file size", func(t *testing.T) {
		tmpFile, err := ioutil.TempFile("", "copy_test")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("test")
		if err != nil {
			log.Fatal(err)
		}

		err = Copy(tmpFile.Name(), outFileName, 5, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	srcFileData := []byte{'1', '2', '3'}
	successTestCases := []struct {
		limit               int64
		offset              int64
		srcFileData         []byte
		expectedOutFileData []byte
	}{
		{0, 0, srcFileData, []byte{'1', '2', '3'}},
		{1000, 0, srcFileData, []byte{'1', '2', '3'}},
		{1, 0, srcFileData, []byte{'1'}},
		{2, 0, srcFileData, []byte{'1', '2'}},
		{0, 1, srcFileData, []byte{'2', '3'}},
		{0, 2, srcFileData, []byte{'3'}},
		{0, 3, srcFileData, []byte{}},
		{1, 1, srcFileData, []byte{'2'}},
		{2, 2, srcFileData, []byte{'3'}},
	}

	for _, tt := range successTestCases {
		t.Run(fmt.Sprintf("success Copy with limit=%d,offset=%d", tt.limit, tt.offset), func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "copy_test")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.Write(tt.srcFileData)
			if err != nil {
				log.Fatal(err)
			}

			err = Copy(tmpFile.Name(), outFileName, tt.offset, tt.limit)
			require.NoError(t, err)

			outFile, err := os.Open(outFileName)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err = outFile.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			outFileData, err := ioutil.ReadFile(outFileName)
			if err != nil {
				log.Fatal(err)
			}

			require.Equal(t, tt.expectedOutFileData, outFileData)
		})
	}
}
