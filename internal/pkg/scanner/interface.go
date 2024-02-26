package scanner

type Repository interface {
	GetRequest(id int) ([]byte, error)
	GetResponse(id int) ([]byte, error)
	GetAllRequests() ([][]byte, error)
}
