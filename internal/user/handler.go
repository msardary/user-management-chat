package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"user-management/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetProfile(service *Service) gin.HandlerFunc {

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

func UpdateProfile(service *Service) gin.HandlerFunc {

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

func GetAllUsers(service *Service) gin.HandlerFunc {

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


