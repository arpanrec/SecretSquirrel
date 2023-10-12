package physical

type Storage interface {
	GetData() (string, error)
	PutData() (bool, error)
	DeleteData() error
}
