package db

type PostgresRepository struct {
	Host string
	Port string
	DBName string
}

func (repo PostgresRepository) FindUser(userName string) (*User, error) {
	return &User{}, nil
}

func (repo PostgresRepository) FindFile(fileName string) (File, error) {
	return File{}, nil
}

func (repo PostgresRepository) FindFiles(userName string) ([]File, error) {
	return []File{}, nil
}

func (repo PostgresRepository) SaveUser(user User) error {
	return nil
}

func (repo PostgresRepository) SaveFile(file File) error {
	return nil
}
