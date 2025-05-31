package tests

import (
	"os"
)

const (
	CartServiceAddress = "http://localhost:8080"

	DefaultUserID = int64(31337)
	DefaultSKU    = int64(1076963)
)

func IsContainerRun() bool {
	_, err := os.Stat("container.lock")

	return err == nil
}
