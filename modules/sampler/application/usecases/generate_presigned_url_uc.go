package samplerusecases

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	"github.com/google/uuid"
)

type GeneratePresignedURLUC struct {
	storageService commonports.StoragePort
}

func NewGeneratePresignedURLUC(
	storageService commonports.StoragePort,
) *GeneratePresignedURLUC {
	return &GeneratePresignedURLUC{
		storageService: storageService,
	}
}

func (uc *GeneratePresignedURLUC) Execute(
	ctx context.Context,
	userId string,
	keyType string,
	name string,
) (string, string, error) {
	ext := filepath.Ext(name)
	baseName := strings.TrimSuffix(name, ext)
	if ext == "" {
		ext = ".png"
	}

	uniqueName := fmt.Sprintf("%s_%s%s", baseName, uuid.NewString(), ext)
	r2Key := fmt.Sprintf("users/%s/samples/%s/%s", userId, keyType, uniqueName)

	contentType := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".webp":
		contentType = "image/webp"
	case ".gif":
		contentType = "image/gif"
	}

	url, err := uc.storageService.GeneratePResignedURL(ctx, r2Key, contentType)
	if err != nil {
		return "", "", err
	}

	return url, r2Key, nil
}
