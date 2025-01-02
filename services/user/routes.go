package user

import (
	"fmt"
	"net/http"

	"github.com/Mikiejoe/go-blog-api/config"
	"github.com/Mikiejoe/go-blog-api/services/auth"
	"github.com/Mikiejoe/go-blog-api/types"
	"github.com/Mikiejoe/go-blog-api/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	userInteface types.UserInTerface
}

func NewHandler(userInteface types.UserInTerface) *Handler {
	return &Handler{
		userInteface: userInteface,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.RegisterHandler).Methods(http.MethodPost)
	router.HandleFunc("/users", h.GetUsersHandler).Methods(http.MethodGet)
	router.HandleFunc("/login", h.LoginHandler).Methods(http.MethodPost)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("getting ser by username")
	u, err := h.userInteface.GetUserByName(payload.Username)
	fmt.Println("user found", u)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("username %s already taken", payload.Username))
		return
	}
	if err != mongo.ErrNoDocuments {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	fmt.Println("getting ser by email")
	_, err = h.userInteface.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("email %s already taken", payload.Email))
		return
	}
	if err != mongo.ErrNoDocuments {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	password, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("couldnt encrypt password"))
		return
	}

	id, err := h.userInteface.CreateUser(
		types.User{
			Firstname: payload.Firstname,
			Lastname:  payload.Lastname,
			Email:     payload.Email,
			Password:  password,
			Location:  payload.Location,
			Username:  payload.Username,
		},
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{"id": id, "password": password})
}

func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.userInteface.GetUsers()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.userInteface.GetUserByName(payload.Username)
	if err == mongo.ErrNoDocuments {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("incorrect username or password"))
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	passmatch := auth.ConparePassword(payload.Password, user.Password)
	if !passmatch {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("incorrect username or password"))
		return
	}
	secret:=[]byte(config.Envs.JWTSecret)
	token,err:= auth.CreateJWT(secret,user.ID.Hex())
	if err!=nil{
		utils.WriteError(w,http.StatusInternalServerError,err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"token":token,
	})

}
