package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter(env *Env) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/auth/register", env.UserRegisterPost).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/auth/login", env.UserSignInPost).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/user/self", env.UserGetSelfInfo).Methods(http.MethodGet)

	return router
}
