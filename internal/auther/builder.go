package auther

import (
	"regexp"

	"github.com/google/uuid"
)

type AuthInfoStorage interface {
	GetClients() ([]string, error)
	GetApis() ([]Api, error)
	GetApiRoutes(uuid.UUID) ([]Route, error)
	GetRouteClients(uuid.UUID) ([]string, error)
}

type Builder struct {
	storage AuthInfoStorage
}

func NewBuilder(storage AuthInfoStorage) *Builder {
	return &Builder{storage: storage}
}

// TODO: Maybe pass pointers
func (builder *Builder) BuildAccessMap() (AccessMap, ClientsSet, error) {
	accesses := AccessMap{}
	clients := ClientsSet{}

	cs, err := builder.storage.GetClients()
	if err != nil {
		return AccessMap{}, ClientsSet{}, nil
	}
	for _, client := range cs {
		clients[client] = struct{}{}
	}

	apis, err := builder.storage.GetApis()
	if err != nil {
		return AccessMap{}, ClientsSet{}, nil
	}

	for _, api := range apis {
		apiRoutes, err := builder.storage.GetApiRoutes(api.ID)
		if err != nil {
			return AccessMap{}, ClientsSet{}, nil
		}

		for _, route := range apiRoutes {
			routeClients, err := builder.storage.GetRouteClients(route.ID)
			if err != nil {
				return AccessMap{}, ClientsSet{}, nil
			}

			if _, ok := accesses[api.Url]; !ok {
				accesses[api.Url] = make(map[string]map[string][]string)
			}

			if _, ok := accesses[api.Url][route.Method]; !ok {
				accesses[api.Url][route.Method] = make(map[string][]string)
			}

			pathRegex := regexp.MustCompile(`({[\w-]+})`).ReplaceAllString(route.Path, `[\w-]+`)
			accesses[api.Url][route.Method][pathRegex] = routeClients
		}
	}

	return accesses, clients, nil
}
