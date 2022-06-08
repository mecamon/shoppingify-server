package storage

type Storage interface {
	UploadImage(file interface{}, fileName string) (string, error)
	DeleteImage(publicID string) (string, error)
}
