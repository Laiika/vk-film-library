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

	mux.HandleFunc("/api/v1/actors/create", middleware.RequireAuth(ar.createActor))
	mux.HandleFunc("/api/v1/actors", middleware.RequireAuth(ar.getAllActors))
	mux.HandleFunc("/api/v1/actors/edit", middleware.RequireAuth(ar.editActor))
	mux.HandleFunc("/api/v1/actors/delete/{id}", middleware.RequireAuth(ar.deleteActor))
}

// @Summary Create actor
// @Description Create actor
// @Tags actors
// @Param input body entity.ActorCreateInput true "information about stored actor"
// @Accept json
// @Produce json
// @Success 201 {object} v1.actorRoutes.createActor.response
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/actors/create [post]
func (ar *actorRoutes) createActor(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

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
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response{Id: id})
	if err != nil {
		ar.log.Errorf("actorRoutes CreateActor: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

// @Summary Get all actors
// @Description Get all actors
// @Tags actors
// @Produce json
// @Success 200 {object} v1.actorRoutes.getAllActors.response
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/actors [get]
func (ar *actorRoutes) getAllActors(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

	role := req.Header.Get(userRoleHeader)
	if role != "admin" && role != "user" {
		ar.log.Errorf("actorRoutes GetAllActors: user does not have the necessary rights %s", role)
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	actors, err := ar.actorService.GetAllActors(context.Background())
	if err != nil {
		ar.log.Errorf("actorRoutes GetAllActors: actorService.GetAllActors %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Actors []*entity.Actor `json:"actors"`
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response{Actors: actors})
	if err != nil {
		ar.log.Errorf("actorRoutes GetAllActors: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

// @Summary Edit actor
// @Description Edit actor
// @Tags actors
// @Param input body entity.Actor true "information about actor"
// @Accept json
// @Success 200
// @Failure 400 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/actors/edit [put]
func (ar *actorRoutes) editActor(w http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

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

	err := ar.actorService.EditActor(context.Background(), &input)
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

// @Summary Delete actor
// @Description Delete actor
// @Tags actors
// @Param id path integer true "Actor id"
// @Success 200
// @Failure 400 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/actors/delete/{id} [delete]
func (ar *actorRoutes) deleteActor(w http.ResponseWriter, req *http.Request) {
	if req.Method != "DELETE" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

	role := req.Header.Get(userRoleHeader)
	if role != "admin" {
		ar.log.Error("actorRoutes DeleteActor: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(req.PathValue("id"))
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
