package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gen-ai-bootcamp-2025/backend_go/internal/models"
	"github.com/gen-ai-bootcamp-2025/backend_go/internal/service"
	"github.com/stretchr/testify/assert"
)

// MockDB implements the necessary database methods for testing
type MockDB struct{}

func (m *MockDB) GetWord(id int64) (*models.Word, error) {
	return &models.Word{ID: id, Japanese: "テスト", Romaji: "tesuto", English: "test"}, nil
}

func (m *MockDB) GetWords(page, perPage int) ([]*models.Word, *models.Pagination, error) {
	words := []*models.Word{
		{ID: 1, Japanese: "テスト", Romaji: "tesuto", English: "test"},
	}
	pagination := &models.Pagination{
		CurrentPage:  page,
		TotalPages:   1,
		TotalItems:   1,
		ItemsPerPage: perPage,
	}
	return words, pagination, nil
}

func (m *MockDB) GetGroup(id int64) (*models.Group, error) {
	return &models.Group{ID: id, Name: "Test Group", WordCount: 10}, nil
}

func (m *MockDB) GetGroups(page, perPage int) ([]*models.Group, *models.Pagination, error) {
	groups := []*models.Group{
		{ID: 1, Name: "Test Group", WordCount: 10},
	}
	pagination := &models.Pagination{
		CurrentPage:  page,
		TotalPages:   1,
		TotalItems:   1,
		ItemsPerPage: perPage,
	}
	return groups, pagination, nil
}

func (m *MockDB) CreateStudySession(groupID, activityID int64) (*models.StudySession, error) {
	return &models.StudySession{ID: 1, GroupID: groupID, StudyActivityID: activityID}, nil
}

func (m *MockDB) GetStudySession(id int64) (*models.StudySession, error) {
	return &models.StudySession{ID: id, GroupID: 1, GroupName: "Test Group"}, nil
}

func (m *MockDB) GetLastStudySession() (*models.StudySession, error) {
	return &models.StudySession{
		ID:       1,
		GroupID:  1,
		GroupName: "Test Group",
	}, nil
}

func (m *MockDB) GetStudyProgress() (*models.StudyProgress, error) {
	return &models.StudyProgress{
		TotalWordsStudied:    10,
		TotalAvailableWords: 100,
	}, nil
}

func (m *MockDB) GetStudySessionsByActivity(activityID int64, page, perPage int) ([]*models.StudySession, *models.Pagination, error) {
	sessions := []*models.StudySession{
		{ID: 1, GroupID: 1, GroupName: "Test Group"},
	}
	pagination := &models.Pagination{
		CurrentPage:  page,
		TotalPages:   1,
		TotalItems:   1,
		ItemsPerPage: perPage,
	}
	return sessions, pagination, nil
}

func (m *MockDB) CreateWordReview(wordID, sessionID int64, correct bool) (*models.WordReviewItem, error) {
	return &models.WordReviewItem{ID: 1, WordID: wordID, StudySessionID: sessionID, Correct: correct}, nil
}

func (m *MockDB) GetQuickStats() (*models.QuickStats, error) {
	return &models.QuickStats{
		SuccessRate:        0.75,
		TotalStudySessions: 5,
		TotalActiveGroups:  2,
		StudyStreakDays:    3,
	}, nil
}

func (m *MockDB) ResetHistory() error {
	return nil
}

func (m *MockDB) FullReset() error {
	return nil
}

func (m *MockDB) GetWordsByGroup(groupID int64, page, perPage int) ([]*models.Word, *models.Pagination, error) {
	words := []*models.Word{
		{ID: 1, Japanese: "テスト", Romaji: "tesuto", English: "test"},
	}
	pagination := &models.Pagination{
		CurrentPage:  page,
		TotalPages:   1,
		TotalItems:   1,
		ItemsPerPage: perPage,
	}
	return words, pagination, nil
}

func (m *MockDB) GetWordsByStudySession(sessionID int64, page, perPage int) ([]*models.Word, *models.Pagination, error) {
	words := []*models.Word{
		{ID: 1, Japanese: "テスト", Romaji: "tesuto", English: "test"},
	}
	pagination := &models.Pagination{
		CurrentPage:  page,
		TotalPages:   1,
		TotalItems:   1,
		ItemsPerPage: perPage,
	}
	return words, pagination, nil
}

func setupTestRouter(t *testing.T) (*gin.Engine, *service.Service) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Initialize mock database
	mockDB := &MockDB{}

	// Initialize service with mock database
	svc := service.NewService(mockDB)

	// Setup router with handlers
	router := gin.New()
	handler := NewHandler(svc)

	// Register routes
	router.GET("/study-sessions/last", handler.GetLastStudySession)
	router.GET("/study/progress", handler.GetStudyProgress)
	router.GET("/stats/quick", handler.GetQuickStats)
	router.GET("/words", handler.GetWords)
	router.GET("/groups", handler.GetGroups)
	router.POST("/study/review", handler.ReviewWord)

	return router, svc
}

func TestGetLastStudySession(t *testing.T) {
	router, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/study-sessions/last", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Test Group", response["group_name"])
}

func TestGetStudyProgress(t *testing.T) {
	router, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/study/progress", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.StudyProgress
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 10, response.TotalWordsStudied)
	assert.Equal(t, 100, response.TotalAvailableWords)
}

func TestGetQuickStats(t *testing.T) {
	router, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats/quick", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.QuickStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0.75, response.SuccessRate)
	assert.Equal(t, 5, response.TotalStudySessions)
}

func TestGetWords(t *testing.T) {
	router, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/words", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	items, ok := response["items"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(items))
}

func TestGetGroups(t *testing.T) {
	router, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	items, ok := response["items"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(items))
}

func TestReviewWordWithInvalidInput(t *testing.T) {
	router, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/study/review", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
