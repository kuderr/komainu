package auther

import (
	"auther/internal/database"
	"context"

	"github.com/gin-gonic/gin"
)

func (s *Service) checkAccessInDB(c *gin.Context, request accessData, clientName string) (bool, error) {
	clientId, err := s.queries.GetClientIdByName(context.Background(), clientName)
	if err != nil {
		return false, err
	}

	apiId, err := s.queries.GetApiIdByUrl(context.Background(), request.ApiUrl)
	if err != nil {
		return false, err
	}

	isApiAdmin, err := s.queries.CountAdminAssociation(context.Background(),
		database.CountAdminAssociationParams{
			ClientID: clientId,
			ApiID:    apiId,
		})
	if err != nil {
		return false, err
	}
	if isApiAdmin != 0 {
		return true, nil
	}

	routeId, err := s.queries.GetApiRouteIdByMethodAndPath(context.Background(),
		database.GetApiRouteIdByMethodAndPathParams{
			ApiID:  apiId,
			Method: request.Method,
			Path:   request.Path,
		})
	//  TODO: search for pattern route
	if err != nil {
		return false, err
	}

	hasRouteAccess, err := s.queries.CountRouteAssociation(context.Background(),
		database.CountRouteAssociationParams{
			ClientID: clientId,
			RouteID:  routeId,
		})
	if err != nil {
		return false, err
	}
	if hasRouteAccess != 0 {
		return true, nil
	}

	return false, nil
}
