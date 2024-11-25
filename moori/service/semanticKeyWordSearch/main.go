package semanticKeyWordSearch

type SvcInterface interface {
	Search(keyword string) ([]interface{}, error)
}

type Service struct {
	service SvcInterface
}

func NewService(service SvcInterface) *Service {
	return &Service{
		service: service,
	}
}
func (s *Service) SearchByText(keyword string) ([]interface{}, error) {
	// TODO - additional service logic can be added here
	return s.service.Search(keyword)
}
