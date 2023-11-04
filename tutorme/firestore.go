package tutorme

import "gopkg.in/guregu/null.v4"

type FirstStoreClient interface {
	CreateClient(id string, photo string, firstName string, lastName string) error
	UpdateClient(id string, photo null.String, firstName null.String, lastName null.String) error
}
