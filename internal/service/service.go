package service

import (
	"context"
	"tools-home/internal/dao"

	"github.com/google/wire"
)

var Provider = wire.NewSet(New)

// Service service.
type Service struct {
	dao dao.Dao
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		dao: d,
	}
	cf = s.Close
	return
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

// Ping ping the resource.
func (s *Service) Ping() error {
	return s.dao.Ping(context.Background())
}
