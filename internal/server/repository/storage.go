package repository

type Storage interface {
	Save(mType string, name string, value float64) error
}
