package document

import (
	"net/http"
	"strconv"

	"github.com/Arun4rangan/api-tutorme/src/auth"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	// DocumentPayload is the struct used to hold payload from /document
	DocumentPayload struct {
		ID          int    `path:"id"`
		Src         string `json:"src" validate:"required,url"`
		Name        string `json:"name" validate:"required,gte=0,lte=100"`
		description string `json:"description" validate:"omitempty,base64"`
	}

	// DocumentOrderPayload is the struct used to hold payload from /document-order
	DocumentOrderPayload struct {
		RefType     string `json:"ref_type" validate:"required,oneof= user session"`
		RefID       string `json:"ref_id" validate:"gte=0,base64"`
		DocumentIDs []int  `json:"document_ids" validate:"required,gt=0,dive,required,numeric"`
	}
)

// CreateDocumentEndpoint view is an endpoint used to create document
func (h *Handler) CreateDocumentEndpoint(c echo.Context) error {
	payload := DocumentPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := auth.GetClaims(c)

	document, err := h.createDocument(
		claims.ClientID,
		payload.Src,
		payload.Name,
		payload.description,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, document)
}

// UpdateDocumentEndpoint view is an endpoint used to update document
func (h *Handler) UpdateDocumentEndpoint(c echo.Context) error {
	payload := DocumentPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := auth.GetClaims(c)

	client, err := h.updateDocument(
		claims.ClientID,
		payload.ID,
		payload.Src,
		payload.Name,
		payload.description,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}

// DeleteDocumentEndpoint view is an endpoint to delete document
func (h *Handler) DeleteDocumentEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed is not valid integer"))
	}
	claims := auth.GetClaims(c)

	err = h.deleteDocument(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetDocumentEndpoint view is an endpoint used to get document
func (h *Handler) GetDocumentEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed is not valid integer"))
	}
	claims := auth.GetClaims(c)

	client, err := h.getDocument(ID, claims.ClientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, client)
}

// CreateDocumentOrderEndpoint view is an endpoint used to create document order
func (h *Handler) CreateDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := auth.GetClaims(c)

	if payload.RefType == "user" && claims.ClientID != payload.RefID {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"You are unauthorized to create this document-order",
		)
	}

	listOfDocuments, err := h.createDocumentOrder(
		claims.ClientID,
		payload.DocumentIDs,
		payload.RefType,
		payload.RefID,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, listOfDocuments)
}

// UpdateDocumentOrderEndpoint is used to update document order
func (h *Handler) UpdateDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := auth.GetClaims(c)

	listOfDocuments, err := h.updateDocumentOrder(
		claims.ClientID,
		payload.DocumentIDs,
		payload.RefID,
		payload.RefType,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, listOfDocuments)
}

// GetDocumentOrderEndpoint is used to get documents in order
func (h *Handler) GetDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := auth.GetClaims(c)

	listOfDocuments, err := h.getDocumentOrder(
		claims.ClientID,
		payload.RefType,
		payload.RefID,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, listOfDocuments)
}
