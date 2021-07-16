package views

import (
	"net/http"
	"strconv"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	// DocumentPayload is the struct used to hold payload from /document
	DocumentPayload struct {
		ID          int    `path:"id"`
		Src         string `json:"src" validate:"required,url"`
		Name        string `json:"name" validate:"required,gte=0,lte=100"`
		Description string `json:"description" validate:"omitempty,base64"`
	}

	// DocumentOrderPayload is the struct used to hold payload from /document-order
	DocumentOrderPayload struct {
		RefType     string `json:"refType" query:"ref_type" validate:"required,oneof= client session"`
		RefID       string `json:"refId" query:"ref_id" validate:"gte=0,base64"`
		DocumentIDs []int  `json:"documentIds" validate:"required,gt=0,dive,required,numeric"`
	}
)

type DocumentView struct {
	DocumentUseCase rfrl.DocumentUseCase
}

// CreateDocumentEndpoint view is an endpoint used to create document
func (dv *DocumentView) CreateDocumentEndpoint(c echo.Context) error {
	payload := DocumentPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return err
	}

	document, err := dv.DocumentUseCase.CreateDocument(
		claims.ClientID,
		payload.Src,
		payload.Name,
		payload.Description,
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

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return err
	}

	client, err := dv.DocumentUseCase.UpdateDocument(
		claims.ClientID,
		payload.ID,
		payload.Src,
		payload.Name,
		payload.Description,
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
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return err
	}

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
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return err
	}

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

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return err
	}

	if payload.RefType == "user" && claims.ClientID != payload.RefID {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"You are unauthorized to create this document-order",
		)
	}

	listOfDocuments, err := dv.DocumentUseCase.CreateDocumentOrder(
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

// UpdateDocumentOrderEndpoint is used to update document order
func (dv *DocumentView) UpdateDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return err
	}

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

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return err
	}

	listOfDocuments, err := dv.DocumentUseCase.GetDocumentOrder(
		claims.ClientID,
		payload.RefID,
		payload.RefType,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, listOfDocuments)
}
