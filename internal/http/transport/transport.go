package transport

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/Intiqo/app-platform/internal/domain"
)

const PageMax = 100

// DecodeQueryOptions decodes the query options
func DecodeQueryOptions(ctx echo.Context) (result domain.QueryOptions) {
	// Get the limit and offset
	limit, offset := GetLimitAndOffset(ctx)

	// Create the query options
	if limit != 0 {
		result.Limit = limit
	}
	if offset != 0 {
		result.Offset = offset
	}

	// Return the result
	return result
}

// DecodeAndValidateRequestBody decodes and validates the request body
func DecodeAndValidateRequestBody(ctx echo.Context, t interface{}) (err error) {
	// Bind the request body
	err = ctx.Bind(t)
	if err != nil {
		return err
	}

	// Validate the request body
	err = ctx.Validate(t)
	if err != nil {
		return err
	}

	// Return the result
	return nil
}

// GetLimitAndOffset gets the limit and offset from the query params
func GetLimitAndOffset(ctx echo.Context) (int64, int64) {
	p := ctx.QueryParam("page")
	s := ctx.QueryParam("size")
	page, _ := strconv.Atoi(p)
	if page < 0 {
		page = 0
	}

	size, _ := strconv.Atoi(s)
	switch {
	case size > PageMax:
		size = PageMax
	case size <= 0:
		size = PageMax
	}

	return int64(size), int64(page * size)
}

// SendResponse sends a response
func SendResponse(ctx echo.Context, status int, data interface{}) error {
	// Create the final response
	var finalResult domain.BaseResponse
	if data != nil {
		finalResult = domain.BaseResponse{
			Data: data,
		}
	}

	// If status is 204, return no content
	if status == http.StatusNoContent {
		return ctx.NoContent(status)
	}

	// Return the result
	return ctx.JSON(status, finalResult)
}

// SendPaginationResponse sends a paginated response
func SendPaginationResponse(ctx echo.Context, status int, data interface{}, total int64) error {
	// Get the page and size
	p, _ := strconv.Atoi(ctx.QueryParam("page"))
	s, _ := strconv.Atoi(ctx.QueryParam("size"))
	page := int64(p)
	if page <= 0 {
		page = 0
	}
	size := int64(s)
	if size <= 0 {
		size = PageMax
	}

	// Create the final response
	finalResult := domain.PaginationResponse{
		BaseResponse: domain.BaseResponse{
			Data: data,
		},
		Page:  page,
		Size:  size,
		Total: total,
	}

	// If status is 204, return no content
	if status == http.StatusNoContent {
		return ctx.NoContent(status)
	}

	// Return the result
	return ctx.JSON(status, finalResult)
}

// CustomValidator custom validator for echo
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates the request body
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
