CREATE TABLE routes_association (
    client_id uuid NOT NULL,
    route_id uuid NOT NULL
);

CREATE TABLE clients (
    id uuid PRIMARY KEY,
    name character varying(100) NOT NULL
);

CREATE TABLE routes (
    id uuid PRIMARY KEY,
    method character varying(100) NOT NULL,
    path character varying(100) NOT NULL,
    api_id uuid NOT NULL
);

CREATE TABLE apis (
    id uuid PRIMARY KEY,
    name character varying(100) NOT NULL,
    url character varying(100) NOT NULL
);