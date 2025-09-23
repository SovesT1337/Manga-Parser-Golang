package database

import (
	"errors"
	
	"gorm.io/gorm"
)

// Репозиторий для работы с контентом
type ContentRepository struct{}

// Создание новой записи
func (r *ContentRepository) Create(content *HContent) error {
	result := DB.Create(content)
	return result.Error
}

// Проверка существования по url_hentaichan
func (r *ContentRepository) ExistsByHentaichanURL(url string) (bool, error) {
	var count int64
	result := DB.Model(&HContent{}).Where("url_hentaichan = ?", url).Count(&count)
	
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// Получение записи по URL (дополнительная полезная функция)
func (r *ContentRepository) GetByHentaichanURL(url string) (*HContent, error) {
	var content HContent
	result := DB.Where("url_hentaichan = ?", url).First(&content)
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &content, result.Error
}