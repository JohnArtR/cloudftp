package storage

type FileManager interface {
	Save(fileContent []byte, dstFilePath string) error
}

type FileService FileManager