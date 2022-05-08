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

func TestHandler_getTasksChannel(t *testing.T) {
	type mockBehavior func(s *mock_service.MockChannel)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockChannel) {
				s.EXPECT().GetTasks().Return([]core.ChannelTaskResponse{
					{
						ApiId:   1111,
						ApiHash: "hash1111",
						Title:   "Test",
						Body:    "TestBody",
						Url:     "TestUrl",
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"tasks":[{"api_id":1111,"api_hash":"hash1111","title":"Test","body":"TestBody","url":"TestUrl"}]}`,
		},
		{
			name: "Service Failure",
			mockBehavior: func(s *mock_service.MockChannel) {
				s.EXPECT().GetTasks().Return([]core.ChannelTaskResponse{}, errors.New("service failure"))
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

			channel := mock_service.NewMockChannel(c)
			testCase.mockBehavior(channel)
			services := &service.Service{Channel: channel}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/getTasksChannel", handler.getTasksChannel)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/getTasksChannel", bytes.NewBuffer([]byte{}))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getChannel(t *testing.T) {
	type mockBehavior func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputApiId          core.ApiIdInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"api_id":1111}`,
			inputApiId: core.ApiIdInput{
				ApiId: 1111,
			},
			mockBehavior: func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput) {
				s.EXPECT().GetByApiId(apiIdInput.ApiId).Return(core.ChannelResponse{
					Id:      1,
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingResponse{
						IsSafeDeal: false,
						IsBudget:   false,
						IsTerm:     false,
						Categories: []int{1, 2},
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"api_id":1111,"api_hash":"hash1111","name":"channel-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"api_id":1111}`,
			inputApiId: core.ApiIdInput{
				ApiId: 1111,
			},
			mockBehavior: func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput) {
				s.EXPECT().GetByApiId(apiIdInput.ApiId).Return(core.ChannelResponse{}, errors.New("service failure"))
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

			channel := mock_service.NewMockChannel(c)
			testCase.mockBehavior(channel, testCase.inputApiId)
			services := &service.Service{Channel: channel}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/getChannel", handler.getChannel)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/getChannel", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_createChannel(t *testing.T) {
	type mockBehavior func(s *mock_service.MockChannel, channelInput core.ChannelInput)
	isFalse := false

	testTable := []struct {
		name                string
		inputBody           string
		inputChannel        core.ChannelInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"api_id":1111,"api_hash":"hash1111","name":"channel-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputChannel: core.ChannelInput{
				ApiId:   1111,
				ApiHash: "hash1111",
				Name:    "channel-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockChannel, channelInput core.ChannelInput) {
				s.EXPECT().Create(channelInput).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"api_hash":"hash1111","name":"channel-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			mockBehavior:        func(s *mock_service.MockChannel, channelInput core.ChannelInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"api_id":1111,"api_hash":"hash1111","name":"channel-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputChannel: core.ChannelInput{
				ApiId:   1111,
				ApiHash: "hash1111",
				Name:    "channel-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockChannel, channelInput core.ChannelInput) {
				s.EXPECT().Create(channelInput).Return(0, errors.New("service failure"))
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

			channel := mock_service.NewMockChannel(c)
			testCase.mockBehavior(channel, testCase.inputChannel)
			services := &service.Service{Channel: channel}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/createChannel", handler.createChannel)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/createChannel", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_updateChannel(t *testing.T) {
	type mockBehavior func(s *mock_service.MockChannel, channelInput core.ChannelInput)
	isFalse := false

	testTable := []struct {
		name                string
		inputBody           string
		inputChannel        core.ChannelInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"api_id":1111,"api_hash":"hash1111","name":"channel-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputChannel: core.ChannelInput{
				ApiId:   1111,
				ApiHash: "hash1111",
				Name:    "channel-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockChannel, channelInput core.ChannelInput) {
				s.EXPECT().Update(channelInput).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"api_hash":"hash1111","name":"channel-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			mockBehavior:        func(s *mock_service.MockChannel, channelInput core.ChannelInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"api_id":1111,"api_hash":"hash1111","name":"channel-1","setting":{"is_safe_deal":false,"is_budget":false,"is_term":false,"categories":[1,2]}}`,
			inputChannel: core.ChannelInput{
				ApiId:   1111,
				ApiHash: "hash1111",
				Name:    "channel-1",
				Setting: core.SettingInput{
					IsSafeDeal: &isFalse,
					IsBudget:   &isFalse,
					IsTerm:     &isFalse,
					Categories: []int{1, 2},
				},
			},
			mockBehavior: func(s *mock_service.MockChannel, channelInput core.ChannelInput) {
				s.EXPECT().Update(channelInput).Return(0, errors.New("service failure"))
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

			channel := mock_service.NewMockChannel(c)
			testCase.mockBehavior(channel, testCase.inputChannel)
			services := &service.Service{Channel: channel}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/updateChannel", handler.updateChannel)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/updateChannel", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_deleteChannel(t *testing.T) {
	type mockBehavior func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputApiId          core.ApiIdInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"api_id":1111}`,
			inputApiId: core.ApiIdInput{
				ApiId: 1111,
			},
			mockBehavior: func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput) {
				s.EXPECT().Delete(apiIdInput.ApiId).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"api_id":1111}`,
			inputApiId: core.ApiIdInput{
				ApiId: 1111,
			},
			mockBehavior: func(s *mock_service.MockChannel, apiIdInput core.ApiIdInput) {
				s.EXPECT().Delete(apiIdInput.ApiId).Return(errors.New("service failure"))
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

			channel := mock_service.NewMockChannel(c)
			testCase.mockBehavior(channel, testCase.inputApiId)
			services := &service.Service{Channel: channel}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/deleteChannel", handler.deleteChannel)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/deleteChannel", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
