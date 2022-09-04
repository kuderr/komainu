-- name: GetClientIdByName :one
SELECT id
FROM clients 
WHERE name = $1;

-- name: GetApiIdByUrl :one
SELECT id 
FROM apis 
WHERE url = $1;

-- name: CountAdminAssociation :one
SELECT count(*)
FROM admins_association 
WHERE client_id = $1 AND api_id = $2;

-- name: GetApiRouteIdByMethodAndPath :one
SELECT id 
FROM routes 
WHERE api_id = $1 AND method = $2 AND path = $3;

-- name: GetApiRoutesByMethod :many
SELECT * 
FROM routes 
WHERE api_id = $1 AND method = $2;

-- name: CountRouteAssociation :one
SELECT count(*)
FROM routes_association 
WHERE client_id = $1 AND route_id = $2;