package ack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/totegamma/concurrent/client"
	"github.com/totegamma/concurrent/x/core"
	"github.com/totegamma/concurrent/x/entity"
	"github.com/totegamma/concurrent/x/key"
	"github.com/totegamma/concurrent/x/util"
)

// Service is the interface for entity service
type Service interface {
	Ack(ctx context.Context, from, to string) error
	GetAcker(ctx context.Context, key string) ([]core.Ack, error)
	GetAcking(ctx context.Context, key string) ([]core.Ack, error)
}

type service struct {
	repository Repository
	entity     entity.Service
	key        key.Service
	config     util.Config
}

// NewService creates a new entity service
func NewService(repository Repository, entity entity.Service, key key.Service, config util.Config) Service {
	return &service{
		repository,
		entity,
		key,
		config,
	}
}

// Ack creates new Ack
func (s *service) Ack(ctx context.Context, document string, signature string) error {
	ctx, span := tracer.Start(ctx, "ServiceAck")
	defer span.End()

	var doc core.AckDocument
	err := json.Unmarshal([]byte(document), &doc)
	if err != nil {
		span.RecordError(err)
		return err
	}

	switch doc.Type {
	case "ack":
		address, err := s.entity.GetAddress(ctx, doc.To)
		if err == nil {
			packet := core.Commit{
				Document:  document,
				Signature: signature,
			}

			packetStr, err := json.Marshal(packet)
			if err != nil {
				span.RecordError(err)
				return err
			}

			resp, err := client.Commit(ctx, address.Domain, string(packetStr))
			if err != nil {
				span.RecordError(err)
				return err
			}

			defer resp.Body.Close()
		}

		return s.repository.Ack(ctx, &core.Ack{
			From:      doc.From,
			To:        doc.To,
			Document:  document,
			Signature: signature,
		})
	case "unack":
		address, err := s.entity.GetAddress(ctx, doc.To)
		if err == nil {

			packet := core.Commit{
				Document:  document,
				Signature: signature,
			}

			packetStr, err := json.Marshal(packet)
			if err != nil {
				span.RecordError(err)
				return err
			}

			resp, err := client.Commit(ctx, address.Domain, string(packetStr))
			if err != nil {
				span.RecordError(err)
				return err
			}
			defer resp.Body.Close()
		}
		return s.repository.Unack(ctx, &core.Ack{
			From:      doc.From,
			To:        doc.To,
			Document:  document,
			Signature: signature,
		})
	default:
		return fmt.Errorf("invalid object type")
	}
}

// GetAcker returns acker
func (s *service) GetAcker(ctx context.Context, user string) ([]core.Ack, error) {
	ctx, span := tracer.Start(ctx, "ServiceGetAcker")
	defer span.End()

	return s.repository.GetAcker(ctx, user)
}

// GetAcking returns acking
func (s *service) GetAcking(ctx context.Context, user string) ([]core.Ack, error) {
	ctx, span := tracer.Start(ctx, "ServiceGetAcking")
	defer span.End()

	return s.repository.GetAcking(ctx, user)
}
