package functions

import (
	"fmt"
	"runtime"
	"strings"
)

type nameParams struct {
	skip int
}

type Option func(*nameParams)

func GetPackageAndFunctionName(opts ...Option) string {
	params := &nameParams{}

	for _, opt := range opts {
		opt(params)
	}

	pc, _, _, ok := runtime.Caller(params.skip + 1)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		methodPathChunks := strings.Split(details.Name(), "/")
		splitted := strings.Split(methodPathChunks[len(methodPathChunks)-1], ".")

		return fmt.Sprintf("%s.%s", splitted[0], splitted[len(splitted)-1])
	}

	return ""
}

func GetFunctionName(opts ...Option) string {
	params := &nameParams{}

	for _, opt := range opts {
		opt(params)
	}

	delimiter := "."
	packageAndFunctionName := GetPackageAndFunctionName(WithSkip(params.skip + 1))

	splitted := strings.Split(packageAndFunctionName, delimiter)
	functionName := splitted[len(splitted)-1]

	return functionName
}

func GetPackageName(opts ...Option) string {
	params := &nameParams{}

	for _, opt := range opts {
		opt(params)
	}

	delimiter := "."
	packageAndFunctionName := GetPackageAndFunctionName(WithSkip(params.skip + 1))

	splitted := strings.Split(packageAndFunctionName, delimiter)
	packageName := splitted[0]

	return packageName
}

func WithSkip(skip int) Option {
	return func(params *nameParams) {
		params.skip = skip
	}
}
