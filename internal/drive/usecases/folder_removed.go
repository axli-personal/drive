package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FolderRemovedResult struct {
	CanAcknowledge bool
}

type folderRemovedHandler struct {
	folderRepo repository.FolderRepository
	fileRepo   repository.FileRepository
	eventRepo  repository.EventRepository
}

func (handler folderRemovedHandler) Handle(
	ctx context.Context,
	event events.FolderRemoved,
) (FolderRemovedResult, error) {
	folderId, err := uuid.Parse(event.FolderId)
	if err != nil {
		return FolderRemovedResult{CanAcknowledge: true}, err
	}

	parent, err := domain.CreateFolderParent(folderId)
	if err != nil {
		return FolderRemovedResult{}, err
	}

	files, err := handler.fileRepo.FindFile(
		ctx,
		repository.FindFileOptions{
			Parent: parent,
		},
	)
	if err != nil {
		if err != repository.ErrNotFound {
			return FolderRemovedResult{CanAcknowledge: false}, err
		}
	}
	for _, file := range files {
		file.ChangeState(domain.StateTrashed)
		handler.fileRepo.UpdateFile(
			ctx,
			file,
			repository.UpdateFileOptions{},
		)
	}

	folders, err := handler.folderRepo.FindFolder(
		ctx,
		repository.FindFolderOptions{
			Parent: parent,
		},
	)
	if err != nil {
		if err != repository.ErrNotFound {
			return FolderRemovedResult{CanAcknowledge: false}, err
		}
	}
	for _, folder := range folders {
		folder.ChangeState(domain.StateTrashed)
		handler.folderRepo.UpdateFolder(
			ctx,
			folder,
			repository.UpdateFolderOptions{},
		)
		handler.eventRepo.PublishFolderRemoved(
			ctx,
			events.FolderRemoved{
				FolderId: folder.Id().String(),
			},
		)
	}

	return FolderRemovedResult{CanAcknowledge: true}, nil
}

type FolderRemovedHandler decorator.Handler[events.FolderRemoved, FolderRemovedResult]

func NewFolderRemovedHandler(
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) FolderRemovedHandler {
	return decorator.WithLogging[events.FolderRemoved, FolderRemovedResult](
		folderRemovedHandler{
			fileRepo: fileRepo,
		},
		logger,
	)
}
