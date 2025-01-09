package must

import (
	"errors"
	"log"
)

func Must[T any](value T, err error) T {
	if err != nil {
		log.Panic(err)
	}

	return value
}

func MustServe(err error) {
	_ = Must(err, err)
}

func MustEnv(value string) string {
	if value == "" {
		log.Panic(ErrEnvNotFound.Error())
	}

	return value
}

var ErrEnvNotFound = errors.New("environment variable not found")
