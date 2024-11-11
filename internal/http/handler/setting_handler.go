package handler

import (
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"

	"github.com/Intiqo/app-platform/internal/domain"
	"github.com/Intiqo/app-platform/internal/http/transport"
)

// SettingHandler represents a handler for the Setting entity
type SettingHandler struct {
	s domain.SettingService
}

// NewSettingHandler creates a new instance of the setting handler
func NewSettingHandler(s domain.SettingService) SettingHandler {
	return SettingHandler{
		s: s,
	}
}

// FindByID finds a setting by ID
//
//	@Summary		Find a setting by id
//	@Description	Find a setting by id
//	@Tags			Setting
//	@ID				findSettingByID
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			id	path		string	true	"Setting ID"
//	@Success		200	{object}	domain.BaseResponse{data=domain.Setting}
//	@Failure		400	{object}	domain.ErrorResponse
//	@Failure		401	{object}	domain.ErrorResponse
//	@Failure		403	{object}	domain.ErrorResponse
//	@Failure		500	{object}	domain.ErrorResponse
//	@Router			/setting/{id} [get]
func (c SettingHandler) FindByID(ctx echo.Context) (err error) {
	// Parse the ID from the path parameter
	id, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		return err
	}

	// Find the setting by ID
	result, err := c.s.FindByID(id)
	if err != nil {
		return err
	}

	// Return the result
	return transport.SendResponse(ctx, http.StatusOK, result)
}

// Filter filters settings by criteria
//
//	@Summary		Filter settings by criteria
//	@Description	Filter settings by criteria. Supports pagination and returns the number of records as total.
//	@Tags			Setting
//	@ID				filterSettingsByCriteria
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			page	query		number									false	"Page Index"
//	@Param			size	query		number									false	"Page Size"
//	@Param			in		body		domain.FilterSettingsByCriteriaInput	true	"Input"
//	@Success		200		{object}	domain.PaginationResponse{data=[]domain.Setting}
//	@Failure		400		{object}	domain.ErrorResponse
//	@Failure		401		{object}	domain.ErrorResponse
//	@Failure		403		{object}	domain.ErrorResponse
//	@Failure		500		{object}	domain.ErrorResponse
//	@Router			/setting/filter [post]
func (c SettingHandler) Filter(ctx echo.Context) (err error) {
	// Parse the input from the request body
	var in domain.FilterSettingsByCriteriaInput
	err = transport.DecodeAndValidateRequestBody(ctx, &in)
	if err != nil {
		return err
	}

	// Decode the query options
	opts := transport.DecodeQueryOptions(ctx)

	// Filter the settings
	result, total, err := c.s.Filter(in, opts)
	if err != nil {
		return err
	}

	// Return the result
	return transport.SendPaginationResponse(ctx, http.StatusOK, result, total)
}
