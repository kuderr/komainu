package auther

import (
	"auther/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type Service struct {
	queries *database.Queries
	secret  string
}

func NewService(queries *database.Queries, secret string) *Service {
	return &Service{queries: queries, secret: secret}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/auth", s.CheckAccess)
	router.POST("/auth/", s.CheckAccess)
}

type accessData struct {
	Token  string `json:"token,omitempty" binding:"required"`
	ApiUrl string `json:"api_url,omitempty" binding:"required"`
	Path   string `json:"path,omitempty" binding:"required"`
	Method string `json:"method,omitempty" binding:"required"`
}

func (s *Service) CheckAccess(c *gin.Context) {
	// Parse request
	var request accessData
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := decodeToken(request.Token, s.secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasAccess, err := s.checkAccessInDB(c, request, token.Subject.ClientName)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	if hasAccess {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Access permit", "access": true, "client": token.Subject.ClientName})
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Access denied", "access": false, "client": token.Subject.ClientName})
	}
}
