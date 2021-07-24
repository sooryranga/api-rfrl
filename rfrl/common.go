package rfrl

import (
	"os"

	"github.com/pkg/errors"
)

const (
	production = "production"
	local      = "local"
	stage      = "stage"
)

const (
	productionInternalError = "Something went wrong!"
)

func GetBackendType() string {
	version, isSet := os.LookupEnv("BACKEND_TYPE")

	if !isSet {
		version = production
	}

	return version
}

func IsProduction() bool {
	return GetBackendType() == production
}

func GetStatusInternalServerError(err error) error {
	backendType := GetBackendType()

	switch backendType {
	case production:
		return errors.New(productionInternalError)
	default:
		return err
	}
}
