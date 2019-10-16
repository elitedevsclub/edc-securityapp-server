package handlers

import (
	"edc-security-app/repos"
	"edc-security-app/types"
	"encoding/json"
	"net/http"
	"web"
)

type UserHttpHandler struct {
	repo *repos.UserRepository
}

func NewUserHttpHandler(repo *repos.UserRepository) *UserHttpHandler {
	return &UserHttpHandler{repo:repo}
}

func (handler *UserHttpHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body := &types.CreateUserOpts{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(body); err != nil {
		web.BadRequest(w, "bad request: malformed json data")
		return
	}
	user, otp, err := handler.repo.CreateUser(body)
	if err != nil {
		web.BadRequest(w, err.Error())
		return
	}
	type response struct {
		User *types.User `json:"user"`
		CodeIdentifier string `json:"code_identifier"`
	}
	web.OK(w, &response{User: user, CodeIdentifier: otp.CodeIdentifier}, "application/json")
}
