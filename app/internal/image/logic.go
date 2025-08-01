package image

import (
	"context"
	"fmt"
	"image-app/internal/core"
	"image-app/internal/model"
	"image-app/internal/pkg/pagination"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type ImageLogic struct {
	deps    *core.Dependency
	imgRepo ImageRepositoryInterface
}

func NewImageLogic(deps *core.Dependency, imgRepo ImageRepositoryInterface) *ImageLogic {
	return &ImageLogic{
		deps:    deps,
		imgRepo: imgRepo,
	}
}

func (l *ImageLogic) UploadImage(ctx context.Context, file io.Reader, fileName string, size int64) (model.ImageMeta, error) {
	// generate uuid to avoid collision in name
	id := uuid.NewString()

	// getting file extension for name alteration
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	// setting up path to file
	dir := fmt.Sprintf("%s/%s-%s.%s", "static", name, id, ext)
	
	// create file to write
	dst, err := os.Create(dir)
	if err != nil {
		l.deps.Logger.ErrorContext(ctx, "failed to create file", slog.Any("error", err))
		return model.ImageMeta{}, err
	}
	defer dst.Close()

	// write the stream into the file
	_, err = io.Copy(dst, file)
	if err != nil {
		l.deps.Logger.ErrorContext(ctx, "failed to write file upload", slog.Any("error", err))
		return model.ImageMeta{}, err
	}

	// prepare to store file metadata
	data := model.ImageMeta{
		ID:       id,
		FileName: fileName,
		URL:      fmt.Sprintf("%s/%s", "http://localhost", dir),
		FileSize: size,
	}
	
	err = l.imgRepo.Store(ctx, data)
	if err != nil {
		l.deps.Logger.ErrorContext(ctx, "failed to store metadata to database", slog.Any("error", err))
		return model.ImageMeta{}, err
	}

	return data, nil
}

func (l *ImageLogic) GetImages(ctx context.Context, params model.ImageSearchParams, page pagination.Page) ([]model.ImageMeta, pagination.Metadata, error) {
	return l.imgRepo.Get(ctx, params, page)
}
