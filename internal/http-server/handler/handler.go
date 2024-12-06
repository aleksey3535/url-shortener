package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"url-shortener/internal/http-server/middleware"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/gorilla/mux"
)

type Handler struct {
	log 	*slog.Logger
	mw 		*middleware.Middleware
	st		storage.StorageI
}

func New(log *slog.Logger, mw *middleware.Middleware, st storage.StorageI) *Handler {
	return &Handler{log: log, mw: mw, st: st}
}

func(h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/short/create", h.mw.Logger(h.CreateHandler)).Methods(http.MethodPost)
	h.log.Debug("URL:'/short/create' METHOD:'POST' INFO:'CREATE NEW SHORT URL' STATUS: ACTIVE")
	router.HandleFunc("/short/delete", h.mw.Logger(h.DeleteHandler)).Methods(http.MethodPost)
	h.log.Debug("URL:'/`short/delete' METHOD:'POST' INFO: 'DELETE SHORT URL' STATUS: ACTIVE")
	router.HandleFunc("/short/{alias}", h.mw.Logger(h.ShortHandler)).Methods(http.MethodGet)
	h.log.Debug("URL:'/short/{alias}' METHOD:'GET' INFO: 'REDIRECT USING SHORT URL' STATUS: ACTIVE")
	return router
}

type CreateInput struct {
	Url		string `json:"url"`
	Alias	string `json:"alias"`

}
type Info struct {
	Message	string `json:"msg"`
}

func(h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateHandler"
	log := h.log.With("op", op)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("failed to read body", sl.Err(err))
		http.Error(w, "Something wrong", http.StatusInternalServerError)
		return
	}
	var input CreateInput
	if err = json.Unmarshal(body, &input); err != nil {
		log.Error("failed to unmarshal body", sl.Err(err))
		http.Error(w, "invalid body syntax", http.StatusBadRequest)
		return
	}
	if err = h.st.SaveURL(input.Url, input.Alias); err != nil {
		log.Error("failed to save url", sl.Err(err))
		http.Error(w, "url already exists", http.StatusConflict)
		return
	}
	mes := Info{Message: "complete"}
	info, _ := json.Marshal(mes)
	w.Write(info)
}

type DeleteInput struct {
	Alias string `json:"alias"`
}

func(h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteHandler"
	log := h.log.With("op", op)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("failed to read body", sl.Err(err))
		http.Error(w, "Something wrong", http.StatusInternalServerError)
		return
	}
	var input DeleteInput
	if err = json.Unmarshal(body, &input); err != nil {
		log.Error("failed to unmarshal body", sl.Err(err))
		http.Error(w, "invalid body syntax", http.StatusBadRequest)
		return
	}
	if err = h.st.DeleteURL(input.Alias); err != nil {
		log.Error("failed to delete url", sl.Err(err))
		http.Error(w, "such url does not exists", http.StatusConflict)
		return
	}
	mes := Info{Message: "complete"}
	info, _ := json.Marshal(mes)
	w.Write(info)
}

func(h *Handler) ShortHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.ShortHandler"
	log := h.log.With("op", op)
	alias := mux.Vars(r)["alias"]
	if alias == "" {
		http.Error(w, "{alias} is required", http.StatusBadRequest)
	}
	urlToRedirect, err := h.st.GetURL(alias)
	if err != nil {
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Error("entered wrong alias", sl.Err(err))
			http.Error(w, "such url does not exists", http.StatusBadRequest)
			return
		}
		log.Error("failed to get url", sl.Err(err))
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, urlToRedirect, http.StatusFound)


}



