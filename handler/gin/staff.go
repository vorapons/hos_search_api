package ginhandler

import (
	"errors"
	"net/http"
	"time"

	"pt_search_hos/domain"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	service domain.StaffService
}

func NewStaffHandler(s domain.StaffService) *StaffHandler {
	return &StaffHandler{service: s}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type createStaffRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

// Login handles POST /staff/login
func (h *StaffHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_INPUT",
			"message": "Invalid request body",
		})
		return
	}

	token, err := h.service.Login(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_INPUT",
				"message": "Login and password are required",
			})
			return
		}
		if errors.Is(err, domain.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid login or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Create handles POST /staff/create
func (h *StaffHandler) Create(c *gin.Context) {
	var req createStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_INPUT",
			"message": "Invalid request body",
		})
		return
	}

	token, err := h.service.CreateStaff(req.Login, req.Password, req.Hospital)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_INPUT",
				"message": "Login, password, and hospital are required",
			})
		case errors.Is(err, domain.ErrStaffExists):
			c.JSON(http.StatusConflict, gin.H{
				"code":    "CONFLICT",
				"message": "A staff account with this login already exists",
			})
		case errors.Is(err, domain.ErrHospitalNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "Hospital not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// Hello handles GET /staff/hello — returns JWT claims for the current user.
func (h *StaffHandler) Hello(c *gin.Context) {
	login, _    := c.Get("login")
	hospital, _ := c.Get("hospital")
	exp, _      := c.Get("exp")

	c.JSON(http.StatusOK, gin.H{
		"login":      login.(string),
		"hospital":   hospital.(string),
		"expires_at": exp.(time.Time).UTC().Format(time.RFC3339),
	})
}

// Logout handles GET /staff/logout
func (h *StaffHandler) Logout(c *gin.Context) {
	tokenRaw, _ := c.Get("token")
	tokenString, _ := tokenRaw.(string)

	if err := h.service.Logout(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
