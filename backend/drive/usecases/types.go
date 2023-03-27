package usecases

import (
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/google/uuid"
)

type FolderLink struct {
	Id   uuid.UUID
	Name string
}

type FileLink struct {
	Id    uuid.UUID
	Name  string
	Bytes int
}

type Children struct {
	Folders []FolderLink
	Files   []FileLink
}

func ToChildren(folders []*domain.Folder, files []*domain.File) Children {
	folderLinks := make([]FolderLink, len(folders))
	filesLinks := make([]FileLink, len(files))

	for i := 0; i < len(folders); i++ {
		folderLinks[i] = FolderLink{
			Id:   folders[i].Id(),
			Name: folders[i].Name(),
		}
	}
	for i := 0; i < len(files); i++ {
		filesLinks[i] = FileLink{
			Id:    files[i].Id(),
			Name:  files[i].Name(),
			Bytes: files[i].Size(),
		}
	}

	return Children{
		Folders: folderLinks,
		Files:   filesLinks,
	}
}
