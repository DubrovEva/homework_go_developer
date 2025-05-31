package config

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	config := &Config{
		Service: Service{
			Host: "localhost",
			Port: "8080",
		},
		Products: Products{
			Host:       "product-service",
			Port:       "8082",
			Token:      "testToken",
			Schema:     "http",
			RetryCount: 3,
		},
	}

	data, err := yaml.Marshal(config)
	require.NoError(t, err)

	// convert data to io.Reader
	result, err := readConfig(bytes.NewReader(data))

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, config, result)
}
