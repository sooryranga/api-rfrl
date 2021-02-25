package document

import (
	"fmt"

	"github.com/pkg/errors"
)

func (h *Handler) createDocument(
	clientID string,
	src string,
	name string,
	description string,
) (*Document, error) {
	document := NewDocument(clientID, src, name, description)
	return CreateDocument(h.db, document)
}

func (h *Handler) updateDocument(
	clientID string,
	ID int,
	src string,
	name string,
	description string,
) (*Document, error) {
	document := NewDocument(clientID, src, name, description)
	return UpdateDocument(h.db, ID, document)
}

func (h *Handler) deleteDocument(
	clientID string,
	ID int,
) error {
	tx, err := h.db.Beginx()
	if err != nil {
		return err
	}

	err = RemoveAndRenumberDocumentsOrder(tx, ID, clientID)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return errors.Wrap(rb, err.Error())
		}
		return err
	}

	err = DeleteDocument(tx, ID, clientID)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return errors.Wrap(rb, err.Error())
		}
		return err
	}

	return nil
}

func (h *Handler) getDocument(
	id int,
	clientID string,
) (*Document, error) {
	return GetDocument(h.db, id, clientID)
}

func (h *Handler) createDocumentOrder(
	clientID string,
	documentIds []int,
	refID string,
	refType string,
) ([]Document, error) {
	check, err := h.checkclientIsInRef(clientID, refType, refID)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf("Unauthorized to create document-order")
	}

	check, err = h.checkDocumentsAreForRef(documentIds, refType, refID)

	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf(
			"Documents are not related to client: Document - %v",
			documentIds,
		)
	}

	return CreateDocumentOrder(h.db, documentIds, refType, refID)
}

func (h *Handler) updateDocumentOrder(
	clientID string,
	documentIds []int,
	refID string,
	refType string,
) ([]Document, error) {
	check, err := h.checkclientIsInRef(clientID, refType, refID)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf("Unauthorized to update document-order")
	}

	check, err = h.checkDocumentsAreForRef(documentIds, refType, refID)

	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.Errorf(
			"Documents are not related to client: Document - %v",
			documentIds,
		)
	}

	return UpdateDocumentOrder(h.db, documentIds, refType, refID)
}

func (h *Handler) getDocumentOrder(
	clientID string,
	refID string,
	refType string,
) ([]Document, error) {
	check, err := h.checkclientIsInRef(clientID, refType, refID)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, errors.New(
			fmt.Sprintf("Client is not part of the ref_type (%s, %s)", refType, refID),
		)
	}

	return GetDocumentOrder(h.db, refType, refID)
}

func (h *Handler) checkDocumentsAreForRef(
	documentIds []int,
	refType string,
	refID string,
) (bool, error) {
	var clientIDs []string
	if refType == ClientRef {
		clientIDs = []string{refID}
	} else if refType == SessionRef {
		// clientIDs = session.GetClientsIdForSession(h.db, refID)
		return false, nil
	} else {
		return false, errors.Errorf("reftype (%v) not found", refType)
	}

	return CheckDocumentsBelongToclients(h.db, clientIDs, documentIds)
}

func (h *Handler) checkclientIsInRef(
	clientID string,
	refType string,
	refID string,
) (bool, error) {
	if refType == ClientRef {
		return refID != clientID, nil
	}
	if refType == SessionRef {
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
