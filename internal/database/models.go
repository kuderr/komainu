// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package database

import (
	"github.com/google/uuid"
)

type Api struct {
	ID  uuid.UUID
	Url string
}

type Client struct {
	ID uuid.UUID
}

type Group struct {
	ID uuid.UUID
}

type GroupClientsAssociation struct {
	GroupID  uuid.UUID
	ClientID uuid.UUID
}

type GroupRoutesAssociation struct {
	GroupID uuid.UUID
	RouteID uuid.UUID
}

type Route struct {
	ID     uuid.UUID
	Method string
	Path   string
	ApiID  uuid.UUID
}

type RoutesAssociation struct {
	ClientID uuid.UUID
	RouteID  uuid.UUID
}
