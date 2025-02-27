package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (a Auth) apiRegister(w http.ResponseWriter, r *http.Request) {
	if a.passwordless {
		a.apiRegisterPasswordless(w, r)
	} else {
		a.apiRegisterUsernameAndPassword(w, r)
	}
}

func (a Auth) apiRegisterPasswordless(w http.ResponseWriter, r *http.Request) {
	email := strings.Trim(utils.Req(r, "email", ""), " ")
	first_name := strings.Trim(utils.Req(r, "first_name", ""), " ")
	last_name := strings.Trim(utils.Req(r, "last_name", ""), " ")

	if first_name == "" {
		api.Respond(w, r, api.Error("First name is required field"))
		return
	}

	if last_name == "" {
		api.Respond(w, r, api.Error("Last name is required field"))
		return
	}

	if email == "" {
		api.Respond(w, r, api.Error("Email is required field"))
		return
	}

	verificationCode := utils.StrRandomFromGamma(LoginCodeLength, LoginCodeGamma)

	json, errJson := utils.ToJSON(map[string]string{
		"email":      email,
		"first_name": first_name,
		"last_name":  last_name,
	})

	if errJson != nil {
		api.Respond(w, r, api.Error("Error serializing data"))
		return
	}

	errTempTokenSave := a.funcTemporaryKeySet(verificationCode, json, 3600)

	if errTempTokenSave != nil {
		api.Respond(w, r, api.Error("token store failed. "+errTempTokenSave.Error()))
		return
	}

	emailContent := a.passwordlessFuncEmailTemplateRegisterCode(email, verificationCode)

	errEmailSent := a.passwordlessFuncEmailSend(email, "Registration Code", emailContent)

	if errEmailSent != nil {
		log.Println(errEmailSent)
		api.Respond(w, r, api.Error("Registration code failed to be send. Please try again later"))
		return
	}

	api.Respond(w, r, api.Success("Registration code was sent successfully"))
}

func (a Auth) apiRegisterUsernameAndPassword(w http.ResponseWriter, r *http.Request) {
	email := strings.Trim(utils.Req(r, "email", ""), " ")
	password := strings.Trim(utils.Req(r, "password", ""), " ")
	first_name := strings.Trim(utils.Req(r, "first_name", ""), " ")
	last_name := strings.Trim(utils.Req(r, "last_name", ""), " ")

	response := a.RegisterWithUsernameAndPassword(email, password, first_name, last_name)

	if response.ErrorMessage != "" {
		api.Respond(w, r, api.Error(response.ErrorMessage))
		return
	}

	api.Respond(w, r, api.Success(response.SuccessMessage))
}
