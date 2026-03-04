package handler

import (
	"errors"

	"pt_search_hos/domain"

	"github.com/gofiber/fiber/v2"
)

type PatientHandler struct {
	service domain.PatientService
}

func NewPatientHandler(s domain.PatientService) *PatientHandler {
	return &PatientHandler{service: s}
}

// GetByID handles GET /patient/search/:id
// The :id can be either a national_id or passport_id.
func (h *PatientHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hospitalID, _ := c.Locals("hospital_id").(string)

	patient, err := h.service.GetPatientByID(id, hospitalID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Patient not found",
			})
		case errors.Is(err, domain.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "Patient ID is required",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			})
		}
	}

	return c.JSON(patient)
}

// Search handles POST /patient/search
func (h *PatientHandler) Search(c *fiber.Ctx) error {
	var input domain.PatientSearchInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_INPUT",
			"message": "Invalid request body",
		})
	}

	hospitalID, _ := c.Locals("hospital_id").(string)

	patients, err := h.service.GetPatientByCondition(input, hospitalID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "At least one search condition is required",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			})
		}
	}

	return c.JSON(patients)
}
