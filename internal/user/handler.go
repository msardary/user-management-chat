package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"user-management/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetProfileHandler
// @Summary Get user profile
// @Description Get the profile of the authenticated user.
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse "Validation errors"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/user [get]
func GetProfileHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
        if !exists {
			response.Error(c, http.StatusUnauthorized, "User ID not found in context")
            return
        }

		user, err := service.GetUserByID(context.Background(), int(userID.(int64)))
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		response.Success(c, http.StatusOK, user, "")

	}

}

// UpdateProfileHandler
// @Summary Update a user's profile by ID
// @Description Updates a user's profile using their ID and request body.
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body UpdateUser true "User update data"
// @Success 204 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse "Invalid input"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/user [put]
func UpdateMyProfileHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		userID, exists := c.Get("userID")
        if !exists {
			response.Error(c, http.StatusUnauthorized, "User ID not found in context")
            return
        }

		var payload UpdateMyProfile

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		id := int(userID.(int64))

		if err := service.UpdateUserSelf(context.Background(), int(id), payload); err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		response.Success(c, http.StatusNoContent, nil, "User updated successfully")

	}

}

// @Summary Get all users
// @Description Get a list of all users (admin only).
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} response.PaginatedData "List of users"
// @Failure 400 {object} response.PaginatedData "Invalid request"
// @Failure 500 {object} response.PaginatedData "Internal server error"
// @Router /admin/api/v1/users [get]
func GetAllUsersHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		page := 1
		perPage := 10

		if p:= c.Query("page"); p != "" {
			fmt.Sscanf(p, "%d", &page)
		}
		if pp := c.Query("per_page"); pp != "" {
			fmt.Sscanf(pp, "%d", &perPage)
		}

		offset := (page - 1) * perPage

		total, err := service.db.CountUsers(context.Background())
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		users, err := service.GetUsers(context.Background(), perPage, offset)
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		response.Paginated(c, users, int(total), page, perPage)

	}

}

// @Summary Update a user profile by ID (admin only)
// @Description Update a user's profile by ID (admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param data body UpdateUser true "User update data"
// @Success 204 {object} response.APIResponse "User updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid input"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /admin/api/v1/users/{id} [put]
func UserUpdateHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		idParam := c.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var payload UpdateUser

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := service.UpdateUser(context.Background(), int(id), payload); err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		response.Success(c, http.StatusNoContent, nil, "User updated successfully")

	}

}

// @Summary Delete a user profile by ID (admin only)
// @Description Delete a user's profile by ID (admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} response.APIResponse "User deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid user ID"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /admin/api/v1/users/{id} [delete]
func UserDeleteHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		idParam := c.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid user ID")
			return
		}

		if err := service.DeleteUser(int(id)); err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		response.Success(c, http.StatusNoContent, nil, "User deleted successfully")

	}

}


