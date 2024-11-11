package transport

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/Intiqo/app-platform/internal/domain"
)

func GetClaimsForContext(ctx echo.Context) (result domain.Claims) {
	// Get the user from the context
	token := ctx.Get("user")

	// Get the user ID and role from the token
	if token != nil {
		u := token.(*jwt.Token)
		jwtClaims := u.Claims.(jwt.MapClaims)
		if jwtClaims["user_id"] != nil && jwtClaims["user_id"].(string) != "" {
			result.UserID = uuid.Must(uuid.FromString(jwtClaims["user_id"].(string)))
		}
		if jwtClaims["role"] != nil && jwtClaims["role"].(string) != "" {
			result.Role = jwtClaims["role"].(string)
		}
	}

	// Return the result
	return result
}
