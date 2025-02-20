package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gen-ai-bootcamp-2025/backend_go/internal/handlers"
	"github.com/gen-ai-bootcamp-2025/backend_go/internal/models"
	"github.com/gen-ai-bootcamp-2025/backend_go/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func setupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	// Dashboard routes
	api.GET("/dashboard/last_study_session", h.GetLastStudySession)
	api.GET("/dashboard/study_progress", h.GetStudyProgress)
	api.GET("/dashboard/quick-stats", h.GetQuickStats)

	// Study activities routes
	api.GET("/study_activities/:id", h.GetStudyActivity)
	api.GET("/study_activities/:id/study_sessions", h.GetStudyActivitySessions)
	api.POST("/study_activities", h.CreateStudyActivity)

	// Words routes
	api.GET("/words", h.GetWords)
	api.GET("/words/:id", h.GetWord)

	// Groups routes
	api.GET("/groups", h.GetGroups)
	api.GET("/groups/:id", h.GetGroup)
	api.GET("/groups/:id/words", h.GetGroupWords)
	api.GET("/groups/:id/study_sessions", h.GetGroupStudySessions)

	// Study sessions routes
	api.GET("/study_sessions", h.GetStudySessions)
	api.GET("/study_sessions/:id", h.GetStudySession)
	api.GET("/study_sessions/:id/words", h.GetStudySessionWords)
	api.POST("/study_sessions/:id/words/:word_id/review", h.ReviewWord)

	// System routes
	api.POST("/reset_history", h.ResetHistory)
	api.POST("/full_reset", h.FullReset)

	return r
}

func main() {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	modelDB := models.NewDB(db)
	svc := service.NewService(modelDB)
	h := handlers.NewHandler(svc)

	r := setupRouter(h)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
