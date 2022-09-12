package auther

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type FakeInfoStorage struct {
	clients      []string
	apis         []Api
	apiRoutes    map[uuid.UUID][]Route
	routeClients map[uuid.UUID][]string
}

func NewFakeStorage() *FakeInfoStorage {
	apiID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	route1ID, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	route2ID, _ := uuid.Parse("22222222-2222-2222-2222-222222222222")

	return &FakeInfoStorage{
		clients:      []string{"test", "foo"},
		apis:         []Api{{apiID, "https://test.com"}},
		apiRoutes:    map[uuid.UUID][]Route{apiID: {{route1ID, "GET", "/test"}, {route2ID, "GET", "/test/{id}"}}},
		routeClients: map[uuid.UUID][]string{route1ID: {"test"}, route2ID: {"test"}},
	}
}

func (s *FakeInfoStorage) GetClients() ([]string, error) {
	return s.clients, nil
}

func (s *FakeInfoStorage) GetApis() ([]Api, error) {
	return s.apis, nil
}

func (s *FakeInfoStorage) GetApiRoutes(apiID uuid.UUID) ([]Route, error) {
	routes, ok := s.apiRoutes[apiID]
	if !ok {
		return nil, fmt.Errorf("Api routes not found")
	}

	return routes, nil
}

func (s *FakeInfoStorage) GetRouteClients(routeID uuid.UUID) ([]string, error) {
	clients, ok := s.routeClients[routeID]
	if !ok {
		return nil, fmt.Errorf("Route clients not found")
	}

	return clients, nil
}

type AutherTestSuite struct {
	suite.Suite
	storage  *FakeInfoStorage
	builder  *Builder
	authInfo *AuthInfo
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AutherTestSuite))
}

func (suite *AutherTestSuite) SetupSuite() {
	suite.storage = NewFakeStorage()
	suite.builder = NewBuilder(suite.storage)

	accesses, clients, err := suite.builder.BuildAccessMap()
	suite.Require().NoError(err)
	suite.Require().Equal(clients, ClientsSet{"test": struct{}{}, "foo": struct{}{}})
	suite.Require().Equal(accesses, AccessMap{"https://test.com": {"GET": {"/test": {"test"}, `/test/[\w-]+`: {"test"}}}})

	suite.authInfo = NewAuthInfo(accesses, clients)
}

func (suite *AutherTestSuite) TestCheckAccessSuccess() {
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "GET",
	}, "test")
	suite.Require().NoError(err)
	suite.Require().Equal(isAccess, true)
}

func (suite *AutherTestSuite) TestCheckPatternAccessSuccess() {
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test/12313124",
		Method: "GET",
	}, "test")
	suite.Require().NoError(err)
	suite.Require().Equal(isAccess, true)
}

func (suite *AutherTestSuite) TestCheckAccessForbidden() {
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "GET",
	}, "foo")
	suite.Require().NoError(err)
	suite.Require().Equal(isAccess, false)
}

func (suite *AutherTestSuite) TestCheckAccessClientNotFound() {
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "GET",
	}, "bar")
	suite.Require().Equal(err, &NotFoundError{msg: "Client bar not found"})
	suite.Require().Equal(isAccess, false)
}

func (suite *AutherTestSuite) TestCheckAccessApiNotFound() {
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test2.com",
		Path:   "/test",
		Method: "GET",
	}, "test")
	suite.Require().Equal(err, &NotFoundError{msg: "Not found api with url https://test2.com"})
	suite.Require().Equal(isAccess, false)
}

func (suite *AutherTestSuite) TestCheckAccesPathsNotFound() {
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "POST",
	}, "test")
	suite.Require().Equal(err, &NotFoundError{msg: "Not found paths for POST https://test.com"})
	suite.Require().Equal(isAccess, false)
}

func (suite *AutherTestSuite) TestUpdate() {
	suite.authInfo.Update(AccessMap{"https://test.com": {"GET": {"/test": {"test"}}}},
		ClientsSet{"test": struct{}{}})

	suite.Require().Equal(suite.authInfo.accesses, AccessMap{"https://test.com": {"GET": {"/test": {"test"}}}})
	suite.Require().Equal(suite.authInfo.clients, ClientsSet{"test": struct{}{}})
}
