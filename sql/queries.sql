-- name: GetClients :many
SELECT id
FROM clients;

-- name: GetClientGroups :many
SELECT groups.id
FROM group_clients_association as gca
JOIN groups
  ON groups.id = gca.group_id
WHERE gca.client_id = $1;

-- name: GetApis :many
SELECT id, url
FROM apis;

-- name: GetApiRoutes :many
SELECT id, path, method
FROM routes 
WHERE api_id = $1;

-- name: GetRouteClients :many
SELECT clients.id 
FROM routes_association AS ra 
JOIN clients
  ON ra.client_id = clients.id 
WHERE ra.route_id = $1;

-- name: GetRouteGroups :many
SELECT groups.id
FROM group_routes_association AS gra
JOIN groups
  ON gra.group_id = groups.id
WHERE gra.route_id = $1;
