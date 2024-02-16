// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/totegamma/concurrent/x/agent"
	"github.com/totegamma/concurrent/x/association"
	"github.com/totegamma/concurrent/x/auth"
	"github.com/totegamma/concurrent/x/character"
	"github.com/totegamma/concurrent/x/collection"
	"github.com/totegamma/concurrent/x/domain"
	"github.com/totegamma/concurrent/x/entity"
	"github.com/totegamma/concurrent/x/jwt"
	"github.com/totegamma/concurrent/x/message"
	"github.com/totegamma/concurrent/x/socket"
	"github.com/totegamma/concurrent/x/stream"
	"github.com/totegamma/concurrent/x/userkv"
	"github.com/totegamma/concurrent/x/util"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func SetupMessageService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, manager socket.Manager, config util.Config) message.Service {
	repository := message.NewRepository(db, mc)
	streamRepository := stream.NewRepository(db, rdb, mc, manager, config)
	entityRepository := entity.NewRepository(db, mc)
	jwtRepository := jwt.NewRepository(rdb)
	service := jwt.NewService(jwtRepository)
	entityService := entity.NewService(entityRepository, config, service)
	streamService := stream.NewService(streamRepository, entityService, config)
	messageService := message.NewService(rdb, repository, streamService)
	return messageService
}

func SetupCharacterService(db *gorm.DB, mc *memcache.Client, config util.Config) character.Service {
	repository := character.NewRepository(db, mc)
	service := character.NewService(repository)
	return service
}

func SetupAssociationService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, manager socket.Manager, config util.Config) association.Service {
	repository := association.NewRepository(db, mc)
	streamRepository := stream.NewRepository(db, rdb, mc, manager, config)
	entityRepository := entity.NewRepository(db, mc)
	jwtRepository := jwt.NewRepository(rdb)
	service := jwt.NewService(jwtRepository)
	entityService := entity.NewService(entityRepository, config, service)
	streamService := stream.NewService(streamRepository, entityService, config)
	messageRepository := message.NewRepository(db, mc)
	messageService := message.NewService(rdb, messageRepository, streamService)
	associationService := association.NewService(repository, streamService, messageService)
	return associationService
}

func SetupStreamService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, manager socket.Manager, config util.Config) stream.Service {
	repository := stream.NewRepository(db, rdb, mc, manager, config)
	entityRepository := entity.NewRepository(db, mc)
	jwtRepository := jwt.NewRepository(rdb)
	service := jwt.NewService(jwtRepository)
	entityService := entity.NewService(entityRepository, config, service)
	streamService := stream.NewService(repository, entityService, config)
	return streamService
}

func SetupDomainHandler(db *gorm.DB, config util.Config) domain.Handler {
	repository := domain.NewRepository(db)
	service := domain.NewService(repository)
	handler := domain.NewHandler(service, config)
	return handler
}

func SetupEntityService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, config util.Config) entity.Service {
	repository := entity.NewRepository(db, mc)
	jwtRepository := jwt.NewRepository(rdb)
	service := jwt.NewService(jwtRepository)
	entityService := entity.NewService(repository, config, service)
	return entityService
}

func SetupSocketHandler(rdb *redis.Client, manager socket.Manager, config util.Config) socket.Handler {
	service := socket.NewService()
	handler := socket.NewHandler(service, rdb, manager)
	return handler
}

func SetupAgent(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, config util.Config) agent.Agent {
	repository := domain.NewRepository(db)
	service := domain.NewService(repository)
	entityRepository := entity.NewRepository(db, mc)
	jwtRepository := jwt.NewRepository(rdb)
	jwtService := jwt.NewService(jwtRepository)
	entityService := entity.NewService(entityRepository, config, jwtService)
	agentAgent := agent.NewAgent(rdb, config, service, entityService)
	return agentAgent
}

func SetupAuthService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, config util.Config) auth.Service {
	repository := auth.NewRepository(db)
	entityRepository := entity.NewRepository(db, mc)
	jwtRepository := jwt.NewRepository(rdb)
	service := jwt.NewService(jwtRepository)
	entityService := entity.NewService(entityRepository, config, service)
	domainRepository := domain.NewRepository(db)
	domainService := domain.NewService(domainRepository)
	authService := auth.NewService(repository, config, entityService, domainService)
	return authService
}

func SetupUserkvHandler(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, config util.Config) userkv.Handler {
	repository := userkv.NewRepository(rdb)
	service := userkv.NewService(repository)
	entityRepository := entity.NewRepository(db, mc)
	jwtRepository := jwt.NewRepository(rdb)
	jwtService := jwt.NewService(jwtRepository)
	entityService := entity.NewService(entityRepository, config, jwtService)
	handler := userkv.NewHandler(service, entityService)
	return handler
}

func SetupCollectionHandler(db *gorm.DB, rdb *redis.Client, config util.Config) collection.Handler {
	repository := collection.NewRepository(db)
	service := collection.NewService(repository)
	handler := collection.NewHandler(service)
	return handler
}

func SetupSocketManager(mc *memcache.Client, db *gorm.DB, rdb *redis.Client, config util.Config) socket.Manager {
	manager := socket.NewManager(mc, rdb, config)
	return manager
}

// wire.go:

var domainHandlerProvider = wire.NewSet(domain.NewHandler, domain.NewService, domain.NewRepository)

var userkvHandlerProvider = wire.NewSet(userkv.NewHandler, userkv.NewService, userkv.NewRepository)

var collectionHandlerProvider = wire.NewSet(collection.NewHandler, collection.NewService, collection.NewRepository)

var jwtServiceProvider = wire.NewSet(jwt.NewService, jwt.NewRepository)

var entityServiceProvider = wire.NewSet(entity.NewService, entity.NewRepository, jwtServiceProvider)

var streamServiceProvider = wire.NewSet(stream.NewService, stream.NewRepository, entityServiceProvider)

var messageServiceProvider = wire.NewSet(message.NewService, message.NewRepository, streamServiceProvider)

var associationServiceProvider = wire.NewSet(association.NewService, association.NewRepository, messageServiceProvider)

var characterServiceProvider = wire.NewSet(character.NewService, character.NewRepository)
