package main

import (
	"errors"
	"log"

	"github.com/gouniverse/utils"
)

func passwordlessUserRegister(username string, first_name string, last_name string) error {
	slug := utils.StrSlugify(username, rune('_'))
	err := jsonStore.Write("users", slug, map[string]string{
		"id":         utils.StrRandomFromGamma(16, "abcdef0123456789"),
		"username":   username,
		"password":   "passwordless_registered", // no need for password
		"first_name": first_name,
		"last_name":  last_name,
	})
	if err != nil {
		return err
	}
	return nil
}

func passwordlessUserFindByEmail(email string) (userID string, err error) {
	slug := utils.StrSlugify(email, rune('_'))
	var user map[string]string
	err = jsonStore.Read("users", slug, &user)
	if err != nil {
		log.Println(err.Error())
		return "not found err", errors.New("unable to find user")
	}

	if user == nil {
		return "not found", errors.New("unable to find user")
	}

	return user["id"], nil
}
