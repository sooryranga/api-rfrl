package document

import (
	"net/http"

	"github.com/Arun4rangan/api-tutorme/src/auth"
	"github.com/labstack/echo/v4"
)

type (
	// ClientPayload is the struct used to hold payload from /client
	DocumentPayload struct {
		ID          string `path:"id"`
		Src         string `json:"src" validate:"required,url"`
		Name        string `json:"name" validate:"required,gte=0,lte=100"`
		description string `json:"description" validate:"omitempty,base64"`
	}

	// EducationPaylod is the struct used to create education
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

	client, err := h.createDocument(
		claims.UserID,
		payload.Src,
		payload.Name,
		payload.description,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, client)
}

// UpdateDocumentEndpoint view is an endpoint used to update document
func (h *Handler) UpdateDocumentEndpoint(c echo.Context) error {
	payload := DocumentPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := auth.GetClaims(c)

	client, err := h.updateDocument(
		claims.UserID,
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
	id := c.Param("id")
	claims := auth.GetClaims(c)

	err := h.deleteDocument(id, claims.UserID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}

// GetDocumentEndpoint view is an endpoint used to get document
func (h *Handler) GetDocumentEndpoint(c echo.Context) error {
	id := c.Param("id")

	claims := auth.GetClaims(c)

	client, err := h.getDocuument(id, claims.UserID)

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

	order, err := h.createDocumentOrder(
		claims.UserID,
		payload.DocumentIDs,
		payload.RefID,
		payload.DocumentIDs,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, order)
}

// UpdateDocumentOrderEndpoint is used to update document order
func (h *Handler) UpdateDocumentOrderEndpoint(c echo.Context) error {
	payload := DocumentOrderPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := auth.GetClaims(c)

	listOfDocuments, err := h.createDocumentOrder(
		claims.UserID,
		payload.DocumentIDs,
		payload.RefID,
		payload.DocumentIDs,
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
		claims.UserID,
		payload.DocumentIDs,
		payload.RefID,
		payload.DocumentIDs,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, listOfDocuments)
}
