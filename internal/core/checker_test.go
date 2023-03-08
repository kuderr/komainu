package checker

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type FakeInfoStorage struct {
	clients      []uuid.UUID
	clientGroups map[uuid.UUID][]uuid.UUID
	apis         []Api
	apiRoutes    map[uuid.UUID][]Route
	routeClients map[uuid.UUID][]uuid.UUID
	routeGroups  map[uuid.UUID][]uuid.UUID
}

func NewFakeStorage() *FakeInfoStorage {
	apiID, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	route1ID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	route2ID, _ := uuid.Parse("22222222-2222-2222-2222-222222222222")
	route3ID, _ := uuid.Parse("22222222-2222-2222-2222-222222222220")

	client1, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	client2, _ := uuid.Parse("44444444-4444-4444-4444-444444444444")

	grp1, _ := uuid.Parse("55555555-5555-5555-5555-555555555555")

	return &FakeInfoStorage{
		clients:      []uuid.UUID{client1, client2},
		clientGroups: map[uuid.UUID][]uuid.UUID{client1: {grp1}},

		apis:         []Api{{apiID, "https://test.com"}},
		apiRoutes:    map[uuid.UUID][]Route{apiID: {{route1ID, "GET", "/test"}, {route2ID, "GET", "/test/{id}"}, {route3ID, "POST", "/test"}}},
		routeClients: map[uuid.UUID][]uuid.UUID{route1ID: {client1, client2}, route2ID: {client1, client2}, route3ID: {client1}},
		routeGroups:  map[uuid.UUID][]uuid.UUID{route1ID: {grp1}},
	}
}

func (s *FakeInfoStorage) GetClients() ([]uuid.UUID, error) {
	return s.clients, nil
}

func (s *FakeInfoStorage) GetClientGroups(clientID uuid.UUID) ([]uuid.UUID, error) {
	groups, ok := s.clientGroups[clientID]
	if !ok {
		return []uuid.UUID{}, nil
	}

	return groups, nil
}

func (s *FakeInfoStorage) GetApis() ([]Api, error) {
	return s.apis, nil
}

func (s *FakeInfoStorage) GetApiRoutes(apiID uuid.UUID) ([]Route, error) {
	routes, ok := s.apiRoutes[apiID]
	if !ok {
		return []Route{}, nil
	}

	return routes, nil
}

func (s *FakeInfoStorage) GetRouteClients(routeID uuid.UUID) ([]uuid.UUID, error) {
	clients, ok := s.routeClients[routeID]
	if !ok {
		return []uuid.UUID{}, nil
	}

	return clients, nil
}

func (s *FakeInfoStorage) GetRouteGroups(routeID uuid.UUID) ([]uuid.UUID, error) {
	clients, ok := s.routeGroups[routeID]
	if !ok {
		return []uuid.UUID{}, nil
	}

	return clients, nil
}

type checkerTestSuite struct {
	suite.Suite
	storage  *FakeInfoStorage
	builder  *Builder
	authInfo *AuthInfo
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(checkerTestSuite))
}

func (suite *checkerTestSuite) SetupSuite() {
	suite.storage = NewFakeStorage()
	suite.builder = NewBuilder(suite.storage)

	client1, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	client2, _ := uuid.Parse("44444444-4444-4444-4444-444444444444")

	grp1, _ := uuid.Parse("55555555-5555-5555-5555-555555555555")

	accesses, clients, err := suite.builder.BuildAccessMap()
	suite.Require().NoError(err)
	suite.Require().Equal(clients, ClientsMap{client1: {client1, grp1}, client2: {client2}})
	suite.Require().Equal(accesses, AccessMap{"https://test.com": {"GET": {"/test": {client1, client2, grp1}, `/test/[\w-]+`: {client1, client2}}, "POST": {"/test": {client1}}}})

	suite.authInfo = NewAuthInfo(accesses, clients)
}

func (suite *checkerTestSuite) TestCheckAccessSuccess() {
	client1, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "GET",
	}, client1)
	suite.Require().NoError(err)
	suite.Require().Equal(isAccess, true)
}

func (suite *checkerTestSuite) TestCheckPatternAccessSuccess() {
	client1, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test/12313124",
		Method: "GET",
	}, client1)
	suite.Require().NoError(err)
	suite.Require().Equal(isAccess, true)
}

func (suite *checkerTestSuite) TestCheckAccessForbidden() {
	client2, _ := uuid.Parse("44444444-4444-4444-4444-444444444444")
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "POST",
	}, client2)
	suite.Require().NoError(err)
	suite.Require().Equal(isAccess, false)
}

func (suite *checkerTestSuite) TestCheckAccessClientNotFound() {
	client3, _ := uuid.Parse("55555555-5555-5555-5555-555555555555")
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "GET",
	}, client3)
	suite.Require().Equal(err, &NotFoundError{msg: "Client with ID 55555555-5555-5555-5555-555555555555 not found"})
	suite.Require().Equal(isAccess, false)
}

func (suite *checkerTestSuite) TestCheckAccessApiNotFound() {
	client1, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test2.com",
		Path:   "/test",
		Method: "GET",
	}, client1)
	suite.Require().Equal(err, &NotFoundError{msg: "Not found api with url https://test2.com"})
	suite.Require().Equal(isAccess, false)
}

func (suite *checkerTestSuite) TestCheckAccesPathsNotFound() {
	client1, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	isAccess, err := suite.authInfo.CheckAccess(&AccessData{
		ApiUrl: "https://test.com",
		Path:   "/test",
		Method: "DELETE",
	}, client1)
	suite.Require().Equal(err, &NotFoundError{msg: "Not found paths for DELETE https://test.com"})
	suite.Require().Equal(isAccess, false)
}

func (suite *checkerTestSuite) TestUpdate() {
	client1, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	client4, _ := uuid.Parse("66666666-6666-6666-6666-666666666666")

	suite.authInfo.Update(
		AccessMap{"https://test.com": {"GET": {"/test": {client1}}}},
		ClientsMap{client4: {client4}},
	)

	suite.Require().Equal(suite.authInfo.accesses, AccessMap{"https://test.com": {"GET": {"/test": {client1}}}})
	suite.Require().Equal(suite.authInfo.clients, ClientsMap{client4: {client4}})
}
