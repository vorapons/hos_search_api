package ginhandler

import (
	"errors"
	"net/http"

	"pt_search_hos/domain"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service domain.PatientService
}

func NewPatientHandler(s domain.PatientService) *PatientHandler {
	return &PatientHandler{service: s}
}

// GetByID handles GET /patient/search/:id
// The :id can be either a national_id or passport_id.
func (h *PatientHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	hospitalIDRaw, _ := c.Get("hospital_id")
	hospitalID, _ := hospitalIDRaw.(string)

	patient, err := h.service.GetPatientByID(id, hospitalID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "NOT_FOUND",
				"message": "Patient not found",
			})
		case errors.Is(err, domain.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_INPUT",
				"message": "Patient ID is required",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, patient)
}

// Search handles POST /patient/search
func (h *PatientHandler) Search(c *gin.Context) {
	var input domain.PatientSearchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_INPUT",
			"message": "Invalid request body",
		})
		return
	}

	hospitalIDRaw, _ := c.Get("hospital_id")
	hospitalID, _ := hospitalIDRaw.(string)

	result, err := h.service.GetPatientByCondition(input, hospitalID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_INPUT",
				"message": "At least one search condition is required",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
