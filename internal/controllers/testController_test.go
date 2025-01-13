package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"userManagement/internal/mocks"
	"userManagement/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	UserController := UserController{userService: mockUserService}
	router.POST("user-signup", UserController.Signup)

	type TypeCase struct {
		name               string
		requestBody        models.User
		expectedStatusCode int
		expectedResponse   string
		returnError        error
		expectSignupCall   bool
	}

	tests := []TypeCase{
		{
			name: "successful signup",
			requestBody: models.User{
				FirstName: "sreejith",
				LastName:  "k",
				Email:     "sreejith@gmail.com",
				Password:  "Password@1",
				Phone:     "1234567890",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   fmt.Sprintf(`{"message":"%v"}`, models.SignupSuccess),
			returnError:        nil,
			expectSignupCall:   true,
		},
		{
			name: "user already exists",
			requestBody: models.User{
				FirstName: "sreejith",
				LastName:  "k",
				Email:     "sreejith@gmail.com",
				Password:  "Password@1",
				Phone:     "1234567890",
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, models.UserExist),
			returnError:        errors.New(models.UserExist),
			expectSignupCall:   true,
		},
		{
			name: "invalid input - empty first name",
			requestBody: models.User{
				FirstName: "",
				LastName:  "k",
				Email:     "sreejith@gmail.com",
				Password:  "Password@1",
				Phone:     "1234567890",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%v"}`, "FirstName is required"),
			returnError:        nil,
			expectSignupCall:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectSignupCall {
				if test.returnError != nil {
					mockUserService.EXPECT().Signup(&test.requestBody).Return(test.returnError).Times(1)
				} else {
					mockUserService.EXPECT().Signup(&test.requestBody).Return(nil).Times(1)
				}
			}

			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/user-signup", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.expectedStatusCode, recorder.Code)

			assert.JSONEq(t, test.expectedResponse, recorder.Body.String())
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	userController := &UserController{userService: mockUserService}

	router.POST("/user-login", userController.Login)

	tests := []struct {
		name               string
		requestBody        models.UserLogin
		mockSetup          func()
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name: "successful login",
			requestBody: models.UserLogin{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				user := &models.User{
					Email:    "test@example.com",
					Password: "password123",
				}
				mockUserService.EXPECT().Login(&models.UserLogin{
					Email:    "test@example.com",
					Password: "password123",
				}).Return(user, nil)
				mockUserService.EXPECT().VerifyPassword(&models.UserLogin{
					Email:    "test@example.com",
					Password: "password123",
				}, user).Return(true)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Login Succesful",
			},
		},
		{
			name: "user not found",
			requestBody: models.UserLogin{
				Email:    "unknown@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockUserService.EXPECT().Login(&models.UserLogin{
					Email:    "unknown@example.com",
					Password: "password123",
				}).Return(nil, errors.New(models.UserNotFound))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"error": models.UserNotFound,
			},
		},
		{
			name: "incorrect password",
			requestBody: models.UserLogin{
				Email:    "test@example.com",
				Password: "wrongPassword",
			},
			mockSetup: func() {
				user := &models.User{
					Email:    "test@example.com",
					Password: "password123",
				}
				mockUserService.EXPECT().Login(&models.UserLogin{
					Email:    "test@example.com",
					Password: "wrongPassword",
				}).Return(user, nil)
				mockUserService.EXPECT().VerifyPassword(&models.UserLogin{
					Email:    "test@example.com",
					Password: "wrongPassword",
				}, user).Return(false)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": models.IncorrectPassword,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mockSetup()

			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/user-login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)

			var actualResponse map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&actualResponse)
			assert.NoError(t, err)

			for key, value := range test.expectedResponse {
				assert.Equal(t, value, actualResponse[key])
			}
		})
	}
}
