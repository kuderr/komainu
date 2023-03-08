package checker

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type AccessMap map[string]map[string]map[string][]uuid.UUID
type ClientsMap map[uuid.UUID][]uuid.UUID

type AuthInfo struct {
	// Map api url -> methods -> routes -> clients slice.
	// Slice of clients, because we will need to search for intersection later.
	// E.g. {"https://test.com": {"GET": {"/test": {"grp1", "user2"}, `/test/[\w-]+`: {"user1", "rob1"}}}}
	// Instead of client names is uuid`s.
	accesses AccessMap

	// Map of clients to its client `entities` – client groups and itself.
	// Through these entities accesses are created in management layer (auth-api).
	// And by them accesses are checked here.
	// E.g. {"user1": {"grp1", "grp2", "user1"}, "rob1": {"rob1", "grp2"}}
	// Instead of client names is uuid`s.
	clients ClientsMap

	// Lock update operation
	mutex sync.RWMutex
}

func NewAuthInfo(accesses AccessMap, clients ClientsMap) *AuthInfo {
	return &AuthInfo{accesses: accesses, clients: clients}
}

func (a *AuthInfo) Update(accesses AccessMap, clients ClientsMap) {
	// Lock read operations to update auth info
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.accesses = accesses
	a.clients = clients
}

func (a *AuthInfo) CheckAccess(request *AccessData, clientID uuid.UUID) (bool, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	clientEntities, ok := a.clients[clientID]
	if !ok {
		return false, &NotFoundError{msg: fmt.Sprintf("Client with ID %s not found", clientID)}
	}

	reqApiUrl := strings.TrimRight(request.ApiUrl, "/")
	reqPath := strings.TrimRight(request.Path, "/")
	reqMethod := strings.ToUpper(request.Method)

	api, ok := a.accesses[reqApiUrl]
	if !ok {
		return false, &NotFoundError{msg: fmt.Sprintf("Not found api with url %s", reqApiUrl)}
	}

	paths, ok := api[reqMethod]
	if !ok {
		return false, &NotFoundError{msg: fmt.Sprintf("Not found paths for %s %s", reqMethod, reqApiUrl)}
	}

	users, ok := paths[reqPath]
	if !ok {
		for p, u := range paths {
			re := regexp.MustCompile(fmt.Sprintf("^%s$", p))
			if re.MatchString(reqPath) {
				users = u
				break
			}
		}
	}

	if hasIntersection(users, clientEntities) {
		return true, nil
	}

	return false, nil
}

func hasIntersection[T comparable](a []T, b []T) bool {
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := hash[v]; ok {
			return true
		}
	}

	return false
}
