package store

import (
	"sort"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
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
WHERE id in ($2) AND client_id in ($1) 
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

	getDocumentOrderQuery string = `
SELCT document.* FROM document_order
JOIN document ON document.id = document_order.document_id
WHERE document_order.ref_type = $1 AND document_order.ref_id = $2
ORDER BY document_order.page
`
)

// GetDocument queries the database for document with id and client id
func (dc *DocumentStore) GetDocument(db tutorme.DB, id int, clientID string) (*tutorme.Document, error) {
	var m tutorme.Document
	err := db.QueryRowx(getDocumentByID, id, clientID).StructScan(&m)

	return &m, errors.Wrap(err, "getDocument")
}

// CreateDocument creates a new row for a document in the database
func (dc *DocumentStore) CreateDocument(db tutorme.DB, doc *tutorme.Document) (*tutorme.Document, error) {
	row := db.QueryRowx(
		insertDocument,
		doc.Src,
		doc.Name,
		doc.Description,
		doc.ClientID,
	)

	var m tutorme.Document

	err := row.StructScan(&m)
	return &m, errors.Wrap(err, "CreateDocument")
}

// UpdateDocument updates a client in the database
func (dc *DocumentStore) UpdateDocument(db tutorme.DB, ID int, doc *tutorme.Document) (*tutorme.Document, error) {
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

	var m tutorme.Document

	err = row.StructScan(&m)
	return &m, errors.Wrap(err, "UpdateDocument")
}

// DeleteDocument deletes a document in the database
func (dc *DocumentStore) DeleteDocument(db tutorme.DB, ID int, clientID string) error {
	_, err := db.Queryx(deleteDocument, ID, clientID)

	return errors.Wrap(err, "DeleteDocument")
}

func (dc *DocumentStore) CheckDocumentsBelongToclients(
	db tutorme.DB,
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

	return count != len(documentIds), err
}

func (dc *DocumentStore) RemoveAndRenumberDocumentsOrder(db tutorme.DB, ID int, clientID string) error {

	rows, err := db.Queryx(removeDocumentOrder, ID)
	var docOrders []tutorme.DocumentOrder

	for rows.Next() {
		var docOrder tutorme.DocumentOrder
		err = rows.StructScan(&docOrder)
		if err != nil {
			return err
		}
		docOrders = append(docOrders, docOrder)
	}

	query := sq.Update("document_order").
		Set("page", sq.Expr("page - 1"))

	for i := 0; i < len(docOrders); i += 1 {
		docOrder := docOrders[i]
		query = query.Where(
			sq.Or{
				sq.And{
					sq.Eq{"ref_type": docOrder.RefType},
					sq.Eq{"ref_if": docOrder.RefID},
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

	return errors.Wrap(err, "RemoveAndRenumberDocuments")
}

func createDocumentOrder(
	db tutorme.DB,
	documentIds []int,
	refType string,
	refID string,
) ([]tutorme.Document, error) {
	query := sq.Insert("document_order").
		Columns("ref_type", "ref_id", "document_id", "page")

	for i := 0; i < len(documentIds); i += 1 {
		query.Values(refType, refID, documentIds[i], i)
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
	db tutorme.DB,
	documentIds []int,
	refType string,
	refID string,
) ([]tutorme.Document, error) {
	return createDocumentOrder(db, documentIds, refType, refID)
}

func (dc *DocumentStore) UpdateDocumentOrder(
	db tutorme.DB,
	documentIds []int,
	refType string,
	refID string,
) ([]tutorme.Document, error) {
	sort.Ints(documentIds)

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
	db tutorme.DB,
	refType string,
	refID string,
) ([]tutorme.Document, error) {
	rows, err := db.Queryx(getDocumentOrderQuery, refType, refID)

	if err != nil {
		return nil, err
	}

	var docs []tutorme.Document
	for rows.Next() {
		var doc tutorme.Document
		err = rows.StructScan(&doc)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, err
}

func (dc *DocumentStore) GetDocumentOrder(
	db tutorme.DB,
	refType string,
	refID string,
) ([]tutorme.Document, error) {
	return getDocumentOrder(db, refType, refID)
}
