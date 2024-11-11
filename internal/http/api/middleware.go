package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/Intiqo/app-platform/internal/domain"
	"github.com/Intiqo/app-platform/internal/http/swagger"
	"github.com/Intiqo/app-platform/internal/http/transport"
)

// SetupMiddleware sets up middleware for the echo server
func (t AppApi) SetupMiddleware(e *echo.Echo) {
	// Set up the validator middleware
	e.Validator = &transport.CustomValidator{Validator: validator.New()}
	// Set up the error handler middleware
	e.HTTPErrorHandler = errorMiddleware
	// Set the request body limit to 10M
	rqsl := t.cfg.RequestBodySizeLimit
	if rqsl == "" {
		rqsl = "10M"
	}
	e.Use(echomiddleware.BodyLimit(rqsl))
	// Recovery middleware recovers from panics anywhere in the chain,
	e.Use(echomiddleware.Recover())
	// Add request ID middleware
	e.Use(echomiddleware.RequestID())
	// Add CORS middleware
	e.Use(echomiddleware.CORS())
	// Add the swagger redirect middleware
	e.Use(swagger.RedirectSwagger)
	// Add the logger middleware
	e.Use(
		echomiddleware.LoggerWithConfig(
			echomiddleware.LoggerConfig{
				Format: "${time_rfc3339} ${id} ${remote_ip} ${method} ${uri} ${latency_human} ${status} ${error}\n",
			},
		),
	)
}

// errorMiddleware absorbs and processes all errors
func errorMiddleware(err error, c echo.Context) {
	switch err.(type) {
	case *echo.HTTPError:
		err := err.(*echo.HTTPError)
		switch err.Code {
		case http.StatusUnauthorized:
			_ = c.JSON(err.Code, domain.UnauthorizedError{
				Code:    domain.ErrorCodeUNAUTHORIZED,
				Message: err.Message.(string),
			})
		case http.StatusForbidden:
			_ = c.JSON(err.Code, domain.ForbiddenAccessError{
				Code:    domain.ErrorCodeFORBIDDENACCESS,
				Message: err.Message.(string),
			})
		case http.StatusNotFound:
			_ = c.JSON(err.Code, domain.NotFoundError{})
		case http.StatusBadRequest:
			_ = c.JSON(err.Code, domain.InvalidRequestError{Message: err.Message.(string)})
		default:
			_ = c.JSON(err.Code, domain.SystemError{Code: domain.ErrorCodeINTERNALSERVERERROR, Message: err.Message.(string)})
		}

	case validator.ValidationErrors:
		var ve error
		fields := make([]string, 0)
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			if e.Tag() == "required" {
				fields = append(fields, fmt.Sprintf("%s is required", e.Field()))
				continue
			}

			if e.Tag() == "e164" {
				fields = append(fields, fmt.Sprintf("%s is an invalid mobile number", e.Field()))
				continue
			}

			if e.Tag() == "email" {
				fields = append(fields, fmt.Sprintf("%s is an invalid email address", e.Field()))
				continue
			}

			if e.Tag() == "oneof" {
				fields = append(fields, fmt.Sprintf("%s must be one of %s", e.Field(), e.Param()))
				continue
			}

			if e.Tag() == "min" {
				fields = append(fields, fmt.Sprintf("%s must be %s characters minimum", e.Field(), e.Param()))
				continue
			}

			if e.Tag() == "max" {
				fields = append(fields, fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param()))
				continue
			}

		}

		ve = domain.ValidationError{
			Code:    domain.ErrorCodeVALIDATIONERROR,
			Message: domain.MessageVALIDATIONFAILED,
			Fields:  fields,
		}
		_ = c.JSON(http.StatusBadRequest, ve)

	case *pgconn.PgError:
		res := domain.SystemError{
			Code:    domain.ErrorCodeINTERNALSERVERERROR,
			Message: err.Error(),
		}
		_ = c.JSON(http.StatusInternalServerError, res)

	case domain.DataNotFoundError:
		res := domain.UserError{
			Code:    domain.ErrorCodeINVALIDREQUEST,
			Message: err.Error(),
		}
		_ = c.JSON(http.StatusBadRequest, res)

	case domain.UserError:
		res := domain.UserError{
			Code:    domain.ErrorCodeINVALIDREQUEST,
			Message: err.Error(),
		}
		_ = c.JSON(http.StatusBadRequest, res)

	case domain.UnauthorizedError:
		res := domain.UnauthorizedError{
			Code:    domain.ErrorCodeUNAUTHORIZED,
			Message: domain.MessageUNAUTHORIZEDACCESS,
		}
		_ = c.JSON(http.StatusUnauthorized, res)

	case domain.ForbiddenAccessError:
		res := domain.ForbiddenAccessError{
			Code:    domain.ErrorCodeFORBIDDENACCESS,
			Message: domain.MessageFORBIDDENACCESS,
		}
		_ = c.JSON(http.StatusForbidden, res)

	default:
		res := domain.SystemError{
			Code:    domain.ErrorCodeINTERNALSERVERERROR,
			Message: err.Error(),
		}
		_ = c.JSON(http.StatusInternalServerError, res)
	}
}
