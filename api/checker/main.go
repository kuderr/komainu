package api

import (
	checker "checker/internal/core"
	"checker/internal/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Token string `json:"token,omitempty" binding:"required"`
	checker.AccessData
}

type Service struct {
	authInfo     *checker.AuthInfo
	JWTPublicKey string
}

func NewService(authInfo *checker.AuthInfo, JWTPublicKey string) *Service {
	return &Service{authInfo: authInfo, JWTPublicKey: JWTPublicKey}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/auth", s.CheckAccess)
	router.POST("/auth/", s.CheckAccess)
}

func (s *Service) CheckAccess(c *gin.Context) {
	// Parse request
	var request AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := tokens.DecodeToken(request.Token, s.JWTPublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID, err := tokens.GetClientID(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasAccess, err := s.authInfo.CheckAccess(&request.AccessData, clientID)
	if err != nil {
		switch err.(type) {
		case *checker.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
	}

	// Build response
	if hasAccess {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Access permit", "access": true, "client": clientID})
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Access denied", "access": false, "client": clientID})
	}
}
