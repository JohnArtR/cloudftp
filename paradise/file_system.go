package paradise

import "github.com/JohnArtR/cloudftp/db"

type FileSystem interface {
	GetFiles(user *db.User) ([]map[string]string, error)
}

type FileManager struct {
	FileSystem
}

type DefaultFileSystem struct {
}

func (dfs DefaultFileSystem) GetFiles(user *db.User) ([]map[string]string, error) {
	files := make([]map[string]string, 0)

	//if p.user == "test" {
	// no op just to use p.user as example
	//}

	for i := 0; i < 5; i++ {
		file := make(map[string]string)
		file["size"] = "90210"
		file["name"] = "paradise.txt"
		files = append(files, file)
	}

	return files, nil
}

func NewDefaultFileSystem() *FileManager {
	fm := FileManager{}
	fm.FileSystem = DefaultFileSystem{}
	return &fm
}
