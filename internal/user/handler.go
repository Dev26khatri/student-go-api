package user

import (
	"errors"
	"net/http"
	"student-go-service/internal/dto"
	"student-go-service/internal/package/utils"
	"student-go-service/internal/response"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service       *Service
	jwtSeceret    string
	jwtExpiration time.Duration
}

func NewHandler(service *Service, jwtSecreet string, jwtExpiration time.Duration) *Handler {
	return &Handler{
		service:       service,
		jwtSeceret:    jwtSecreet,
		jwtExpiration: jwtExpiration,
	}
}
func toUserResponse(user *User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid request body", "Fix your Inputs")
		return
	}

	user, err := h.service.Register(req)

	if err != nil {
		if errors.Is(err, utils.ErrEmailAlreadyExists) {
			response.Error(
				c,
				http.StatusConflict,
				"Email already exists", "Using Another Work Email")
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to register user", err.Error())
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"User registered successfully",
		toUserResponse(user),
	)
}

func (h *Handler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Request Body", err.Error())

		return
	}
	user, token, err := h.service.Login(req, h.jwtSeceret, h.jwtExpiration)

	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Invalid email or password ", err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "Login Succesfully", dto.LoginResponse{Token: token, User: toUserResponse(user)})
}
