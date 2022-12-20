package auth

import (
	"errors"
	"net/http"
	"strings"
	"videown-server/internal/ginlet/resp"
	"videown-server/internal/service/auth"

	"github.com/gin-gonic/gin"
)

const BEARER_PREFIX = "Bearer "

func jwtBearerRequired(c *gin.Context) (bool, *auth.CustomClaims) {
	bearer := c.Request.Header.Get("Authorization")
	if bearer == "" {
		resp.ErrorWithHttpStatus(c, errors.New("Authorization field in header cannot be found"), http.StatusUnauthorized)
		c.Abort()
		return false, nil
	}
	jwtStr := strings.TrimPrefix(bearer, BEARER_PREFIX)
	// decode jwt => claims
	claims, err := auth.Jwth().ParseToken(jwtStr)
	if err != nil {
		resp.ErrorWithHttpStatus(c, err, http.StatusUnauthorized)
		c.Abort()
		return false, nil
	}
	return true, claims
}

func AuthRequiredForUser(c *gin.Context) {
	if ok, claims := jwtBearerRequired(c); ok && claims != nil {
		auth.StoreAuth(c, claims, false)
	}
}

func AuthRequiredForAdmin(c *gin.Context) {
	if ok, claims := jwtBearerRequired(c); ok && claims != nil {
		if !claims.IsAdmin {
			resp.ErrorWithHttpStatus(c, errors.New("no access permissions"), http.StatusForbidden)
			c.Abort()
			return
		}
		auth.StoreAuth(c, claims, true)
	}
}
