package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/service"
	"vk-film-library/pkg/logger"
)

type actorRoutes struct {
	actorService service.Actor
	log          *logger.Logger
}

func newActorRoutes(mux *http.ServeMux, actorService service.Actor, middleware *AuthMiddleware, log *logger.Logger) {
	ar := &actorRoutes{
		actorService: actorService,
		log:          log,
	}

	mux.HandleFunc("/api/v1/actors/create", middleware.RequireAuth(ar.CreateActor))
	mux.HandleFunc("/api/v1/actors", middleware.RequireAuth(ar.GetAllActors))
	mux.HandleFunc("/api/v1/actors/edit", middleware.RequireAuth(ar.EditActor))
	mux.HandleFunc("/api/v1/actors/delete", middleware.RequireAuth(ar.DeleteActor))
}

func (ar *actorRoutes) CreateActor(w http.ResponseWriter, req *http.Request) {
	role := req.Header.Get(userRoleHeader)
	if role != "admin" {
		ar.log.Error("actorRoutes CreateActor: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	var input entity.ActorCreateInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		ar.log.Errorf("actorRoutes CreateActor: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := ar.actorService.CreateActor(context.Background(), &input)
	if err != nil {
		ar.log.Errorf("actorRoutes CreateActor: actorService.CreateActor %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	w.WriteHeader(http.StatusCreated)
	jsonResp, err := json.Marshal(response{Id: id})
	if err != nil {
		ar.log.Errorf("actorRoutes CreateActor: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (ar *actorRoutes) GetAllActors(w http.ResponseWriter, req *http.Request) {
	role := req.Header.Get(userRoleHeader)
	if role != "admin" && role != "user" {
		ar.log.Error("actorRoutes GetAllActors: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	actors, err := ar.actorService.GetAllActors(context.Background())
	if err != nil {
		ar.log.Errorf("actorRoutes GetAllActors: actorService.GetAllActors %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(actors)
	if err != nil {
		ar.log.Errorf("actorRoutes GetAllActors: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (ar *actorRoutes) EditActor(w http.ResponseWriter, req *http.Request) {
	role := req.Header.Get(userRoleHeader)
	if role != "admin" {
		ar.log.Error("actorRoutes EditActor: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	var input entity.Actor
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		ar.log.Errorf("actorRoutes EditActor: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		ar.log.Errorf("actorRoutes EditActor: cannot get actor id %v", err)
		http.Error(w, "cannot get actor id", http.StatusBadRequest)
		return
	}
	input.Id = id
	err = ar.actorService.EditActor(context.Background(), &input)
	if err != nil {
		ar.log.Errorf("actorRoutes EditActor: actorService.EditActor %v", err)
		if err == service.ErrActorNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ar *actorRoutes) DeleteActor(w http.ResponseWriter, req *http.Request) {
	role := req.Header.Get(userRoleHeader)
	if role != "admin" {
		ar.log.Error("actorRoutes DeleteActor: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		ar.log.Errorf("actorRoutes DeleteActor: cannot get actor id %v", err)
		http.Error(w, "cannot get actor id", http.StatusBadRequest)
		return
	}

	err = ar.actorService.DeleteActor(context.Background(), id)
	if err != nil {
		ar.log.Errorf("actorRoutes DeleteActor: actorService.DeleteActor %v", err)
		if err == service.ErrActorNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
