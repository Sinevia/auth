package auth

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (a Auth) apiPasswordReset(w http.ResponseWriter, r *http.Request) {
	token := strings.Trim(utils.Req(r, "token", ""), " ")
	password := strings.Trim(utils.Req(r, "password", ""), " ")
	passwordConfirm := strings.Trim(utils.Req(r, "password_confirm", ""), " ")

	if token == "" {
		api.Respond(w, r, api.Error("Token is required field"))
		return
	}

	if password == "" {
		api.Respond(w, r, api.Error("Password is required field"))
		return
	}

	if password != passwordConfirm {
		api.Respond(w, r, api.Error("Passwords do not match"))
		return
	}

	userID, errToken := a.funcTemporaryKeyGet(token)

	if errToken != nil {
		api.Respond(w, r, api.Error("Link not valid of expired"))
		return
	}

	if userID == "" {
		api.Respond(w, r, api.Error("Link not valid of expired"))
		return
	}

	errPasswordChange := a.funcUserPasswordChange(userID, password)

	if errPasswordChange != nil {
		api.Respond(w, r, api.Error("authentication failed. "+errPasswordChange.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("login success", map[string]interface{}{
		"token": token,
	}))
}
