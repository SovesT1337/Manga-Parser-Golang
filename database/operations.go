package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ContentCreateNew(url string) (*Content, error) {
	c := &Content{UrlHentaichan: url, Status: "New"}
	return c, DB.Create(c).Error
}

func ContentExistsByURL(url string) (bool, error) {
	var count int64
	result := DB.Model(&Content{}).Where("url_hentaichan = ?", url).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func ContentGetByURL(url string) (*Content, error) {
	var content Content
	result := DB.Where("url_hentaichan = ?", url).First(&content)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &content, result.Error
}

func ContentGetByID(id uint) (*Content, error) {
	var content Content
	result := DB.First(&content, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &content, result.Error
}

func ContentClaimNew() (*Content, error) {
	var c Content
	tx := DB.Begin()
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("status = ?", "New").Order("id asc").First(&c).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(&c).Update("status", "Processing").Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func ContentMarkParsed(id uint, telegraphURL string) error {
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"url_telegraph": telegraphURL,
		"status":        "Parsed",
		"last_error":    "",
	}).Error
}

func ContentUpdateMeta(id uint, name, series, author, translator, tagsJSON string) error {
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"name":       name,
		"series":     series,
		"author":     author,
		"translator": translator,
		"tags_json":  tagsJSON,
	}).Error
}

func ContentFindDue(limit int) ([]Content, error) {
	var rows []Content
	q := DB.Where("status = ? AND scheduled_at <= NOW()", "Confirmed").Order("scheduled_at asc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func ContentLastScheduledAt() (*time.Time, error) {
	var row Content
	res := DB.Where("status = ? AND scheduled_at IS NOT NULL", "Confirmed").Order("scheduled_at desc").Limit(1).First(&row)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return row.ScheduledAt, nil
}

func ContentMarkSent(id uint) error {
	now := time.Now()
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"status":  "Sent",
		"sent_at": &now,
	}).Error
}

func ContentMarkError(id uint, errMsg string) error {
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"status":     "Error",
		"last_error": errMsg,
	}).Error
}

func ContentFindParsedPendingReview(limit int) ([]Content, error) {
	var rows []Content
	q := DB.Where("status = ? AND review_sent_at IS NULL", "Parsed").Order("id asc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func ContentMarkReviewSent(id uint) error {
	now := time.Now()
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"review_sent_at": &now,
	}).Error
}

func ContentMarkConfirmed(id uint) error {
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"status": "Confirmed",
	}).Error
}

func ContentMarkCancelled(id uint) error {
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"status": "Cancelled",
	}).Error
}

func ContentMarkConfirmedAndSchedule(id uint, scheduleAt time.Time) error {
	return DB.Model(&Content{}).Where("id = ?", id).Updates(map[string]any{
		"status":       "Confirmed",
		"scheduled_at": scheduleAt,
		"last_error":   "",
	}).Error
}

// Administrators

func AdminExists(userID int64) (bool, error) {
	var count int64
	res := DB.Model(&Administrator{}).Where("telegram_user_id = ?", userID).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count > 0, nil
}

func AdminList() ([]Administrator, error) {
	var rows []Administrator
	if err := DB.Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func AdminAdd(userID int64, username string) error {
	a := &Administrator{TelegramUserID: userID, Username: username}
	return DB.Create(a).Error
}
