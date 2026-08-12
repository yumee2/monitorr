package api

import (
	"context"
	"monitorr/internal/storage"
	"time"
)

type CheckReader interface {
	FindAll(ctx context.Context, interval time.Duration) (map[string][]storage.ServiceStatus, error)
	CountUptime(ctx context.Context, serviceName string, interval time.Duration) (int, error)
	FindTotalByName(ctx context.Context, serviceName string, interval time.Duration) ([]storage.ServiceStatus, error)
}

type Server struct {
	repo CheckReader
}

func NewServer(storage CheckReader) *Server {
	return &Server{repo: storage}
}

func (s *Server) GetAllServices(ctx context.Context, interval time.Duration) (map[string][]storage.ServiceStatus, error) {
	return s.repo.FindAll(ctx, interval)
}

func (s *Server) GetUptime(ctx context.Context, serviceName string, interval time.Duration) (int, error) {
	return s.repo.CountUptime(ctx, serviceName, interval)
}

func (s *Server) GetTotalByName(ctx context.Context, serviceName string, interval time.Duration) ([]storage.ServiceStatus, error) {
	return s.repo.FindTotalByName(ctx, serviceName, interval)
}
