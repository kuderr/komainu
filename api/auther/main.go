package autherApi

import (
	"auther/internal/auther"
	"auther/internal/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct {
	authInfo *auther.AuthInfo
	secret   string
}

func NewService(authInfo *auther.AuthInfo, secret string) *Service {
	return &Service{authInfo: authInfo, secret: secret}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/auth", s.CheckAccess)
	router.POST("/auth/", s.CheckAccess)
}

func (s *Service) CheckAccess(c *gin.Context) {
	// Parse request
	var request auther.AccessData
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := tokens.DecodeToken(request.Token, s.secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasAccess, err := s.authInfo.CheckAccess(&request, token.Subject.ClientName)
	if err != nil {
		switch err.(type) {
		case *auther.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
	}

	// Build response
	if hasAccess {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Access permit", "access": true, "client": token.Subject.ClientName})
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Access denied", "access": false, "client": token.Subject.ClientName})
	}
}
