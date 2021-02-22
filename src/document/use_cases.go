package document

func (h *Handler) createDocument(
	userID string,
	src string,
	name string,
	description string,
) (*Document, error) {
	document := NewDocument(userID, src, name, description)
	return CreateDocument(h.db, document)
}

func (h *Handler) updateDocument(
	userID string,
	ID int,
	src string,
	name string,
	description string,
) (*Document, error) {
	document := NewDocument(userID, src, name, description)
	return UpdateDocument(h.db, document)
}

func (h *Handler) deleteDocument(
	userID string,
	ID int,
) error {
	return DeleteDocument(h.db, ID, userID)
}

func (h *Handler) getDocument(
	id int,
	userID string,
) (*Document, error) {
	return getDocument(h.db, ID, userID)
}

func (h *Handler) createDocumentOrder(
	userId string,
	documentIds []int,
	refID string,
	refType string,
) ([]Document, error) {
	if err := h.checkUserIsInRef(userId, refType, refID); err != nil {
		return err
	}

	if err := h.checkDocumentsAreForRef(documentIds, refType, refID); err != nil {
		return nil, err
	}

	return CreateDocumentOrder(h.db, documentIds, refType, refID)
}

func (h *Handler) updateDocumentOrder(
	userId string,
	documentIds []int,
	refID string,
	refType string,
) ([]Document, error) {
	if err := h.checkUserIsInRef(userId, refType, refID); err != nil {
		return err
	}

	if err := checkDocumentsAreForRef(documentIds, refType, refID); err != nil {
		return nil, err
	}

	return UpdateDocumentOrder(h.db, documentIds, refType, refID)
}

func (h *Handler) getDocumentOrder(
	userID string,
	refID string,
	refType string,
) ([]Document, error) {
	if err := h.checkUserIsInRef(userId, refType, refID); err != nil {
		return err
	}

	return GetDocumentOrder(h.db, refType, refID)
}
