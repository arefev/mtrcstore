package repository

type Storage interface {
	Save(mType string, name string, value float64) error
	FindGauge(name string) (gauge, error)
	FindCounter(name string) (counter, error)
	Get() map[string]string
}