package service

import (
	"github.com/gen-ai-bootcamp-2025/backend_go/internal/models"
)

type Service struct {
	db models.DBInterface
}

func NewService(db models.DBInterface) *Service {
	return &Service{db: db}
}

func (s *Service) GetWord(id int64) (*models.Word, error) {
	return s.db.GetWord(id)
}

func (s *Service) GetWords(page int) (*models.PaginatedResponse, error) {
	perPage := 100
	words, pagination, err := s.db.GetWords(page, perPage)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Items:      words,
		Pagination: *pagination,
	}, nil
}

func (s *Service) GetGroup(id int64) (*models.Group, error) {
	return s.db.GetGroup(id)
}

func (s *Service) GetGroups(page int) (*models.PaginatedResponse, error) {
	perPage := 100
	groups, pagination, err := s.db.GetGroups(page, perPage)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Items:      groups,
		Pagination: *pagination,
	}, nil
}

func (s *Service) CreateStudySession(groupID, activityID int64) (*models.StudySession, error) {
	return s.db.CreateStudySession(groupID, activityID)
}

func (s *Service) GetStudySession(id int64) (*models.StudySession, error) {
	return s.db.GetStudySession(id)
}

func (s *Service) ReviewWord(wordID, sessionID int64, correct bool) (*models.WordReviewItem, error) {
	return s.db.CreateWordReview(wordID, sessionID, correct)
}

func (s *Service) GetQuickStats() (*models.QuickStats, error) {
	return s.db.GetQuickStats()
}

func (s *Service) ResetHistory() error {
	return s.db.ResetHistory()
}

func (s *Service) FullReset() error {
	return s.db.FullReset()
}

func (s *Service) GetLastStudySession() (*models.StudySession, error) {
	return s.db.GetLastStudySession()
}

func (s *Service) GetStudyProgress() (*models.StudyProgress, error) {
	return s.db.GetStudyProgress()
}

func (s *Service) GetStudySessionsByActivity(activityID int64, page int) (*models.PaginatedResponse, error) {
	perPage := 100
	sessions, pagination, err := s.db.GetStudySessionsByActivity(activityID, page, perPage)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Items:      sessions,
		Pagination: *pagination,
	}, nil
}

func (s *Service) GetWordsByGroup(groupID int64, page int) (*models.PaginatedResponse, error) {
	perPage := 100
	words, pagination, err := s.db.GetWordsByGroup(groupID, page, perPage)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Items:      words,
		Pagination: *pagination,
	}, nil
}

func (s *Service) GetWordsByStudySession(sessionID int64, page int) (*models.PaginatedResponse, error) {
	perPage := 100
	words, pagination, err := s.db.GetWordsByStudySession(sessionID, page, perPage)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Items:      words,
		Pagination: *pagination,
	}, nil
}
