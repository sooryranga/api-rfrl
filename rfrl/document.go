package rfrl

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// Document model
type Document struct {
	ID          int         `db:"id" json:"id"`
	CreatedAt   time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updatedAt"`
	Src         string      `db:"src" json:"src"`
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	ClientID    string      `db:"client_id"`
}

type DocumentOrder struct {
	ID         int       `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	RefType    string    `db:"ref_type" json:"refType"`
	RefID      string    `db:"ref_id" json:"refId"`
	DocumentID int       `db:"document_id" json:"documentId"`
	Page       int       `db:"page" json"page"`
}

// NewDocument creates new client model struct
func NewDocument(
	clientID string,
	src string,
	name string,
	description string,
) *Document {
	document := Document{
		Name:        name,
		Src:         src,
		ClientID:    clientID,
		Description: null.NewString(description, description != ""),
	}
	return &document
}

func NewDocumentOrder(
	refType string,
	refID string,
	documentID int,
	page int,
) *DocumentOrder {
	return &DocumentOrder{
		RefType:    refType,
		RefID:      refID,
		DocumentID: documentID,
		Page:       page,
	}
}

const (
	ClientRef  string = "client"
	SessionRef string = "session"
)

type DocumentStore interface {
	GetDocument(db DB, id int, clientID string) (*Document, error)
	CreateDocument(db DB, doc *Document) (*Document, error)
	UpdateDocument(db DB, ID int, doc *Document) (*Document, error)
	DeleteDocument(db DB, ID int, clientID string) error
	CheckDocumentsBelongToclients(db DB, clientIDs []string, documentIds []int) (bool, error)
	RemoveAndRenumberDocumentsOrder(db DB, ID int, clientID string) error
	CreateDocumentOrder(db DB, documentIds []int, refType string, refID string) ([]Document, error)
	UpdateDocumentOrder(db DB, documentIds []int, refType string, refID string) ([]Document, error)
	GetDocumentOrder(db DB, refType string, refID string) ([]Document, error)
}

type DocumentUseCase interface {
	CreateDocument(clientID string, src string, name string, description string) (*Document, error)
	UpdateDocument(clientID string, ID int, src string, name string, description string) (*Document, error)
	DeleteDocument(clientID string, ID int) error
	GetDocument(id int, clientID string) (*Document, error)
	CreateDocumentOrder(clientID string, documentIds []int, refID string, refType string) ([]Document, error)
	UpdateDocumentOrder(clientID string, documentIds []int, refID string, refType string) ([]Document, error)
	GetDocumentOrder(clientID string, refID string, refType string) ([]Document, error)
}
