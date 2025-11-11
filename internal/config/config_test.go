package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		envs     map[string]string
		expected Config
		panic    assert.PanicAssertionFunc
	}{
		{
			desc: "sucess",
			envs: map[string]string{
				"SERVER_ADDRESS": "localhost",
				"SERVER_PORT":    "8000",
			},
			expected: Config{
				ServerAddress: "localhost",
				ServerPort:    8000,
			},
			panic: assert.NotPanics,
		},
		{
			desc: "failure with critical invalid configurations",
			envs: map[string]string{
				"SERVER_ADDRESS": "localhost",
				"SERVER_PORT":    "invalid",
			},
			expected: Config{},
			panic:    assert.Panics,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			for env, value := range tC.envs {
				old := os.Getenv(env)
				os.Setenv(env, value)
				t.Cleanup(func() {
					os.Setenv(env, old)
				})
			}

			tC.panic(t, func() {
				res := InitConfig()
				assert.Equal(t, tC.expected, res)
			})
		})
	}
}
