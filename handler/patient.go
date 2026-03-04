package handler

import (
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
func (h *PatientHandler) GetByID(c *fiber.Ctx) error {
	// TODO: implement
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Search handles POST /patient/search
func (h *PatientHandler) Search(c *fiber.Ctx) error {
	// TODO: implement
	return c.SendStatus(fiber.StatusNotImplemented)
}
