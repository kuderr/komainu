package auther

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type AccessMap map[string]map[string]map[string][]string
type ClientsSet map[string]struct{}

type AuthInfo struct {
	// Map api -> methods -> routes -> clients.
	// Slice of clients, because we will need to search for intersection later.
	// E.g. {"https://test.com": {"GET": {"/test": ["rob1", "grp1"], `/test/[\w-]+`: ["user1", "rob1"]}}}
	accesses AccessMap

	// Set of all clients/robots in system
	clients ClientsSet

	// Lock update operation
	mutex sync.RWMutex
}

func NewAuthInfo(accesses AccessMap, clients ClientsSet) *AuthInfo {
	return &AuthInfo{accesses: accesses, clients: clients}
}

func (a *AuthInfo) Update(accesses AccessMap, clients ClientsSet) {
	// Lock read operations to update auth info
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.accesses = accesses
	a.clients = clients
}

func (a *AuthInfo) CheckAccess(request *AccessData, clientName string) (bool, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if _, ok := a.clients[clientName]; !ok {
		return false, &NotFoundError{msg: fmt.Sprintf("Client %s not found", clientName)}
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

	if contains(users, clientName) {
		return true, nil
	}

	return false, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// FUTURE
// Hash has complexity: O(n * x) where x is a factor of hash function efficiency (between 1 and 2)
func HashGeneric[T comparable](a []T, b []T) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}
