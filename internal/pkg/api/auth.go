package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func ExtractToken(c *gin.Context) (string, string, error) {
	authorizationHeader := c.GetHeader("Authorization")

	if len(authorizationHeader) == 0 {
		return "", "", fmt.Errorf("authorization header is not set")
	}

	components := strings.Split(authorizationHeader, " ")

	if len(components) != 2 {
		return "", "", fmt.Errorf("invalid authorization header")
	}

	return components[0], components[1], nil
}
