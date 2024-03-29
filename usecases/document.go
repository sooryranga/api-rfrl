package usecases

import (
	"fmt"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// DocumentUseCase holds all business related functions for document
type DocumentUseCase struct {
	db            *sqlx.DB
	documentStore rfrl.DocumentStore
}

// NewDocumentUseCase creates new ClientUseCase
func NewDocumentUseCase(db sqlx.DB, documentStore rfrl.DocumentStore) *DocumentUseCase {
	return &DocumentUseCase{&db, documentStore}
}

// CreateDocument creates a new document
func (dc *DocumentUseCase) CreateDocument(
	clientID string,
	src string,
	name string,
	description string,
) (*rfrl.Document, error) {
	document := rfrl.NewDocument(clientID, src, name, description)
	return dc.documentStore.CreateDocument(dc.db, document)
}

// UpdateDocument updates an existing document
func (dc *DocumentUseCase) UpdateDocument(
	clientID string,
	ID int,
	src string,
	name string,
	description string,
) (*rfrl.Document, error) {
	document := rfrl.NewDocument(clientID, src, name, description)
	return dc.documentStore.UpdateDocument(dc.db, ID, document)
}

// DeleteDocument deletes existing document
func (dc *DocumentUseCase) DeleteDocument(
	clientID string,
	ID int,
) error {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = dc.db.Beginx()

	if *err != nil {
		return errors.Wrap(*err, "DeleteDocument")
	}

	defer rfrl.HandleTransactions(tx, err)

	*err = dc.documentStore.RemoveAndRenumberDocumentsOrder(tx, ID, clientID)

	if *err != nil {
		return *err
	}

	*err = dc.documentStore.DeleteDocument(tx, ID, clientID)

	if *err != nil {
		return *err
	}

	return nil
}

// GetDocument gets an existing document
func (dc *DocumentUseCase) GetDocument(
	id int,
	clientID string,
) (*rfrl.Document, error) {
	return dc.documentStore.GetDocument(dc.db, id, clientID)
}

// CreateDocumentOrder creates a new order of documents
func (dc *DocumentUseCase) CreateDocumentOrder(
	clientID string,
	documentIds []int,
	refID string,
	refType string,
) ([]rfrl.Document, error) {
	check, err := dc.checkclientIsInRef(clientID, refType, refID)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf("Unauthorized to create document-order")
	}

	check, err = dc.checkDocumentsAreForRef(documentIds, refType, refID)

	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf(
			"Documents are not related to client: Document - %v",
			documentIds,
		)
	}

	return dc.documentStore.CreateDocumentOrder(dc.db, documentIds, refType, refID)
}

// UpdateDocumentOrder updates existing document order (ie reshuffling)
func (dc *DocumentUseCase) UpdateDocumentOrder(
	clientID string,
	documentIDs []int,
	refID string,
	refType string,
) ([]rfrl.Document, error) {
	check, err := dc.checkclientIsInRef(clientID, refType, refID)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf("Unauthorized to update document-order")
	}

	check, err = dc.checkDocumentsAreForRef(documentIDs, refType, refID)

	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf(
			"Documents are not related to client: Document - %v",
			documentIDs,
		)
	}

	return dc.documentStore.UpdateDocumentOrder(dc.db, documentIDs, refType, refID)
}

// GetDocumentOrder grabs document in order
func (dc *DocumentUseCase) GetDocumentOrder(
	clientID string,
	refID string,
	refType string,
) ([]rfrl.Document, error) {
	check, err := dc.checkclientIsAbleToViewDocument(clientID, refType, refID)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.New(
			fmt.Sprintf("Client is not part of the ref_type (%s, %s)", refType, refID),
		)
	}

	return dc.documentStore.GetDocumentOrder(dc.db, refType, refID)
}

func (dc *DocumentUseCase) checkDocumentsAreForRef(
	documentIds []int,
	refType string,
	refID string,
) (bool, error) {
	var clientIDs []string
	if refType == rfrl.ClientRef {
		clientIDs = []string{refID}
	} else if refType == rfrl.SessionRef {
		// clientIDs = session.GetClientsIdForSession(h.db, refID)
		return false, nil
	} else {
		return false, errors.Errorf("reftype (%v) not found", refType)
	}

	return dc.documentStore.CheckDocumentsBelongToclients(dc.db, clientIDs, documentIds)
}

func (dc *DocumentUseCase) checkclientIsAbleToViewDocument(
	clientID string,
	refType string,
	refID string,
) (bool, error) {
	if refType == rfrl.ClientRef {
		return true, nil
	}
	if refType == rfrl.SessionRef {
		//TODO: Need to figure out how to handle checkclientIsAbleToViewDocument
		// clientIDs, err := session.GetclientsIdForSession(h.db, refID)

		// for i := 0; i < len(clientIDs); i += 1 {
		// 	if clientIDs[i] == clientID {
		// 		return true, nil
		// 	}
		// }
		// return false, err
		return false, nil
	}
	return false, nil
}

func (dc *DocumentUseCase) checkclientIsInRef(
	clientID string,
	refType string,
	refID string,
) (bool, error) {
	if refType == rfrl.ClientRef {
		return refID == clientID, nil
	}
	if refType == rfrl.SessionRef {
		//TODO: Need to figure out how to handle checkclientIsAbleToViewDocument
		// clientIDs, err := session.GetclientsIdForSession(h.db, refID)

		// for i := 0; i < len(clientIDs); i += 1 {
		// 	if clientIDs[i] == clientID {
		// 		return true, nil
		// 	}
		// }
		// return false, err
		return false, nil
	}

	return false, nil
}
