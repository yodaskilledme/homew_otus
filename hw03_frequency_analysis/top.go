package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

type WordCount struct {
	word      string
	numberCnt int
}

type PairList []WordCount

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].numberCnt < p[j].numberCnt
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
// just to fix travis ci cache...
func Top10(str string) []string {
	if str == "" {
		return []string{}
	}

	wordsCount := countWords(str)
	sortedWords := sortWords(wordsCount)

	res := make([]string, 0, 10)
	for key, value := range sortedWords {
		if key == 10 {
			break
		}
		res = append(res, value.word)
	}

	return res
}

func countWords(str string) map[string]int {
	counter := make(map[string]int)
	wordSlice := strings.Fields(str)
	for key := range wordSlice {
		counter[wordSlice[key]]++
	}

	return counter
}

func sortWords(stings map[string]int) PairList {
	sortedWords := make(PairList, 0, len(stings))

	for word, numberCnt := range stings {
		sortedWords = append(sortedWords, WordCount{
			word:      word,
			numberCnt: numberCnt,
		})
	}

	sort.Sort(sort.Reverse(sortedWords))

	return sortedWords
}
