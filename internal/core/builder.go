package checker

import (
	"regexp"

	"github.com/google/uuid"
)

type AuthInfoStorage interface {
	GetClients() ([]uuid.UUID, error)
	GetClientGroups(uuid.UUID) ([]uuid.UUID, error)
	GetApis() ([]Api, error)
	GetApiRoutes(uuid.UUID) ([]Route, error)
	GetRouteClients(uuid.UUID) ([]uuid.UUID, error)
	GetRouteGroups(uuid.UUID) ([]uuid.UUID, error)
}

type Builder struct {
	storage AuthInfoStorage
}

func NewBuilder(storage AuthInfoStorage) *Builder {
	return &Builder{storage: storage}
}

// TODO: Maybe pass pointers
func (builder *Builder) BuildAccessMap() (AccessMap, ClientsMap, error) {
	accesses := AccessMap{}
	clients := ClientsMap{}

	cs, err := builder.storage.GetClients()
	if err != nil {
		return AccessMap{}, ClientsMap{}, err
	}

	// Build client entities
	for _, clientID := range cs {
		groups, err := builder.storage.GetClientGroups(clientID)
		if err != nil {
			return AccessMap{}, ClientsMap{}, err
		}
		clients[clientID] = append([]uuid.UUID{clientID}, groups...)
	}

	apis, err := builder.storage.GetApis()
	if err != nil {
		return AccessMap{}, ClientsMap{}, err
	}

	for _, api := range apis {
		apiRoutes, err := builder.storage.GetApiRoutes(api.ID)
		if err != nil {
			return AccessMap{}, ClientsMap{}, err
		}

		for _, route := range apiRoutes {
			routeClients, err := builder.storage.GetRouteClients(route.ID)
			if err != nil {
				return AccessMap{}, ClientsMap{}, err
			}

			routeGroups, err := builder.storage.GetRouteGroups(route.ID)
			if err != nil {
				return AccessMap{}, ClientsMap{}, err
			}

			routeClientEntities := append(routeClients, routeGroups...)
			if len(routeClientEntities) == 0 {
				// No Accesses for route
				continue
			}

			if _, ok := accesses[api.Url]; !ok {
				accesses[api.Url] = make(map[string]map[string][]uuid.UUID)
			}

			if _, ok := accesses[api.Url][route.Method]; !ok {
				accesses[api.Url][route.Method] = make(map[string][]uuid.UUID)
			}

			pathRegex := regexp.MustCompile(`({[\w-]+})`).ReplaceAllString(route.Path, `[\w-]+`)
			accesses[api.Url][route.Method][pathRegex] = routeClientEntities
		}
	}

	return accesses, clients, nil
}
