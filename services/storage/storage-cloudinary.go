package storage

import (
	"context"
	"errors"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/mecamon/shoppingify-server/config"
)

var storage *Cloudinary

type Cloudinary struct {
	service *cloudinary.Cloudinary
}

func InitStorage() (Storage, error) {
	app := config.Get()
	cld, err := cloudinary.NewFromParams(app.StorageCloudName, app.StorageAPIKey, app.StorageAPISecret)
	if err != nil {
		return nil, err
	}
	storage = &Cloudinary{service: cld}
	return storage, nil
}

func GetStorage() (Storage, error) {
	if storage == nil {
		return nil, errors.New("instance not created yet. Make sure to call InitStorage() before calling this method")
	}
	return storage, nil
}

func (s *Cloudinary) UploadImage(file interface{}, filename string) (string, error) {
	ctx := context.Background()
	uploadResult, err := storage.service.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			PublicID:       filename,
			UniqueFilename: true,
			Folder:         "shoppingify-app",
			Overwrite:      true,
		})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func (s *Cloudinary) DeleteImage(publicID string) (string, error) {
	ctx := context.Background()
	resp, err := storage.service.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image"})

	if err != nil {
		return "", err
	}

	return resp.Result, nil
}
