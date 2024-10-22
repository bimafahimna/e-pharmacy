package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

const maxFileSize = 500 * 1024

type cloudinaryMiddleware struct {
	config *config.CloudinaryConfig
}

func NewCloudinaryMiddleware(config *config.CloudinaryConfig) *cloudinaryMiddleware {
	return &cloudinaryMiddleware{
		config: config,
	}
}

func (m *cloudinaryMiddleware) Upload(useCase string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, handlerFile, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: apperror.ErrFileNotFound.Error(),
			})
			return
		}

		if handlerFile.Size > maxFileSize {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: apperror.ErrFileSizeExceededLimit.Error(),
			})
			return
		}

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var allowedContentTypes map[string]bool
		switch useCase {
		case "product", "partner":
			allowedContentTypes = map[string]bool{
				"image/png": true,
			}
		case "payment":
			allowedContentTypes = map[string]bool{
				"image/png":       true,
				"image/jpg":       true,
				"image/jpeg":      true,
				"application/pdf": true,
			}
		}

		contentType := handlerFile.Header.Get("Content-Type")
		if !allowedContentTypes[contentType] {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: apperror.ErrUnsupportedFileType.Error(),
			})
			return
		}

		cld, err := cloudinary.NewFromParams(m.config.CloudName, m.config.ApiKey, m.config.ApiSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Message: apperror.ErrInternalServerError.Error(),
			})
			return
		}

		name := time.Now().Format("2006-01-02_15:04:05") + "_" + handlerFile.Filename
		result, err := cld.Upload.Upload(timeout, file, uploader.UploadParams{
			PublicID: name,
			Folder:   m.config.UploadFolder,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Message: apperror.ErrInternalServerError.Error(),
			})
			return
		}

		ctx.Set("image_url", result.SecureURL)
		ctx.Next()
	}
}
