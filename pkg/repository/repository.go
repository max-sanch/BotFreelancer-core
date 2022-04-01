package repository

type Authentication interface {

}

type User interface {

}

type Repository struct {

}

func NewRepository() *Repository {
	return &Repository{}
}
