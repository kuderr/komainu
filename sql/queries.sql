-- name: GetClients :many
SELECT name
FROM clients;

-- name: GetApis :many
SELECT id, url
FROM apis;

-- name: GetApiRoutes :many
SELECT id, path, method
FROM routes 
WHERE api_id = $1;

-- name: GetRouteClients :many
SELECT c.name 
FROM routes_association AS ra 
JOIN clients AS c ON ra.client_id=c.id 
WHERE route_id = $1;