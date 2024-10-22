package config

import "os"

type CloudinaryConfig struct {
	CloudName    string
	ApiKey       string
	ApiSecret    string
	UploadFolder string
}

func initCloudinaryConfig() CloudinaryConfig {
	return CloudinaryConfig{
		CloudName:    os.Getenv("CLOUDINARY_CLOUD_NAME"),
		ApiKey:       os.Getenv("CLOUDINARY_API_KEY"),
		ApiSecret:    os.Getenv("CLOUDINARY_API_SECRET"),
		UploadFolder: os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),
	}
}
