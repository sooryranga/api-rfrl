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

func getBackendType() string {
	version, isSet := os.LookupEnv("BACKEND_TYPE")

	if !isSet {
		version = production
	}

	return version
}

func GetStatusInternalServerError(err error) string {
	backendType := getBackendType()

	switch backendType {
	case production:
		return productionInternalError
	default:
		return err.Error()
	}
}
