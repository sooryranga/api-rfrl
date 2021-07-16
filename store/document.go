package store

import (
	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// DocumentStore holds all store related functions for document
type DocumentStore struct{}

// NewDocumentStore creates new DocumentStore
func NewDocumentStore() *DocumentStore {
	return &DocumentStore{}
}

const (
	getDocumentByID string = `
SELECT * FROM document
WHERE id = $1 AND client_id = $2
	`
	insertDocument string = `
INSERT INTO document (src, name, description, client_id)
VALUES ($1, $2, $3, $4)
RETURNING *
	`
	checkDocumentsBelongToclients string = `
SELECT COUNT(*) FROM document 
WHERE id in (?) AND client_id in (?) 
`
	deleteDocument string = `
DELETE FROM document WHERE id = $1 AND client_id = $2
`
	removeDocumentFromDocumentOrder string = `
DELETE FROM document_order 
WHERE document_id = $1
RETURNING ref_type, ref_type, page
`

	removeDocumentOrder string = `
DELETE FROM document_order 
WHERE ref_type = $1 AND ref_id = $2
`

	deleteDocumentOrderWithDocumentID string = `
DELETE FROM document_order 
WHERE document_id = $1 
RETURNING *
`

	getDocumentOrderQuery string = `
SELECT document.* FROM document_order
JOIN document ON document.id = document_order.document_id
WHERE document_order.ref_type = $1 AND document_order.ref_id = $2
ORDER BY document_order.page
`
)

// GetDocument queries the database for document with id and client id
func (dc *DocumentStore) GetDocument(db rfrl.DB, id int, clientID string) (*rfrl.Document, error) {
	var m rfrl.Document
	err := db.QueryRowx(getDocumentByID, id, clientID).StructScan(&m)

	return &m, errors.Wrap(err, "getDocument")
}

// CreateDocument creates a new row for a document in the database
func (dc *DocumentStore) CreateDocument(db rfrl.DB, doc *rfrl.Document) (*rfrl.Document, error) {
	row := db.QueryRowx(
		insertDocument,
		doc.Src,
		doc.Name,
		doc.Description,
		doc.ClientID,
	)

	var m rfrl.Document

	err := row.StructScan(&m)
	return &m, errors.Wrap(err, "CreateDocument")
}

// UpdateDocument updates a client in the database
func (dc *DocumentStore) UpdateDocument(db rfrl.DB, ID int, doc *rfrl.Document) (*rfrl.Document, error) {
	query := sq.Update("document")
	if doc.Description.Valid {
		query = query.Set("description", doc.Description)
	}

	sql, args, err := query.
		Set("src", doc.Src).
		Set("name", doc.Name).
		Where(sq.Eq{"id": ID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(
		sql,
		args...,
	)

	var m rfrl.Document

	err = row.StructScan(&m)
	return &m, errors.Wrap(err, "UpdateDocument")
}

// DeleteDocument deletes a document in the database
func (dc *DocumentStore) DeleteDocument(db rfrl.DB, ID int, clientID string) error {
	_, err := db.Queryx(deleteDocument, ID, clientID)

	return errors.Wrap(err, "DeleteDocument")
}

func (dc *DocumentStore) CheckDocumentsBelongToclients(
	db rfrl.DB,
	clientIDs []string,
	documentIds []int,
) (bool, error) {
	query, args, err := sqlx.In(checkDocumentsBelongToclients, documentIds, clientIDs)

	if err != nil {
		return false, err
	}
	query = db.Rebind(query)
	row := db.QueryRowx(query, args...)

	var count int
	err = row.Scan(&count)

	return count == len(documentIds), err
}

func (dc *DocumentStore) RemoveAndRenumberDocumentsOrder(db rfrl.DB, ID int, clientID string) error {
	rows, err := db.Queryx(deleteDocumentOrderWithDocumentID, ID)

	if err != nil {
		return err
	}

	docOrders := make([]rfrl.DocumentOrder, 0)

	for rows.Next() {
		var docOrder rfrl.DocumentOrder
		err = rows.StructScan(&docOrder)
		if err != nil {
			return err
		}
		docOrders = append(docOrders, docOrder)
	}

	query := sq.Update("document_order").
		Set("page", sq.Expr("page - 1"))

	for i := 0; i < len(docOrders); i++ {
		docOrder := docOrders[i]
		query = query.Where(
			sq.Or{
				sq.And{
					sq.Eq{"ref_type": docOrder.RefType},
					sq.Eq{"ref_id": docOrder.RefID},
					sq.Expr("page > ?", docOrder.Page),
				},
			},
		)
	}

	sql, args, err := query.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Queryx(
		sql,
		args...,
	)

	return err
}

func createDocumentOrder(
	db rfrl.DB,
	documentIds []int,
	refType string,
	refID string,
) ([]rfrl.Document, error) {
	query := sq.Insert("document_order").
		Columns("ref_type", "ref_id", "document_id", "page")

	for i := 0; i < len(documentIds); i++ {
		query = query.Values(refType, refID, documentIds[i], i)
	}

	sql, args, err := query.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = db.Queryx(
		sql,
		args...,
	)

	return getDocumentOrder(db, refType, refID)
}

func (dc *DocumentStore) CreateDocumentOrder(
	db rfrl.DB,
	documentIds []int,
	refType string,
	refID string,
) ([]rfrl.Document, error) {
	return createDocumentOrder(db, documentIds, refType, refID)
}

func (dc *DocumentStore) UpdateDocumentOrder(
	db rfrl.DB,
	documentIds []int,
	refType string,
	refID string,
) ([]rfrl.Document, error) {
	_, err := db.Queryx(removeDocumentOrder, refType, refID)
	if err != nil {
		return nil, err
	}

	docs, err := createDocumentOrder(
		db,
		documentIds,
		refType,
		refID,
	)

	return docs, err
}

func getDocumentOrder(
	db rfrl.DB,
	refType string,
	refID string,
) ([]rfrl.Document, error) {
	rows, err := db.Queryx(getDocumentOrderQuery, refType, refID)

	if err != nil {
		return nil, err
	}

	docs := make([]rfrl.Document, 0)

	for rows.Next() {
		var doc rfrl.Document
		err = rows.StructScan(&doc)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, err
}

func (dc *DocumentStore) GetDocumentOrder(
	db rfrl.DB,
	refType string,
	refID string,
) ([]rfrl.Document, error) {
	return getDocumentOrder(db, refType, refID)
}
