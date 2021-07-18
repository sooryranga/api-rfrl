package rfrl

import "os"

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

func GetStatusInternalServerError(err error) string {
	backendType := GetBackendType()

	switch backendType {
	case production:
		return productionInternalError
	default:
		return err.Error()
	}
}
