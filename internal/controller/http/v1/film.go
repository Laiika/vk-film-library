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

type filmRoutes struct {
	filmService service.Film
	log         *logger.Logger
}

func newFilmRoutes(mux *http.ServeMux, filmService service.Film, middleware *AuthMiddleware, log *logger.Logger) {
	ar := &filmRoutes{
		filmService: filmService,
		log:         log,
	}

	mux.HandleFunc("/api/v1/films/create", middleware.RequireAuth(ar.createFilm))
	mux.HandleFunc("/api/v1/films/name", middleware.RequireAuth(ar.getFilmsByName))
	mux.HandleFunc("/api/v1/films/actor", middleware.RequireAuth(ar.getFilmsByActor))
	mux.HandleFunc("/api/v1/films/delete/{id}", middleware.RequireAuth(ar.deleteFilm))
}

// @Summary Create film
// @Description Create film
// @Tags films
// @Param input body entity.FilmCreateInput true "information about stored film"
// @Accept json
// @Produce json
// @Success 201 {object} v1.filmRoutes.createFilm.response
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/films/create [post]
func (fr *filmRoutes) createFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

	role := req.Header.Get(userRoleHeader)
	if role != "admin" {
		fr.log.Error("filmRoutes CreateFilm: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	var input entity.FilmCreateInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		fr.log.Errorf("filmRoutes CreateFilm: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := fr.filmService.CreateFilm(context.Background(), &input)
	if err != nil {
		fr.log.Errorf("filmRoutes CreateFilm: filmService.CreateFilm %v", err)
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
		fr.log.Errorf("filmRoutes CreateFilm: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

// @Summary Get films by name
// @Description Get films by part of name
// @Tags films
// @Param input body entity.NamePart true "information about film name"
// @Accept json
// @Produce json
// @Success 200 {object} v1.filmRoutes.getFilmsByName.response
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/films/name [post]
func (fr *filmRoutes) getFilmsByName(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

	role := req.Header.Get(userRoleHeader)
	if role != "admin" && role != "user" {
		fr.log.Error("filmRoutes GetFilmsByName: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	var input entity.NamePart
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByName: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	films, err := fr.filmService.GetFilmsByName(context.Background(), input.Name)
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByName: filmService.GetFilmsByName %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Films []*entity.Film `json:"films"`
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response{Films: films})
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByName: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

// @Summary Get films by actor name
// @Description Get films by part of actor name
// @Tags films
// @Param input body entity.NamePart true "information about actor name"
// @Accept json
// @Produce json
// @Success 200 {object} v1.filmRoutes.getFilmsByActor.response
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/films/actor [post]
func (fr *filmRoutes) getFilmsByActor(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

	role := req.Header.Get(userRoleHeader)
	if role != "admin" && role != "user" {
		fr.log.Error("filmRoutes GetFilmsByActor: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	var input entity.NamePart
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByActor: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	films, err := fr.filmService.GetFilmsByActor(context.Background(), input.Name)
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByActor: filmService.GetFilmsByActor %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Films []*entity.Film `json:"films"`
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response{Films: films})
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByActor: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

// @Summary Delete film
// @Description Delete film
// @Tags films
// @Param id path integer true "Film id"
// @Success 200
// @Failure 400 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Security JWT
// @Router /api/v1/films/delete/{id} [delete]
func (fr *filmRoutes) deleteFilm(w http.ResponseWriter, req *http.Request) {
	if req.Method != "DELETE" {
		http.Error(w, "incorrect http method", http.StatusBadRequest)
		return
	}

	role := req.Header.Get(userRoleHeader)
	if role != "admin" {
		fr.log.Error("filmRoutes DeleteFilm: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		fr.log.Errorf("filmRoutes DeleteFilm: cannot get film id %v", err)
		http.Error(w, "cannot get film id", http.StatusBadRequest)
		return
	}
	fr.log.Println(id)

	err = fr.filmService.DeleteFilm(context.Background(), id)
	if err != nil {
		fr.log.Errorf("filmRoutes DeleteFilm: filmService.DeleteFilm %v", err)
		if err == service.ErrFilmNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
