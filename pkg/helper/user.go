package helper

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func GetLoginUserFromContext(c *gin.Context) (*entity.User, bool) {
	userValue, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "failed to get user login", nil)
		return nil, false
	}

	user, ok := userValue.(*entity.User)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user login", nil)
		return nil, false
	}

	return user, true
}
