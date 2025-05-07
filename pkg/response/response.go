package response

import (
	"errors"
	"net/http"
	"strconv"
	"user-management/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginatedData struct {
	Items 		interface{} `json:"items"`
	Total 		int 	   	`json:"total"`
	Page 		int 	   	`json:"page"`
	PerPage 	int 	   	`json:"per_page"`
	TotalPages 	int 		`json:"total_pages"`
}

func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	
	logSuccess(c, strconv.Itoa(http.StatusOK), "Request completed successfully")
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})

}

func Created(c *gin.Context, data interface{}) {
	
	logSuccess(c, strconv.Itoa(http.StatusCreated), "Request completed successfully")
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})

}

func Error(c *gin.Context, statusCode int, errors interface{}) {
	
	logError(c, statusCode, errors)
	c.JSON(statusCode, APIResponse{
		Success: false,
		Errors:  errors,
	})
	
}

func ValidationError(c *gin.Context, err error) {
	
	errors := validateErrors(err)
	logError(c, http.StatusBadRequest, errors)
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Errors:  errors,
	})
}

func Paginated(c *gin.Context, items interface{}, total, page, perPage int) {
	
	totalPages := (total + perPage - 1) / perPage

	logSuccess(c, strconv.Itoa(http.StatusOK), "Request completed successfully")
    c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Data: PaginatedData{
            Items:      items,
            Total:      total,
            Page:       page,
            PerPage:    perPage,
            TotalPages: totalPages,
        },
    })

}

func ParseDBError(c *gin.Context, err error) (int, string) {
	
	var pgErr *pgconn.PgError

	logError(c, http.StatusInternalServerError, err)
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique violation
			return http.StatusConflict, "Duplicate value entered."
		case "23503": // foreign key violation
			return http.StatusBadRequest, "No such record found."
		case "22P02": // invalid input syntax
			return http.StatusBadRequest, "Invalid input syntax"
		default:
			return http.StatusInternalServerError, "Internal server error"
		}

	}

	return http.StatusInternalServerError, "Internal server error"

}

func validateErrors(errs error) map[string]string {
	
	ValidationErrors := make(map[string]string)

	errors := errs.(validator.ValidationErrors)

	for _, e := range errors {
		ValidationErrors[e.Field()] = e.Translate(validation.Trans)
	}
	return ValidationErrors

}

func logSuccess(c *gin.Context, status, message string) {
	
	logger.WithFields(logrus.Fields{
		"status": status,
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
		"ip":     c.ClientIP(),
	}).Info(message)

}

func logError(c *gin.Context, statusCode int, err interface{}) {
	
	if statusCode >= 500 {
		logger.WithFields(logrus.Fields{
			"status": statusCode,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
		}).Error(err)
	} else if statusCode >= 400 {
		logger.WithFields(logrus.Fields{
			"status": statusCode,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
		}).Warn(err)
	}
	
}