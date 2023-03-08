CREATE TABLE routes_association (
    client_id uuid NOT NULL,
    route_id uuid NOT NULL
);

CREATE TABLE clients (
    id uuid PRIMARY KEY
);

CREATE TABLE routes (
    id uuid PRIMARY KEY,
    method character varying(100) NOT NULL,
    path character varying(100) NOT NULL,
    api_id uuid NOT NULL
);

CREATE TABLE apis (
    id uuid PRIMARY KEY,
    url character varying(100) NOT NULL
);

CREATE TABLE group_clients_association (
    group_id uuid NOT NULL,
    client_id uuid NOT NULL
);

CREATE TABLE group_routes_association (
    group_id uuid NOT NULL,
    route_id uuid NOT NULL
);

CREATE TABLE groups (
    id uuid PRIMARY KEY
);