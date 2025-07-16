package image

import (
	"context"
	"image-app/internal/model"
	"image-app/internal/pkg/pagination"
	"io"
)

type ImageRepositoryInterface interface {
	Store(ctx context.Context, data model.ImageMeta) error
	Get(ctx context.Context, params model.ImageSearchParams, page pagination.Page) ([]model.ImageMeta, pagination.Metadata, error)
}

type ImageLogicInterface interface {
	UploadImage(ctx context.Context, file io.Reader, fileName string, size int64) (model.ImageMeta, error)
	GetImages(ctx context.Context, params model.ImageSearchParams, page pagination.Page) ([]model.ImageMeta, pagination.Metadata, error)
}
