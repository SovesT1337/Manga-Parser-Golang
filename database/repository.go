package database

// Интерфейс для работы с контентом
type ContentRepositoryInterface interface {
	Create(content *HContent) error
	ExistsByHentaichanURL(url string) (bool, error)
	GetByHentaichanURL(url string) (*HContent, error)
}

// Создание экземпляра репозитория
func NewContentRepository() ContentRepositoryInterface {
	return &ContentRepository{}
}