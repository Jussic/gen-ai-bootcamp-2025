package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type DB struct {
	*sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{db}
}

// Word operations
func (db *DB) GetWord(id int64) (*Word, error) {
	word := &Word{}
	err := db.QueryRow(`
		SELECT id, japanese, romaji, english, parts 
		FROM words WHERE id = ?`, id).Scan(
		&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts)
	if err != nil {
		return nil, err
	}

	if word.Parts.Valid {
		err = json.Unmarshal([]byte(word.Parts.String), &word.PartsMap)
		if err != nil {
			return nil, err
		}
	}

	return word, nil
}

func (db *DB) GetWords(page, perPage int) ([]*Word, *Pagination, error) {
	offset := (page - 1) * perPage
	words := []*Word{}

	rows, err := db.Query(`
		SELECT id, japanese, romaji, english, parts 
		FROM words LIMIT ? OFFSET ?`, perPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		word := &Word{}
		err := rows.Scan(&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts)
		if err != nil {
			return nil, nil, err
		}
		if word.Parts.Valid {
			err = json.Unmarshal([]byte(word.Parts.String), &word.PartsMap)
			if err != nil {
				return nil, nil, err
			}
		}
		words = append(words, word)
	}

	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, nil, err
	}

	pagination := &Pagination{
		CurrentPage:  page,
		TotalPages:   (total + perPage - 1) / perPage,
		TotalItems:   total,
		ItemsPerPage: perPage,
	}

	return words, pagination, nil
}

// Group operations
func (db *DB) GetGroup(id int64) (*Group, error) {
	group := &Group{}
	err := db.QueryRow(`
		SELECT g.id, g.name, COUNT(wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		WHERE g.id = ?
		GROUP BY g.id`, id).Scan(&group.ID, &group.Name, &group.WordCount)
	return group, err
}

func (db *DB) GetGroups(page, perPage int) ([]*Group, *Pagination, error) {
	offset := (page - 1) * perPage
	groups := []*Group{}

	rows, err := db.Query(`
		SELECT g.id, g.name, COUNT(wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		GROUP BY g.id
		LIMIT ? OFFSET ?`, perPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		group := &Group{}
		err := rows.Scan(&group.ID, &group.Name, &group.WordCount)
		if err != nil {
			return nil, nil, err
		}
		groups = append(groups, group)
	}

	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
	if err != nil {
		return nil, nil, err
	}

	pagination := &Pagination{
		CurrentPage:  page,
		TotalPages:   (total + perPage - 1) / perPage,
		TotalItems:   total,
		ItemsPerPage: perPage,
	}

	return groups, pagination, nil
}

// Study Session operations
func (db *DB) CreateStudySession(groupID, activityID int64) (*StudySession, error) {
	result, err := db.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)`, groupID, activityID, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.GetStudySession(id)
}

func (db *DB) GetStudySession(id int64) (*StudySession, error) {
	session := &StudySession{}
	err := db.QueryRow(`
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id,
			g.name as group_name,
			(SELECT COUNT(*) FROM word_review_items WHERE study_session_id = s.id) as review_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		WHERE s.id = ?`, id).Scan(
		&session.ID, &session.GroupID, &session.CreatedAt,
		&session.StudyActivityID, &session.GroupName, &session.ReviewItemCount)
	return session, err
}

func (db *DB) GetLastStudySession() (*StudySession, error) {
	session := &StudySession{}
	err := db.QueryRow(`
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id,
			g.name as group_name,
			(SELECT COUNT(*) FROM word_review_items WHERE study_session_id = s.id) as review_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		ORDER BY s.created_at DESC
		LIMIT 1`).Scan(
		&session.ID, &session.GroupID, &session.CreatedAt,
		&session.StudyActivityID, &session.GroupName, &session.ReviewItemCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

func (db *DB) GetStudyProgress() (*StudyProgress, error) {
	progress := &StudyProgress{}
	
	// Get total words studied (unique words that have been reviewed)
	err := db.QueryRow(`
		SELECT COUNT(DISTINCT word_id)
		FROM word_review_items`).Scan(&progress.TotalWordsStudied)
	if err != nil {
		return nil, err
	}

	// Get total available words
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM words`).Scan(&progress.TotalAvailableWords)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (db *DB) GetStudySessionsByActivity(activityID int64, page, perPage int) ([]*StudySession, *Pagination, error) {
	offset := (page - 1) * perPage
	sessions := []*StudySession{}

	rows, err := db.Query(`
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id,
			g.name as group_name,
			(SELECT COUNT(*) FROM word_review_items WHERE study_session_id = s.id) as review_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		WHERE s.study_activity_id = ?
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?`, activityID, perPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		session := &StudySession{}
		err := rows.Scan(
			&session.ID, &session.GroupID, &session.CreatedAt,
			&session.StudyActivityID, &session.GroupName, &session.ReviewItemCount)
		if err != nil {
			return nil, nil, err
		}
		sessions = append(sessions, session)
	}

	var total int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM study_sessions
		WHERE study_activity_id = ?`, activityID).Scan(&total)
	if err != nil {
		return nil, nil, err
	}

	pagination := &Pagination{
		CurrentPage:  page,
		TotalPages:   (total + perPage - 1) / perPage,
		TotalItems:   total,
		ItemsPerPage: perPage,
	}

	return sessions, pagination, nil
}

// Word Review operations
func (db *DB) CreateWordReview(wordID, sessionID int64, correct bool) (*WordReviewItem, error) {
	result, err := db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)`, wordID, sessionID, correct, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	review := &WordReviewItem{
		ID:             id,
		WordID:         wordID,
		StudySessionID: sessionID,
		Correct:        correct,
		CreatedAt:      time.Now(),
	}
	return review, nil
}

// Statistics operations
func (db *DB) GetQuickStats() (*QuickStats, error) {
	stats := &QuickStats{}

	// Get success rate
	err := db.QueryRow(`
		SELECT COALESCE(AVG(CASE WHEN correct THEN 100.0 ELSE 0.0 END), 0)
		FROM word_review_items`).Scan(&stats.SuccessRate)
	if err != nil {
		return nil, err
	}

	// Get total study sessions
	err = db.QueryRow(`
		SELECT COUNT(*) FROM study_sessions`).Scan(&stats.TotalStudySessions)
	if err != nil {
		return nil, err
	}

	// Get total active groups
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT group_id) 
		FROM study_sessions 
		WHERE created_at >= datetime('now', '-30 days')`).Scan(&stats.TotalActiveGroups)
	if err != nil {
		return nil, err
	}

	// Get study streak
	err = db.QueryRow(`
		WITH RECURSIVE dates(date) AS (
			SELECT date(MAX(created_at)) FROM study_sessions
			UNION ALL
			SELECT date(date, '-1 day')
			FROM dates
			WHERE date > date((
				SELECT MIN(created_at) FROM study_sessions
			))
		)
		SELECT COUNT(*)
		FROM dates d
		WHERE EXISTS (
			SELECT 1 FROM study_sessions
			WHERE date(created_at) = d.date
		)`).Scan(&stats.StudyStreakDays)

	return stats, err
}

// System operations
func (db *DB) ResetHistory() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM study_activities")
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) FullReset() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tables := []string{
		"word_review_items",
		"study_sessions",
		"study_activities",
		"words_groups",
		"words",
		"groups",
	}

	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM " + table)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *DB) GetWordsByGroup(groupID int64, page, perPage int) ([]*Word, *Pagination, error) {
	offset := (page - 1) * perPage
	words := []*Word{}

	rows, err := db.Query(`
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 1) as correct_count,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 0) as wrong_count
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		LIMIT ? OFFSET ?`, groupID, perPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		word := &Word{}
		var correctCount, wrongCount int
		err := rows.Scan(
			&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts,
			&correctCount, &wrongCount)
		if err != nil {
			return nil, nil, err
		}
		if word.Parts.Valid {
			err = json.Unmarshal([]byte(word.Parts.String), &word.PartsMap)
			if err != nil {
				return nil, nil, err
			}
		}
		words = append(words, word)
	}

	var total int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM words_groups
		WHERE group_id = ?`, groupID).Scan(&total)
	if err != nil {
		return nil, nil, err
	}

	pagination := &Pagination{
		CurrentPage:  page,
		TotalPages:   (total + perPage - 1) / perPage,
		TotalItems:   total,
		ItemsPerPage: perPage,
	}

	return words, pagination, nil
}

func (db *DB) GetWordsByStudySession(sessionID int64, page, perPage int) ([]*Word, *Pagination, error) {
	offset := (page - 1) * perPage
	words := []*Word{}

	rows, err := db.Query(`
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 1) as correct_count,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.word_id = w.id AND wri.correct = 0) as wrong_count
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		GROUP BY w.id
		LIMIT ? OFFSET ?`, sessionID, perPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		word := &Word{}
		var correctCount, wrongCount int
		err := rows.Scan(
			&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts,
			&correctCount, &wrongCount)
		if err != nil {
			return nil, nil, err
		}
		if word.Parts.Valid {
			err = json.Unmarshal([]byte(word.Parts.String), &word.PartsMap)
			if err != nil {
				return nil, nil, err
			}
		}
		words = append(words, word)
	}

	var total int
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT word_id)
		FROM word_review_items
		WHERE study_session_id = ?`, sessionID).Scan(&total)
	if err != nil {
		return nil, nil, err
	}

	pagination := &Pagination{
		CurrentPage:  page,
		TotalPages:   (total + perPage - 1) / perPage,
		TotalItems:   total,
		ItemsPerPage: perPage,
	}

	return words, pagination, nil
}
