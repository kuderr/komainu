package auther

import (
	"auther/config"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type AuthData struct {
	Token  string `json:"token"`
	ApiUrl string `json:"api_url"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type TokenBody struct {
	ClientName string `json:"client_name"`
}

type DecodedToken struct {
	Subject TokenBody `json:"sub"`
	jwt.StandardClaims
}

func checkAccess(authData *AuthData) (bool, error) {
	token, err := decodeToken(authData.Token)
	if err != nil {
		return false, err
	}

	var clientId string
	clientId, err = getClientIDByName(token.Subject.ClientName)
	if err != nil {
		return false, err
	}

	var apiId string
	apiId, err = getApiIDByUrl(authData.ApiUrl)
	if err != nil {
		return false, err
	}

	var isApiAdmin bool
	isApiAdmin, err = isApiAdminAssociationExist(clientId, apiId)
	if err != nil {
		return false, err
	}
	if isApiAdmin {
		return true, nil
	}

	var routeId string
	routeId, err = getApiRouteIDByMethodAndPath(apiId, authData.Method, authData.Path)
	if err != nil {
		return false, err
	}

	var hasRouteAccess bool
	hasRouteAccess, err = isRouteAccessExist(clientId, routeId)
	if err != nil {
		return false, err
	}
	if hasRouteAccess {
		return true, nil
	}

	return false, nil
}

func getClientIDByName(clientName string) (string, error) {
	var clientId string
	err := config.DB.QueryRow("SELECT id FROM clients WHERE name = $1", clientName).Scan(&clientId)
	if err != nil {
		return "", err
	}

	return clientId, nil
}

func getApiIDByUrl(apiUrl string) (string, error) {
	var apiId string
	err := config.DB.QueryRow("SELECT id FROM apis WHERE url = $1", apiUrl).Scan(&apiId)
	if err != nil {
		return "", err
	}

	return apiId, nil
}

func isApiAdminAssociationExist(clientId string, apiId string) (bool, error) {
	var count int
	err := config.DB.QueryRow("SELECT count(*) FROM admins_association WHERE client_id = $1 AND api_id = $2", clientId, apiId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count >= 1, nil
}

func getApiRouteIDByMethodAndPath(apiId string, method string, path string) (string, error) {
	var routeId string
	err := config.DB.QueryRow("SELECT id FROM routes WHERE api_id = $1 AND method = $2 AND path = $3", apiId, method, path).Scan(&routeId)
	if err != nil {
		return "", err
	}

	// if err != nil && err != sql.ErrNoRows {
	// 	return "", err
	// }

	// routeId, err = getPatternRoute(apiId, method, path)
	// if err != nil {
	// 	return "", err
	// }

	return routeId, nil
}

func isRouteAccessExist(clientId string, routeId string) (bool, error) {
	var count int
	err := config.DB.QueryRow("SELECT count(*) FROM routes_association WHERE client_id = $1 AND route_id = $2", clientId, routeId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count >= 1, nil
}

//  TODO !!!
// func getPatternRoute(apiId string, method string, path string) (string, error) {
// 	rows, err := config.DB.Query("SELECT id, path FROM routes WHERE api_id = $1 AND method = $2", apiId, method)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		route := Route{}
// 		re := regexp.MustCompile(`"({[\\\w-]+})"`)

// 		routePath := re.ReplaceAll([]byte())
// 		err := rows.Scan(&route.ID, &route.Path)
// 		if err != nil {
// 			return "", err
// 		}

// 	}

// 	if err = rows.Err(); err != nil {
// 		return "", err
// 	}
// }

func decodeToken(token string) (DecodedToken, error) {
	creds := strings.Replace(token, "Bearer ", "", 1)

	tk := &DecodedToken{}
	decoded, err := jwt.ParseWithClaims(creds, tk, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("APP_SECRET")), nil
	})

	if err != nil {
		log.Println(err)
		return *tk, err
	}

	if !decoded.Valid {
		log.Println("Invalid token")
		return *tk, errors.New("Invalid Token")
	}

	return *tk, nil
}
