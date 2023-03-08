package checker

import (
	"checker/internal/database"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type DatabaseAuthInfoStorage struct {
	queries *database.Queries
}

func NewDatabaseAuthInfoStorage(queries *database.Queries) *DatabaseAuthInfoStorage {
	return &DatabaseAuthInfoStorage{queries: queries}
}

func (db *DatabaseAuthInfoStorage) GetClients() ([]uuid.UUID, error) {
	clients, err := db.queries.GetClients(context.Background())
	if err != nil {
		if err == pgx.ErrNoRows {
			return []uuid.UUID{}, nil
		}

		return nil, err
	}

	return clients, nil
}
func (db *DatabaseAuthInfoStorage) GetClientGroups(clientID uuid.UUID) ([]uuid.UUID, error) {
	groups, err := db.queries.GetClientGroups(context.Background(), clientID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []uuid.UUID{}, nil
		}

		return nil, err
	}

	return groups, nil
}

func (db *DatabaseAuthInfoStorage) GetApis() ([]Api, error) {
	apis, err := db.queries.GetApis(context.Background())
	if err != nil {
		if err == pgx.ErrNoRows {
			return []Api{}, nil
		}

		return nil, err
	}

	var items []Api
	for _, api := range apis {
		items = append(items, Api{ID: api.ID, Url: api.Url})
	}

	return items, nil
}

func (db *DatabaseAuthInfoStorage) GetApiRoutes(apiID uuid.UUID) ([]Route, error) {
	routes, err := db.queries.GetApiRoutes(context.Background(), apiID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []Route{}, nil
		}

		return nil, err
	}

	var items []Route
	for _, route := range routes {
		items = append(items, Route{ID: route.ID, Method: route.Method, Path: route.Path})
	}

	return items, nil
}

func (db *DatabaseAuthInfoStorage) GetRouteClients(routeID uuid.UUID) ([]uuid.UUID, error) {
	clients, err := db.queries.GetRouteClients(context.Background(), routeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []uuid.UUID{}, nil
		}

		return nil, err
	}

	return clients, nil
}

func (db *DatabaseAuthInfoStorage) GetRouteGroups(routeID uuid.UUID) ([]uuid.UUID, error) {
	groups, err := db.queries.GetRouteGroups(context.Background(), routeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []uuid.UUID{}, nil
		}

		return nil, err
	}

	return groups, nil
}
