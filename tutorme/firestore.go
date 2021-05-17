package tutorme

import "gopkg.in/guregu/null.v4"

type FireStoreClient interface {
	CreateClient(id string, photo string, firstName string, lastName string) error
	UpdateClient(id string, photo null.String, firstName null.String, lastName null.String) error
	UpdateCode(sessionID int, codeID int, result string) error
	CreateCode(sessionID int, codeID int) error
}
