package image

import (
	"context"
	"image-app/internal/core"
	"image-app/internal/model"
	"image-app/internal/pkg/pagination"
	"log/slog"

	"github.com/huandu/go-sqlbuilder"
)

type ImageRepository struct {
	deps *core.Dependency
}

func NewImageRepository(deps *core.Dependency) *ImageRepository {
	return &ImageRepository{
		deps: deps,
	}
}

func (r *ImageRepository) Store(ctx context.Context, data model.ImageMeta) error {
	q := `INSERT INTO project.images (id, file_name, url, file_size, created_at) VALUES ($1, $2, $3, $4, now())`

	tx, err := r.deps.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, q, data.ID, data.FileName, data.URL, data.FileSize)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *ImageRepository) Get(ctx context.Context, params model.ImageSearchParams, page pagination.Page) ([]model.ImageMeta, pagination.Metadata, error) {
	var (
		result []model.ImageMeta
		meta   pagination.Metadata
	)

	// base query
	countQ := sqlbuilder.NewSelectBuilder()
	countQ.Select("COUNT(1)").From("project.images")

	q := sqlbuilder.NewSelectBuilder()
	q = q.Select("id", "file_name", "url", "file_size", "created_at").From("project.images")

	if params.FileName != "" {
		q.Where(
			q.ILike("file_name", "%"+params.FileName+"%"),
		)
	}

	if params.ID != "" {
		q.Where(
			q.ILike("id", "%"+params.ID+"%"),
		)
	}

	// pagination
	page.Compute()
	q.Limit(page.Limit)
	q.Offset(page.Offset)

	// build and exec query
	query, args := q.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.deps.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return result, meta, err
	}
	defer rows.Close()

	var temp model.SQLImageMeta
	for rows.Next() {
		err := rows.StructScan(&temp)
		if err != nil {
			r.deps.Logger.WarnContext(ctx, "failed to scan image metadata", slog.Any("error", err))
			continue
		}

		result = append(result, model.ImageMeta{
			ID:        temp.ID.String,
			FileName:  temp.FileName.String,
			URL:       temp.URL.String,
			FileSize:  temp.FileSize.Int64,
			CreatedAt: temp.CreatedAt.Time,
		})
	}

	// build metadata
	var total int64
	countQ.WhereClause = q.WhereClause
	query, args = countQ.BuildWithFlavor(sqlbuilder.PostgreSQL)
	err = r.deps.DB.QueryRowxContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return result, meta, err
	}

	meta.Compute(total, page.Size, page.Page)

	return result, meta, nil
}
