package image

import (
	"image-app/internal/core"
	"image-app/internal/model"
	"image-app/internal/pkg/pagination"
	"image-app/internal/pkg/xhttp"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ImageHandler struct {
	deps     *core.Dependency
	imgLogic ImageLogicInterface
}

func NewImageHandler(deps *core.Dependency, imgLogic ImageLogicInterface) *ImageHandler {
	return &ImageHandler{
		deps:     deps,
		imgLogic: imgLogic,
	}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	file, header, err := r.FormFile("image")
	if err != nil {
		h.deps.Logger.WarnContext(ctx, "failed to read request body", slog.Any("error", err))
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to read request body for image upload",
		}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	if strings.ToLower(filepath.Ext(header.Filename)) != ".jpeg" {
		h.deps.Logger.WarnContext(ctx, "invalid image format/extension")
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Message: "invalid image format/extension",
		}, http.StatusBadRequest)
		return
	}

	if header.Size > 10*1024*1024 {
		h.deps.Logger.WarnContext(ctx, "invalid image size, must below 10 MB")
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Message: "invalid image size, must below 10 MB",
		}, http.StatusBadRequest)
		return
	}

	data, err := h.imgLogic.UploadImage(ctx, file, header.Filename, header.Size)
	if err != nil {
		h.deps.Logger.ErrorContext(ctx, "failed to upload image", slog.Any("error", err))
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to upload image",
		}, http.StatusInternalServerError)
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "image uploaded",
		Data:    data,
	}, http.StatusCreated)
}

func (h *ImageHandler) GetImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	fileName := chi.URLParam(r, "filename")

	page, err := pagination.ParsePaginationRequest(r)
	if err != nil {
		h.deps.Logger.WarnContext(ctx, "failed to parse pagination params", slog.Any("error", err))
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to parse pagination params",
		}, http.StatusInternalServerError)
		return
	}
	data, meta, err := h.imgLogic.GetImages(ctx, model.ImageSearchParams{
		ID:       id,
		FileName: fileName,
	}, page)
	if err != nil {
		h.deps.Logger.ErrorContext(ctx, "failed to get image(s)", slog.Any("error", err))
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to get image(s)",
		}, http.StatusInternalServerError)
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseListResponse{
		Message:  "images fetched",
		Data:     data,
		Metadata: meta,
	}, http.StatusOK)
}
