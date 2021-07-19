package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/IvanKyrylov/cost-management-api/internal/apperror"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/user/"
)

type Handler struct {
	Logger      *log.Logger
	UserService Service
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(userURL, apperror.Middleware(h.GetUser))
	router.HandleFunc(usersURL, apperror.Middleware(h.GetAllUsers))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}
	h.Logger.Println("GET USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get uuid from path")
	userId, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		return apperror.BadRequestError("id query parameter is required positive integers")
	}

	user, err := h.UserService.GetById(r.Context(), userId)
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}

	h.Logger.Println("GET USERS")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get limit from URL")
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		return apperror.BadRequestError("limit query parameter is required positive integers")
	}

	h.Logger.Println("get limit from URL")
	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		return apperror.BadRequestError("page query parameter is required positive integers")
	}

	users, err := h.UserService.GetAll(r.Context(), limit, page)
	if err != nil {
		return err
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(usersBytes)
	return nil
}
