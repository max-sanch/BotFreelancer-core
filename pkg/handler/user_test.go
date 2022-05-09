package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/service"
	mock_service "github.com/max-sanch/BotFreelancer-core/pkg/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_getTasksUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().GetTasks().Return([]core.UserTaskResponse{
					{
						TgId:  1111,
						Title: "Test",
						Body:  "TestBody",
						Url:   "TestUrl",
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"tasks":[{"tg_id":1111,"title":"Test","body":"TestBody","url":"TestUrl"}]}`,
		},
		{
			name: "Service Failure",
			mockBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().GetTasks().Return([]core.UserTaskResponse{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user)
			services := &service.Service{User: user}
			handler := NewHandler(services)

			// Test Server
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/getTasksUser", handler.getTasksUser)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/getTasksUser", bytes.NewBuffer([]byte{}))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, tgIdInput core.TgIdInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputTgId           core.TgIdInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"tg_id":1111}`,
			inputTgId: core.TgIdInput{
				TgId: 1111,
			},
			mockBehavior: func(s *mock_service.MockUser, tgIdInput core.TgIdInput) {
				s.EXPECT().GetByTgId(tgIdInput.TgId).Return(core.UserResponse{
					Id:       1,
					TgId:     1111,
					Username: "user-1",
					Setting: core.SettingResponse{
						IsSafeDeal: false,
						IsBudget:   false,
						IsTerm:     false,
						Categories: []int{1, 2},
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"tg_id":1111,"username":"user-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockUser, tgIdInput core.TgIdInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"tg_id":1111}`,
			inputTgId: core.TgIdInput{
				TgId: 1111,
			},
			mockBehavior: func(s *mock_service.MockUser, tgIdInput core.TgIdInput) {
				s.EXPECT().GetByTgId(tgIdInput.TgId).Return(core.UserResponse{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputTgId)
			services := &service.Service{User: user}
			handler := NewHandler(services)

			// Test Server
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/getUser", handler.getUser)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/getUser", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_createUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, userInput core.UserInput)
	isFalse := false

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           core.UserInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"tg_id":1111,"username":"user-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputUser: core.UserInput{
				TgId:     1111,
				Username: "user-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockUser, userInput core.UserInput) {
				s.EXPECT().Create(userInput).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"username":"user-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			mockBehavior:        func(s *mock_service.MockUser, userInput core.UserInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"tg_id":1111,"username":"user-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputUser: core.UserInput{
				TgId:     1111,
				Username: "user-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockUser, userInput core.UserInput) {
				s.EXPECT().Create(userInput).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputUser)
			services := &service.Service{User: user}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/createUser", handler.createUser)

			// Test Request
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/createUser", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_updateUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, userInput core.UserInput)
	isFalse := false

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           core.UserInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"tg_id":1111,"username":"user-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputUser: core.UserInput{
				TgId:     1111,
				Username: "user-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockUser, userInput core.UserInput) {
				s.EXPECT().Update(userInput).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"username":"user-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			mockBehavior:        func(s *mock_service.MockUser, userInput core.UserInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"tg_id":1111,"username":"user-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputUser: core.UserInput{
				TgId:     1111,
				Username: "user-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockUser, userInput core.UserInput) {
				s.EXPECT().Update(userInput).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputUser)
			services := &service.Service{User: user}
			handler := NewHandler(services)

			// Test Server
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/updateUser", handler.updateUser)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/updateUser", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
