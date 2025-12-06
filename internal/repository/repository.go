package repository

type User struct {
    ID    string
    Login string
    Hash  string
}

type Repository struct{}

func NewRepository() *Repository { 
	return &Repository{} 
}

