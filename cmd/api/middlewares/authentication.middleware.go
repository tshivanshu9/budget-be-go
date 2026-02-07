package middlewares

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/models"
	"gorm.io/gorm"
)

type AppMiddleware struct {
	DB *gorm.DB
}

func (am *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return common.SendUnauthorizedResponse(c, nil)
		}

		authHeaderSplit := strings.Split(authHeader, " ")
		accessToken := authHeaderSplit[1]

		claims, err := common.ParseJWT(accessToken)
		if err != nil || common.IsClaimExpired(claims) {
			return common.SendUnauthorizedResponse(c, nil)
		}

		var user models.UserModel
		result := am.DB.First(&user, claims.ID)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.Error != nil {
			return common.SendUnauthorizedResponse(c, nil)
		}

		c.Set("user", user)
		return next(c)
	}
}
