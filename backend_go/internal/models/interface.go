package models

// DBInterface defines the interface that both the real DB and mock DB must implement
type DBInterface interface {
	GetWord(id int64) (*Word, error)
	GetWords(page, perPage int) ([]*Word, *Pagination, error)
	GetGroup(id int64) (*Group, error)
	GetGroups(page, perPage int) ([]*Group, *Pagination, error)
	CreateStudySession(groupID, activityID int64) (*StudySession, error)
	GetStudySession(id int64) (*StudySession, error)
	GetLastStudySession() (*StudySession, error)
	GetStudyProgress() (*StudyProgress, error)
	GetStudySessionsByActivity(activityID int64, page, perPage int) ([]*StudySession, *Pagination, error)
	CreateWordReview(wordID, sessionID int64, correct bool) (*WordReviewItem, error)
	GetQuickStats() (*QuickStats, error)
	ResetHistory() error
	FullReset() error
	GetWordsByGroup(groupID int64, page, perPage int) ([]*Word, *Pagination, error)
	GetWordsByStudySession(sessionID int64, page, perPage int) ([]*Word, *Pagination, error)
}
