package physical

type Storage interface {
	GetData(l string) (string, error)
	PutData(l string, d string) (bool, error)
	DeleteData(l string) error
}
