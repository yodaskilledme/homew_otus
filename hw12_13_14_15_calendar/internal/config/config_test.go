package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	t.Run("read config", func(t *testing.T) {
		got, err := ParseConfig("../../configs/config.yml")

		require.NoError(t, err)

		require.NotEmpty(t, got.Http)
		require.NotEmpty(t, got.Http.Host)
		require.NotEmpty(t, got.Http.Port)

		require.NotEmpty(t, got.Database)
		require.NotEmpty(t, got.Database.Host)
		require.NotEmpty(t, got.Database.Port)
		require.NotEmpty(t, got.Database.User)
		require.NotEmpty(t, got.Database.Password)
		require.NotEmpty(t, got.Database.Name)

		require.NotEmpty(t, got.Logger)
		require.NotEmpty(t, got.Logger.Level)
		require.NotEmpty(t, got.Logger.Output)

		require.NotEmpty(t, got.Grpc)
		require.NotEmpty(t, got.Grpc.Host)
		require.NotEmpty(t, got.Grpc.Port)

		require.NotEmpty(t, got.Amqp)
		require.NotEmpty(t, got.Amqp.Uri)
		require.NotEmpty(t, got.Amqp.Queue)
	})

	t.Run("read config, wrong path", func(t *testing.T) {
		_, err := ParseConfig("../../configs/configggg.yml")

		require.Error(t, err)
	})
}