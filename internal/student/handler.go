package student

import (
	"errors"
	"log"
	"net/http"
	"student-go-service/internal/dto"
	"student-go-service/internal/package/utils"
	"student-go-service/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler processes HTTP requests for student operations.
type Handler struct {
	service *Service
}

// NewHandler creates a student handler with its service dependency.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// StudentIDParams contains the student ID from a request path.
type StudentIDParams struct {
	Id uint `uri:"studentId" binding:"required"`
}

// StudentQuery contains filters and pagination from a list request.
type StudentQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	SortBy string `form:"sort_by"`
	Order  string `form:"order"`
}

// toStudentResponse converts a database student into an API response.
func toStudentResponse(student *Student, photoUrl string) dto.StudentResponse {
	return dto.StudentResponse{
		ID:        student.ID,
		Name:      student.Name,
		Email:     student.Email,
		Age:       student.Age,
		PhotoUrl:  photoUrl,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
	}
}

// toAllStudents converts multiple students into API responses.
func toAllStudents(students []Student) []dto.StudentResponse {
	response := make([]dto.StudentResponse, 0, len(students))
	for i := range students {
		response = append(response, toStudentResponse(&students[i], ""))
	}
	return response
}

// Create validates and creates a student from the request body.
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateStudentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}

	if req.Photo != nil {
		if err := utils.ValidateImage(req.Photo); err != nil {
			response.Error(c, http.StatusBadRequest, "Image Formate is not valid", err.Error())
			return
		}

		student, err := h.service.Create(req)
		if err != nil {

			if errors.Is(err, ErrEmailAlreadyExist) {
				response.Error(c, http.StatusConflict, ErrEmailAlreadyExist.Error(), "Using another work email!")
				return
			}

			response.Error(c, http.StatusInternalServerError, "Failed To Create Student", err.Error())
			return
		}

		photoUrl, err := h.service.GetPhotoUrl(c, student.Photo)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed To Generate Photo Url", err.Error())
			return
		}

		response.Success(c, http.StatusCreated, "Student Created Successfully", toStudentResponse(student, photoUrl))
	}
}

// GetAll returns students using filters and pagination.
func (h *Handler) GetAll(c *gin.Context) {
	var query StudentQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Request Query", err.Error())
		return
	}
	sortBy := utils.GetSortColumn(query.SortBy)
	order := utils.GetSortOrder(query.Order)

	pagination := utils.NewPagination(
		query.Page,
		query.Limit,
	)
	student, total, err := h.service.GetAll(c.Request.Context(), pagination, query.Search, sortBy, order)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch students", err.Error())
		return
	}
	meta := utils.NewPaginationMeta(pagination, total)

	response.Success(c, http.StatusOK, "Students Fetched successfully ", gin.H{
		"items": student,
		"meta":  meta,
	})
}

// GetById returns one student identified by the request path.
func (h *Handler) GetById(c *gin.Context) {
	var params StudentIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}
	log.Println("Handler", params.Id)
	student, err := h.service.GetById(params.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Student Not Found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to Fetch Student!", err.Error())
		return

	}
	photoUrl, err := h.service.GetPhotoUrl(c, student.Photo)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed To Generate Photo Url", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Student Fetched Successfully!", toStudentResponse(student, photoUrl))
}

// Update validates and saves changes to an existing student.
func (h *Handler) Update(c *gin.Context) {

	var params StudentIDParams

	if err := c.ShouldBindUri(&params); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid Request Params",
			err.Error(),
		)
		return
	}

	var req dto.UpdateStudentRequest

	if err := c.ShouldBind(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid Request Body",
			err.Error(),
		)
		return
	}

	// Validate new image BEFORE service update.
	if req.Photo != nil {

		if err := utils.ValidateImage(req.Photo); err != nil {
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid Image",
				err.Error(),
			)
			return
		}
	}

	student, err := h.service.Update(
		params.Id,
		&req,
	)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				"Student Not Found",
				err.Error(),
			)
			return
		}

		if errors.Is(err, ErrEmailAlreadyExist) {
			response.Error(
				c,
				http.StatusConflict,
				ErrEmailAlreadyExist.Error(),
				"Using Another Work Email",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to Update Student",
			err.Error(),
		)
		return
	}

	// Generate presigned URL for the new/current photo.
	photoURL, err := h.service.GetPhotoUrl(
		c.Request.Context(),
		student.Photo,
	)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed To Generate Photo URL",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Successfully Updated Student",
		toStudentResponse(student, photoURL),
	)
}

// Delete removes a student identified by the request path.
func (h *Handler) Delete(c *gin.Context) {
	var params StudentIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Request Body", err.Error())
		return
	}
	log.Println("Handler", params.Id)
	if err := h.service.Delete(params.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Student Not Found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed To Delete Student", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "DELETED", params.Id)

}
