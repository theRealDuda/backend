package handlers

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func normalizeUserID(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func getSessionUserID(c *gin.Context) (int, bool) {
	sess := sessions.Default(c)
	return normalizeUserID(sess.Get("user_id"))
}

func getUserID(c *gin.Context) (int, bool) {
	if ctxUserID, ok := c.Get("user_id"); ok {
		if id, ok := normalizeUserID(ctxUserID); ok {
			return id, true
		}
	}

	return getSessionUserID(c)
}
