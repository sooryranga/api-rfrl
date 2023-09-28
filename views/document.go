package views

import (
	"net/http"
	"strconv"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
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

type DocumentView struct {
	DocumentUseCase tutorme.DocumentUseCase
}

// CreateDocumentEndpoint view is an endpoint used to create document
func (dv *DocumentView) CreateDocumentEndpoint(c echo.Context) error {
	payload := DocumentPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := tutorme.GetClaims(c)

	document, err := dv.DocumentUseCase.CreateDocument(
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
func (dv *DocumentView) UpdateDocumentEndpoint(c echo.Context) error {
	payload := DocumentPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := tutorme.GetClaims(c)

	client, err := dv.DocumentUseCase.UpdateDocument(
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
func (dv *DocumentView) DeleteDocumentEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed is not valid integer"))
	}
	claims := tutorme.GetClaims(c)

	err = dv.DocumentUseCase.DeleteDocument(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetDocumentEndpoint view is an endpoint used to get document
func (dv *DocumentView) GetDocumentEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed is not valid integer"))
	}
	claims := tutorme.GetClaims(c)

	client, err := dv.DocumentUseCase.GetDocument(ID, claims.ClientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, client)
}

// CreateDocumentOrderEndpoint view is an endpoint used to create document order
func (dv *DocumentView) CreateDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := tutorme.GetClaims(c)

	if payload.RefType == "user" && claims.ClientID != payload.RefID {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"You are unauthorized to create this document-order",
		)
	}

	listOfDocuments, err := dv.DocumentUseCase.CreateDocumentOrder(
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
func (dv *DocumentView) UpdateDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := tutorme.GetClaims(c)

	listOfDocuments, err := dv.DocumentUseCase.UpdateDocumentOrder(
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
func (dv *DocumentView) GetDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := tutorme.GetClaims(c)

	listOfDocuments, err := dv.DocumentUseCase.GetDocumentOrder(
		claims.ClientID,
		payload.RefType,
		payload.RefID,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, listOfDocuments)
}
