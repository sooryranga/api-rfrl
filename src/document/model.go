package client

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
	UserId      string         `db:"user_id"`
}

type DocumentOrder struct {
	ID         int    `db:"id" json:"id"`
	RefType    string `db:"ref_type" json:"ref_type"`
	RefID      string `db:"ref_id" json:"ref_id"`
	DocumentID int    `db:"document_id" json:"document_id"`
	Order      int    `db:"order" json"order"`
}

// NewDocument creates new client model struct
func NewDocument(
	src string,
	name string,
	description string,
	userId string,
) *Document {
	document := Document{Name: name, Src: src, UserId: userId}
	if description != "" {
		document.Description = sql.NullString{String: description, Valid: true}
	}
	return &document
}

func NewDocumentOrder(
	refType string,
	refID string,
	documentID int,
	order int,
) *DocumentOrder {
	return &DocumentOrder{
		RefType:    refType,
		RefID:      refID,
		DocumentID: documentID,
		Order:      order,
	}
}
