package repository

type Storage interface {
	Save(mType string, name string, value float64) error
	Find(mType string, name string) (float64, error)
	Get() map[string]float64
}
