package handler

import (
	"errors"
	"time"
	"pt_search_hos/domain"

	"github.com/gofiber/fiber/v2"
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
func (h *StaffHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_INPUT",
			"message": "Invalid request body",
		})
	}

	token, err := h.service.Login(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "Login and password are required",
			})
		}
		if errors.Is(err, domain.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "Invalid login or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		})
	}

	return c.JSON(fiber.Map{"token": token})
}

// Create handles POST /staff/create
func (h *StaffHandler) Create(c *fiber.Ctx) error {
	var req createStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_INPUT",
			"message": "Invalid request body",
		})
	}

	token, err := h.service.CreateStaff(req.Login, req.Password, req.Hospital)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "Login, password, and hospital are required",
			})
		case errors.Is(err, domain.ErrStaffExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"code":    "CONFLICT",
				"message": "A staff account with this login already exists",
			})
		case errors.Is(err, domain.ErrHospitalNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Hospital not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
}

// Hello handles GET /staff/hello — returns JWT claims for the current user.
func (h *StaffHandler) Hello(c *fiber.Ctx) error {
	login, _ := c.Locals("login").(string)
	hospital, _ := c.Locals("hospital").(string)
	exp, _ := c.Locals("exp").(time.Time)

	return c.JSON(fiber.Map{
		"login":      login,
		"hospital":   hospital,
		"expires_at": exp.UTC().Format(time.RFC3339),
	})
}

// Logout handles GET /staff/logout
func (h *StaffHandler) Logout(c *fiber.Ctx) error {
	tokenString, _ := c.Locals("token").(string)

	if err := h.service.Logout(tokenString); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		})
	}

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}
