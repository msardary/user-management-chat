package auth

import (
	"context"

	"net/http"
	"strings"
	"time"
	"user-management/pkg/response"
	"user-management/pkg/utils"
	"user-management/pkg/validation"

	"github.com/gin-gonic/gin"
)

type RefreshTokenStruct struct {
	UserID    int32     `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RegisterHandler
// @Summary Register a new user
// @Description Register a new user by providing username, email, and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param data body RegisterPayload true "User registration data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse "Validation errors"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/register [post]
func RegisterHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		var payload RegisterPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := validation.Validate.Struct(payload); err != nil {
			response.ValidationError(c, err)
			return
		}

		payload.Username = strings.ToLower(payload.Username)
		payload.Password, _ = utils.HashPassword(payload.Password)
		err := service.Register(context.Background(), payload)
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		response.Success(c, http.StatusCreated, nil, "User registered successfully")
	}

}

// LoginHandler
// @Summary User login
// @Description Login using username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param data body LoginPayload true "User login data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse "Invalid input"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/login [post]
func LoginHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		var req LoginPayload

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		user, err := service.GetUserByUsername(context.Background(), req.Username)
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		if(!utils.CheckPasswordHash(req.Password, user.PasswordHash)) {
			response.Error(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		accessToken, _ := utils.GenerateAccessToken(user.ID, user.Username, user.IsAdmin)
		refreshToken, _ := utils.GenerateRefreshToken(user.ID, user.Username, user.IsAdmin)

		params := RefreshTokenStruct{
			UserID: int32(user.ID),
			TokenHash: utils.HashToken(refreshToken),
			ExpiresAt: time.Now().Add(time.Hour*24*7),
		}

		err = service.SaveRefreshToken(context.Background(), params)
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		data := map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}
		response.Success(c, http.StatusOK, data, "Login successful")
	}

}

// RefreshTokenHandler
// @Summary Refresh access token
// @Description Provide a valid refresh token to get new access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param data body RefreshTokenPayload true "Refresh token data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse "Invalid input"
// @Failure 401 {object} response.APIResponse "Unauthorized or expired token"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/refresh [post]
func RefreshTokenHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		var req RefreshTokenPayload

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		hashedToken := utils.HashToken(req.RefreshToken)
		refresh, err := service.FindRefreshToken(context.Background(), hashedToken)
		if err != nil || refresh.Revoked || refresh.ExpiresAt.Before(time.Now()) {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}
		
		userInfo, err := service.GetUserByID(context.Background(), int(refresh.UserID))
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}
		newAccessToken, _ := utils.GenerateAccessToken(userInfo.ID, userInfo.Username, userInfo.IsAdmin)
		newRefreshToken, _ := utils.GenerateRefreshToken(userInfo.ID, userInfo.Username, userInfo.IsAdmin)
		
		params := RefreshTokenStruct{
			UserID: refresh.UserID,
			TokenHash: utils.HashToken(newRefreshToken),
			ExpiresAt: time.Now().Add(time.Hour*24*7),
		}
		err = service.SaveRefreshToken(context.Background(), params)
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		data := map[string]string{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		}
		response.Success(c, http.StatusOK, data, "Refresh token generated successfully")
	}

}

// LogoutHandler
// @Summary Logout user
// @Description Logs out the currently authenticated user by removing their refresh token
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/user/logout [get]
func LogoutHandler(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
        if !exists {
			response.Error(c, http.StatusUnauthorized, "User ID not found in context")
            return
        }

		id := int(userID.(int64))
		err := service.DeleteRefreshTokenByUserID(context.Background(), int32(id))
		if err != nil {
			status, msg := response.ParseDBError(c, err)
			response.Error(c, status, msg)
			return
		}

		response.Success(c, http.StatusOK, nil, "Logged out successfully")
	}

}