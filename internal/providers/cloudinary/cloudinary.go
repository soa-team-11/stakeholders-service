package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func Init(apiKey, apiSecret, cloudName string) error {
	var err error
	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	return err
}

func UploadImage(file multipart.File) (string, error) {
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: "profiles",
	})
	if err != nil {
		return "", err
	}
	fmt.Printf("Upload result: %+v\n", uploadResult)
	return uploadResult.SecureURL, nil
}
