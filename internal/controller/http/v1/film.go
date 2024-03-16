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

	mux.HandleFunc("/api/v1/films/create", middleware.RequireAuth(ar.CreateFilm))
	mux.HandleFunc("/api/v1/films/name", middleware.RequireAuth(ar.GetFilmsByName))
	mux.HandleFunc("/api/v1/films/actor", middleware.RequireAuth(ar.GetFilmsByActor))
	mux.HandleFunc("/api/v1/films/delete", middleware.RequireAuth(ar.DeleteFilm))
}

func (fr *filmRoutes) CreateFilm(w http.ResponseWriter, req *http.Request) {
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
	jsonResp, err := json.Marshal(response{Id: id})
	if err != nil {
		fr.log.Errorf("filmRoutes CreateFilm: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (fr *filmRoutes) GetFilmsByName(w http.ResponseWriter, req *http.Request) {
	role := req.Header.Get(userRoleHeader)
	if role != "admin" && role != "user" {
		fr.log.Error("filmRoutes GetFilmsByName: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	type input struct {
		part string `json:"name_part"`
	}
	var name input
	if err := json.NewDecoder(req.Body).Decode(&name); err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByName: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	films, err := fr.filmService.GetFilmsByName(context.Background(), name.part)
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByName: filmService.GetFilmsByName %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(films)
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByName: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (fr *filmRoutes) GetFilmsByActor(w http.ResponseWriter, req *http.Request) {
	role := req.Header.Get(userRoleHeader)
	if role != "admin" && role != "user" {
		fr.log.Error("filmRoutes GetFilmsByActor: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	type input struct {
		part string `json:"name_part"`
	}
	var name input
	if err := json.NewDecoder(req.Body).Decode(&name); err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByActor: invalid request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	films, err := fr.filmService.GetFilmsByActor(context.Background(), name.part)
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByActor: filmService.GetFilmsByActor %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(films)
	if err != nil {
		fr.log.Errorf("filmRoutes GetFilmsByActor: cannot marshal response %v", err)
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (fr *filmRoutes) DeleteFilm(w http.ResponseWriter, req *http.Request) {
	role := req.Header.Get(userRoleHeader)
	if role != "admin" {
		fr.log.Error("filmRoutes DeleteFilm: user does not have the necessary rights")
		http.Error(w, "you do not have the necessary rights", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		fr.log.Errorf("filmRoutes DeleteFilm: cannot get film id %v", err)
		http.Error(w, "cannot get film id", http.StatusBadRequest)
		return
	}

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
