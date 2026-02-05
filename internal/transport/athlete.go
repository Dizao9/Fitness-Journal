package transport

import (
	"net/http"

	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
)

type AuthService interface {
	Register(req dto.RegisterUser) (string, error)
}

type Handler struct {
	AuthSvc AuthService
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

}
