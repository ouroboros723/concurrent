package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/totegamma/concurrent/client"
	"github.com/totegamma/concurrent/x/core"
	"github.com/totegamma/concurrent/x/util"
)

// Service is the interface for host service
type Service interface {
	Upsert(ctx context.Context, host core.Domain) (core.Domain, error)
	GetByFQDN(ctx context.Context, key string) (core.Domain, error)
	GetByCCID(ctx context.Context, key string) (core.Domain, error)
	List(ctx context.Context) ([]core.Domain, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, host core.Domain) error
	UpdateScrapeTime(ctx context.Context, id string, scrapeTime time.Time) error
}

type service struct {
	repository Repository
	config     util.Config
}

// NewService creates a new host service
func NewService(repository Repository, config util.Config) Service {
	return &service{repository, config}
}

// Upsert creates new host
func (s *service) Upsert(ctx context.Context, host core.Domain) (core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Service.Upsert")
	defer span.End()

	return s.repository.Upsert(ctx, host)
}

// GetByFQDN returns domain by FQDN
func (s *service) GetByFQDN(ctx context.Context, fqdn string) (core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Service.GetByFQDN")
	defer span.End()

	domain, err := s.repository.GetByFQDN(ctx, fqdn)
	if err == nil {
		if domain.DimensionID != s.config.Concurrent.Dimension {
			return core.Domain{}, fmt.Errorf("domain is not in the same dimension")
		}
		return domain, nil
	}

	domain, err = client.GetDomain(ctx, fqdn)
	if err != nil {
		return core.Domain{}, err
	}

	_, err = s.repository.Upsert(ctx, domain)
	if err != nil {
		return core.Domain{}, err
	}

	if domain.DimensionID != s.config.Concurrent.Dimension {
		return core.Domain{}, fmt.Errorf("domain is not in the same dimension")
	}

	return domain, nil
}

// GetByCCID returns domain by CCID
func (s *service) GetByCCID(ctx context.Context, key string) (core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Service.GetByCCID")
	defer span.End()

	return s.repository.GetByCCID(ctx, key)
}

// List returns list of domains
func (s *service) List(ctx context.Context) ([]core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Service.List")
	defer span.End()

	return s.repository.GetList(ctx)
}

// Delete deletes a domain
func (s *service) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "Domain.Service.Delete")
	defer span.End()

	return s.repository.Delete(ctx, id)
}

// Update updates a domain
func (s *service) Update(ctx context.Context, host core.Domain) error {
	ctx, span := tracer.Start(ctx, "Domain.Service.Update")
	defer span.End()

	return s.repository.Update(ctx, host)
}

// UpdateScrapeTime updates a domain's scrape time
func (s *service) UpdateScrapeTime(ctx context.Context, id string, scrapeTime time.Time) error {
	ctx, span := tracer.Start(ctx, "Domain.Service.UpdateScrapeTime")
	defer span.End()

	return s.repository.UpdateScrapeTime(ctx, id, scrapeTime)
}
