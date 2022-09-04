CREATE TABLE routes_association (
    client_id uuid PRIMARY KEY,
    route_id uuid PRIMARY KEY
);

CREATE TABLE admins_association (
    client_id uuid PRIMARY KEY,
    api_id uuid PRIMARY KEY
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