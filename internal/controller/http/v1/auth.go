package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/service"
	"vk-film-library/pkg/logger"
)

type authRoutes struct {
	authService service.Auth
	log         *logger.Logger
}

func newAuthRoutes(mux *http.ServeMux, authService service.Auth, log *logger.Logger) {
	ar := &authRoutes{
		authService: authService,
		log:         log,
	}

	mux.HandleFunc("/sign-in", ar.signIn)
	mux.HandleFunc("/sign-up", ar.signUpUser)
	mux.HandleFunc("/admin/sign-up", ar.signUpAdmin)
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password" validate:"required,password"`
}

func (ar *authRoutes) signUpUser(w http.ResponseWriter, req *http.Request) {
	var input signInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		ar.log.Errorf("authRoutes signUpUser: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := ar.authService.CreateUser(context.Background(), &entity.CreateInput{
		Username: input.Username,
		Password: input.Password,
		Role:     "user",
	})
	if err != nil {
		ar.log.Errorf("authRoutes signUpUser: authService.CreateUser %v", err)
		if err == service.ErrUserAlreadyExists {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	w.WriteHeader(http.StatusCreated)
	jsonResp, err := json.Marshal(response{Id: id})
	if err != nil {
		ar.log.Errorf("authRoutes signUpUser: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (ar *authRoutes) signUpAdmin(w http.ResponseWriter, req *http.Request) {
	var input signInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		ar.log.Errorf("authRoutes signUpAdmin: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := ar.authService.CreateUser(context.Background(), &entity.CreateInput{
		Username: input.Username,
		Password: input.Password,
		Role:     "admin",
	})
	if err != nil {
		ar.log.Errorf("authRoutes signUpAdmin: authService.CreateUser %v", err)
		if err == service.ErrUserAlreadyExists {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	w.WriteHeader(http.StatusCreated)
	jsonResp, err := json.Marshal(response{Id: id})
	if err != nil {
		ar.log.Errorf("authRoutes signUpAdmin: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (ar *authRoutes) signIn(w http.ResponseWriter, req *http.Request) {
	var input signInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		ar.log.Errorf("authRoutes signUpIn: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := ar.authService.GenerateToken(context.Background(), &entity.AuthInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		ar.log.Errorf("authRoutes signIn: authService.GenerateToken %v", err)
		if err == service.ErrUserNotFound {
			http.Error(w, "invalid username or password", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(response{Token: token})
	if err != nil {
		ar.log.Errorf("authRoutes signIn: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
