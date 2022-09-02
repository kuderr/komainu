package auther

import (
	"auther/responses"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CheckAccess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var authData AuthData
	err := json.NewDecoder(r.Body).Decode(&authData)
	if err != nil {
		log.Println(err)
		resp := responses.ErrorResponse{Err: "Invalid auth data"}
		responses.SendError(w, resp, 400)
		return
	}

	hasAccess, err := checkAccess(&authData)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if hasAccess {
		err = responses.SendData(w, "OK", 200)
	} else {
		err = responses.SendData(w, "Access Denied", 403)
	}

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

}
