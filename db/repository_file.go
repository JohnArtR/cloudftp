package db

type FileFinder interface {
	FindFile(fileName string) (File, error)
	FindFiles (userName string) ([]File, error)
}

type FileSaver interface {
	SaveFile(file File) error
	SaveFiles(file []File) error
}

type FileRepository interface {
	FileFinder
	FileSaver
}
