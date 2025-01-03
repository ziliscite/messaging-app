package must

import "errors"

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

func MustServe(err error) {
	_ = Must(err, err)
}

func MustEnv(value string) string {
	if value == "" {
		panic(ErrEnvNotFound.Error())
	}

	return value
}

var ErrEnvNotFound = errors.New("environment variable not found")
