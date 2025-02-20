package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gen-ai-bootcamp-2025/backend_go/internal/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// GetLastStudySession returns the most recent study session
func (h *Handler) GetLastStudySession(c *gin.Context) {
	session, err := h.svc.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no study sessions found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

// GetStudyProgress returns study progress statistics
func (h *Handler) GetStudyProgress(c *gin.Context) {
	progress, err := h.svc.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

// GetQuickStats returns quick overview statistics
func (h *Handler) GetQuickStats(c *gin.Context) {
	stats, err := h.svc.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetStudyActivity returns a specific study activity
func (h *Handler) GetStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// TODO: Implement getting study activity
	activity := gin.H{
		"id":           id,
		"name":         "Vocabulary Quiz",
		"thumbnail_url": "https://example.com/thumbnail.jpg",
		"description":  "Practice your vocabulary with flashcards",
	}
	c.JSON(http.StatusOK, activity)
}

// GetStudyActivitySessions returns study sessions for a specific activity
func (h *Handler) GetStudyActivitySessions(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	response, err := h.svc.GetStudySessionsByActivity(activityID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// CreateStudyActivity creates a new study activity
func (h *Handler) CreateStudyActivity(c *gin.Context) {
	var req struct {
		GroupID         int64 `json:"group_id" binding:"required"`
		StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.svc.CreateStudySession(req.GroupID, req.StudyActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       session.ID,
		"group_id": session.GroupID,
	})
}

// GetWords returns a paginated list of words
func (h *Handler) GetWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	response, err := h.svc.GetWords(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetWord returns a specific word
func (h *Handler) GetWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	word, err := h.svc.GetWord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}

// GetGroups returns a paginated list of groups
func (h *Handler) GetGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	response, err := h.svc.GetGroups(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetGroup returns a specific group
func (h *Handler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	group, err := h.svc.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

// GetGroupWords returns words for a specific group
func (h *Handler) GetGroupWords(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	response, err := h.svc.GetWordsByGroup(groupID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetGroupStudySessions returns study sessions for a specific group
func (h *Handler) GetGroupStudySessions(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{
		"items": []gin.H{},
		"pagination": gin.H{
			"current_page":   1,
			"total_pages":    1,
			"total_items":    0,
			"items_per_page": 100,
		},
	})
}

// GetStudySessions returns all study sessions
func (h *Handler) GetStudySessions(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{
		"items": []gin.H{},
		"pagination": gin.H{
			"current_page":   1,
			"total_pages":    1,
			"total_items":    0,
			"items_per_page": 100,
		},
	})
}

// GetStudySession returns a specific study session
func (h *Handler) GetStudySession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	session, err := h.svc.GetStudySession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

// GetStudySessionWords returns words for a specific study session
func (h *Handler) GetStudySessionWords(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	response, err := h.svc.GetWordsByStudySession(sessionID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// ReviewWord records a word review in a study session
func (h *Handler) ReviewWord(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word id"})
		return
	}

	var req struct {
		Correct bool `json:"correct" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.svc.ReviewWord(wordID, sessionID, req.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// ResetHistory resets study history
func (h *Handler) ResetHistory(c *gin.Context) {
	if err := h.svc.ResetHistory(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Study history has been reset",
	})
}

// FullReset performs a complete system reset
func (h *Handler) FullReset(c *gin.Context) {
	if err := h.svc.FullReset(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "System has been fully reset",
	})
}
