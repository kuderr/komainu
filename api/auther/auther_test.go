// [WIP]

package auther

import (
	"auther/config"
	"auther/internal/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	router  *gin.Engine
	queries *database.Queries
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupSuite() {
	cfg, err := config.Read("env", ".env", ".")
	suite.Require().NoError(err)

	postgres, err := database.NewPostgres(cfg.PostgresUrl)
	suite.Require().NoError(err)

	suite.queries = database.New(postgres.DB)
	service := NewService(suite.queries, "xxx")

	suite.router = gin.Default()
	service.RegisterHandlers(suite.router)

}

type response struct {
	message string
	access  bool
	client  string
}

// dummy test token, encoded with xxx secret above
const Token string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOnsiY2xpZW50X25hbWUiOiJ0ZXN0In19.7UXTvMN42-p2b76cGn0R17YW4P9UDpGT07D8pdfneuM"

func (suite *ServiceTestSuite) TestCheckAccess() {
	request := accessData{
		Token:  Token,
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "/GET",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("GET", "/auth", &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	var resp response
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&resp))
	suite.Require().Equal(resp.access, true)
	suite.Require().Equal(resp.message, "Access permit")
	suite.Require().Equal(resp.client, "")
}
