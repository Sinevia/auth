package auth

import (
	"log"

	"github.com/gouniverse/utils"
	validator "github.com/gouniverse/validator"
)

type RegisterUsernameAndPasswordResponse struct {
	ErrorMessage   string
	SuccessMessage string
	Token          string
}

func (a Auth) RegisterWithUsernameAndPassword(email string, password string, firstName string, lastName string) (response RegisterUsernameAndPasswordResponse) {
	if firstName == "" {
		response.ErrorMessage = "First name is required field"
		return response
	}

	if lastName == "" {
		response.ErrorMessage = "Last name is required field"
		return response
	}

	if email == "" {
		response.ErrorMessage = "Email is required field"
		return response
	}

	if password == "" {
		response.ErrorMessage = "Password is required field"
		return response
	}

	if !validator.IsEmail(email) {
		response.ErrorMessage = "This is not a valid email: " + email
		return response
	}

	if a.funcUserRegister == nil {
		response.ErrorMessage = "registration failed. FuncUserRegister function not defined"
		return response
	}

	if !a.enableVerification {
		err := a.funcUserRegister(email, password, firstName, lastName)

		if err != nil {
			response.ErrorMessage = "registration failed. " + err.Error()
			return response
		}

		response.SuccessMessage = "registration success"
		return response
	}

	verificationCode := utils.StrRandomFromGamma(LoginCodeLength, LoginCodeGamma)

	json, errJson := utils.ToJSON(map[string]string{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"password":   password,
	})

	if errJson != nil {
		response.ErrorMessage = "Error serializing data"
		return response
	}

	errTempTokenSave := a.funcTemporaryKeySet(verificationCode, json, 3600)

	if errTempTokenSave != nil {
		response.ErrorMessage = "token store failed. " + errTempTokenSave.Error()
		return response
	}

	emailContent := a.funcEmailTemplateRegisterCode(email, verificationCode)

	errEmailSent := a.funcEmailSend(email, "Registration Code", emailContent)

	if errEmailSent != nil {
		log.Println(errEmailSent)
		response.ErrorMessage = "Registration code failed to be send. Please try again later"
		return response
	}

	response.SuccessMessage = "Registration code was sent successfully"
	return response
}
