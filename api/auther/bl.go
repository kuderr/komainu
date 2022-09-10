package auther

import (
	"auther/internal/database"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (s *Service) checkAccessInDB(c *gin.Context, request accessData, clientName string) (bool, error) {
	apiUrl := strings.TrimRight(request.ApiUrl, "/")
	path := strings.TrimRight(request.Path, "/")
	method := strings.ToUpper(request.Method)

	clientId, err := s.queries.GetClientIdByName(context.Background(), clientName)
	if err != nil {
		return false, err
	}

	apiId, err := s.queries.GetApiIdByUrl(context.Background(), apiUrl)
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
			Method: method,
			Path:   path,
		})
	// TODO: refactor
	if err != nil {
		if err != pgx.ErrNoRows {
			return false, err
		} else {
			// Search for pattern route
			routes, err := s.queries.GetApiRoutesByMethod(context.Background(),
				database.GetApiRoutesByMethodParams{
					ApiID:  apiId,
					Method: method,
				})
			if err != nil {
				return false, err
			}

			var found bool
			routeId, found = searchPatternRoute(routes, path)
			if !found {
				// TODO: use custom errors
				return false, pgx.ErrNoRows
			}
		}
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

func searchPatternRoute(routes []database.Route, path string) (uuid.UUID, bool) {
	var patternRouteId uuid.UUID
	var found bool

	for _, route := range routes {
		pathRegex := regexp.MustCompile(`({[\w-]+})`).ReplaceAllString(route.Path, `[\w-]+`)
		re := regexp.MustCompile(fmt.Sprintf("^%s$", pathRegex))
		if re.MatchString(path) {
			patternRouteId = route.ID
			found = true
			break
		}
	}

	return patternRouteId, found
}
