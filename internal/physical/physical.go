package physical

type KVData struct {
	Value    string            `json:"value"`
	Metadata map[string]string `json:"metadata"`
}

type KeyValueStorage interface {
	ListKeys(key *string) ([]string, error)

	ListVersions(key *string) ([]int, error)

	GetLatestVersion(key *string) (int, error)

	GetNextVersion(key *string) (int, error)

	Get(key *string, version *int) (*KVData, error)

	Save(key *string, keyValue *KVData) error

	Update(key *string, keyValue *KVData, version *int) error

	Delete(key *string, version *int) error
}
