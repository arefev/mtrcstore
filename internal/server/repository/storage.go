package repository

type Storage interface {
	Save(mType string, name string, value string) error
}