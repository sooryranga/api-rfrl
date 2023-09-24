package document

import (
	"database/sql"
	"time"
)

// Document model
type Document struct {
	ID          int            `db:"id" json:"id" mapstructure:"id"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at" mapstructure:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at" mapstructure:"updated_at"`
	Src         string         `db:"src" json:"src" mapstructure:"src"`
	Name        string         `db:"name" json:"name" mapstructure:"name"`
	Description sql.NullString `db:"description" json:"description" mapstructure:"description"`
	ClientID    string         `db:"client_id"`
}

type DocumentOrder struct {
	ID         int       `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	RefType    string    `db:"ref_type" json:"ref_type"`
	RefID      string    `db:"ref_id" json:"ref_id"`
	DocumentID int       `db:"document_id" json:"document_id"`
	Page       int       `db:"page" json"page"`
}

// NewDocument creates new client model struct
func NewDocument(
	clientID string,
	src string,
	name string,
	description string,
) *Document {
	document := Document{Name: name, Src: src, ClientID: clientID}
	if description != "" {
		document.Description = sql.NullString{String: description, Valid: true}
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
