package blog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mikiejoe/go-blog-api/middlewares"
	"github.com/Mikiejoe/go-blog-api/types"
	"github.com/Mikiejoe/go-blog-api/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	store     types.BlogInterface
	userStore types.UserInTerface
}

func NewHandler(store types.BlogInterface, userStore types.UserInTerface) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/blogs", h.GetBlogsHandler).Methods(http.MethodGet)
	router.HandleFunc("/blogs", middlewares.AuthMiddleWare(h.CreateBlogHandler, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/blogs/{id}", h.GetBlogHandler).Methods(http.MethodGet)

}

func (h Handler) GetBlogsHandler(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.store.GetBlogs()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, blogs)
}

func (h Handler) CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.BlogPayload
	fmt.Println("inside the handler")
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	fmt.Println("past parse json")
	userId := middlewares.GetUseridFromCtx(r.Context())
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("something went wrong"))
		return
	}

	docId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	res, err := h.store.CreateBlog(types.Blog{
		Title:     payload.Title,
		Content:   payload.Content,
		UserId:    docId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"blogId": res,
	})

}
func (h *Handler) UpdateBlogHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.BlogPayload
	fmt.Println("inside the handler")
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id param is required"))
		return
	}
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userId := middlewares.GetUseridFromCtx(r.Context())
	docId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	blog, err := h.store.GetBlogByID(id)
	if err != mongo.ErrNoDocuments {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blog with id %v not found", id))

		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if blog.ID != docId {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you are not allowed to modify this blog"))
	}
	err = h.store.UpdateBlog(id, blog)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteSuccess(w, http.StatusNoContent, "blog updated successfully")
}

func (h *Handler) GetBlogHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println("id from path ", id)
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id param is required"))
		return
	}
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	blog, err := h.store.GetBlogByID(id)
	if err == mongo.ErrNoDocuments {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blog with id %v not found", id))

		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, blog)
}
