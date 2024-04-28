// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package concurrent

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/totegamma/concurrent/client"
	"github.com/totegamma/concurrent/x/ack"
	"github.com/totegamma/concurrent/x/agent"
	"github.com/totegamma/concurrent/x/association"
	"github.com/totegamma/concurrent/x/auth"
	"github.com/totegamma/concurrent/x/domain"
	"github.com/totegamma/concurrent/x/entity"
	"github.com/totegamma/concurrent/x/jwt"
	"github.com/totegamma/concurrent/x/key"
	"github.com/totegamma/concurrent/x/message"
	"github.com/totegamma/concurrent/x/profile"
	"github.com/totegamma/concurrent/x/schema"
	"github.com/totegamma/concurrent/x/semanticid"
	"github.com/totegamma/concurrent/x/socket"
	"github.com/totegamma/concurrent/x/store"
	"github.com/totegamma/concurrent/x/subscription"
	"github.com/totegamma/concurrent/x/timeline"
	"github.com/totegamma/concurrent/x/userkv"
	"github.com/totegamma/concurrent/x/util"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func SetupJwtService(rdb *redis.Client) jwt.Service {
	repository := jwt.NewRepository(rdb)
	service := jwt.NewService(repository)
	return service
}

func SetupAckService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, config util.Config) ack.Service {
	repository := ack.NewRepository(db)
	service := SetupEntityService(db, rdb, mc, client2, config)
	keyService := SetupKeyService(db, rdb, mc, client2, config)
	ackService := ack.NewService(repository, client2, service, keyService, config)
	return ackService
}

func SetupKeyService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, config util.Config) key.Service {
	repository := key.NewRepository(db, mc)
	service := SetupEntityService(db, rdb, mc, client2, config)
	keyService := key.NewService(repository, service, config)
	return keyService
}

func SetupMessageService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, manager socket.Manager, config util.Config) message.Service {
	service := SetupSchemaService(db)
	repository := message.NewRepository(db, mc, service)
	entityService := SetupEntityService(db, rdb, mc, client2, config)
	timelineService := SetupTimelineService(db, rdb, mc, client2, manager, config)
	keyService := SetupKeyService(db, rdb, mc, client2, config)
	messageService := message.NewService(repository, client2, entityService, timelineService, keyService, config)
	return messageService
}

func SetupProfileService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, config util.Config) profile.Service {
	service := SetupSchemaService(db)
	repository := profile.NewRepository(db, mc, service)
	keyService := SetupKeyService(db, rdb, mc, client2, config)
	semanticidService := SetupSemanticidService(db)
	profileService := profile.NewService(repository, keyService, semanticidService)
	return profileService
}

func SetupAssociationService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, manager socket.Manager, config util.Config) association.Service {
	service := SetupSchemaService(db)
	repository := association.NewRepository(db, mc, service)
	entityService := SetupEntityService(db, rdb, mc, client2, config)
	timelineService := SetupTimelineService(db, rdb, mc, client2, manager, config)
	messageService := SetupMessageService(db, rdb, mc, client2, manager, config)
	keyService := SetupKeyService(db, rdb, mc, client2, config)
	associationService := association.NewService(repository, client2, entityService, timelineService, messageService, keyService, config)
	return associationService
}

func SetupTimelineService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, manager socket.Manager, config util.Config) timeline.Service {
	service := SetupSchemaService(db)
	repository := timeline.NewRepository(db, rdb, mc, client2, service, manager, config)
	entityService := SetupEntityService(db, rdb, mc, client2, config)
	domainService := SetupDomainService(db, client2, config)
	semanticidService := SetupSemanticidService(db)
	subscriptionService := SetupSubscriptionService(db)
	timelineService := timeline.NewService(repository, entityService, domainService, semanticidService, subscriptionService, config)
	return timelineService
}

func SetupDomainService(db *gorm.DB, client2 client.Client, config util.Config) domain.Service {
	repository := domain.NewRepository(db)
	service := domain.NewService(repository, client2, config)
	return service
}

func SetupEntityService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, config util.Config) entity.Service {
	service := SetupSchemaService(db)
	repository := entity.NewRepository(db, mc, service)
	jwtService := SetupJwtService(rdb)
	entityService := entity.NewService(repository, client2, config, jwtService)
	return entityService
}

func SetupSocketHandler(rdb *redis.Client, manager socket.Manager, config util.Config) socket.Handler {
	service := socket.NewService()
	handler := socket.NewHandler(service, rdb, manager)
	return handler
}

func SetupAgent(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, manager socket.Manager, config util.Config) agent.Agent {
	service := SetupStoreService(db, rdb, mc, client2, manager, config)
	agentAgent := agent.NewAgent(rdb, service, config)
	return agentAgent
}

func SetupAuthService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, config util.Config) auth.Service {
	service := SetupEntityService(db, rdb, mc, client2, config)
	domainService := SetupDomainService(db, client2, config)
	keyService := SetupKeyService(db, rdb, mc, client2, config)
	authService := auth.NewService(config, service, domainService, keyService)
	return authService
}

func SetupUserkvService(db *gorm.DB) userkv.Service {
	repository := userkv.NewRepository(db)
	service := userkv.NewService(repository)
	return service
}

func SetupSocketManager(mc *memcache.Client, db *gorm.DB, rdb *redis.Client, config util.Config) socket.Manager {
	manager := socket.NewManager(mc, rdb, config)
	return manager
}

func SetupSchemaService(db *gorm.DB) schema.Service {
	repository := schema.NewRepository(db)
	service := schema.NewService(repository)
	return service
}

func SetupStoreService(db *gorm.DB, rdb *redis.Client, mc *memcache.Client, client2 client.Client, manager socket.Manager, config util.Config) store.Service {
	repository := store.NewRepository(rdb)
	service := SetupKeyService(db, rdb, mc, client2, config)
	entityService := SetupEntityService(db, rdb, mc, client2, config)
	messageService := SetupMessageService(db, rdb, mc, client2, manager, config)
	associationService := SetupAssociationService(db, rdb, mc, client2, manager, config)
	profileService := SetupProfileService(db, rdb, mc, client2, config)
	timelineService := SetupTimelineService(db, rdb, mc, client2, manager, config)
	ackService := SetupAckService(db, rdb, mc, client2, config)
	subscriptionService := SetupSubscriptionService(db)
	storeService := store.NewService(repository, service, entityService, messageService, associationService, profileService, timelineService, ackService, subscriptionService, config)
	return storeService
}

func SetupSubscriptionService(db *gorm.DB) subscription.Service {
	service := SetupSchemaService(db)
	repository := subscription.NewRepository(db, service)
	subscriptionService := subscription.NewService(repository)
	return subscriptionService
}

func SetupSemanticidService(db *gorm.DB) semanticid.Service {
	repository := semanticid.NewRepository(db)
	service := semanticid.NewService(repository)
	return service
}

// wire.go:

// Lv0
var jwtServiceProvider = wire.NewSet(jwt.NewService, jwt.NewRepository)

var schemaServiceProvider = wire.NewSet(schema.NewService, schema.NewRepository)

var domainServiceProvider = wire.NewSet(domain.NewService, domain.NewRepository)

var semanticidServiceProvider = wire.NewSet(semanticid.NewService, semanticid.NewRepository)

var userKvServiceProvider = wire.NewSet(userkv.NewService, userkv.NewRepository)

// Lv1
var entityServiceProvider = wire.NewSet(entity.NewService, entity.NewRepository, SetupJwtService, SetupSchemaService)

var subscriptionServiceProvider = wire.NewSet(subscription.NewService, subscription.NewRepository, SetupSchemaService)

// Lv2
var keyServiceProvider = wire.NewSet(key.NewService, key.NewRepository, SetupEntityService)

var timelineServiceProvider = wire.NewSet(timeline.NewService, timeline.NewRepository, SetupEntityService, SetupDomainService, SetupSchemaService, SetupSemanticidService, SetupSubscriptionService)

// Lv3
var profileServiceProvider = wire.NewSet(profile.NewService, profile.NewRepository, SetupKeyService, SetupSchemaService, SetupSemanticidService)

var authServiceProvider = wire.NewSet(auth.NewService, SetupEntityService, SetupDomainService, SetupKeyService)

var ackServiceProvider = wire.NewSet(ack.NewService, ack.NewRepository, SetupEntityService, SetupKeyService)

// Lv4
var messageServiceProvider = wire.NewSet(message.NewService, message.NewRepository, SetupEntityService, SetupTimelineService, SetupKeyService, SetupSchemaService)

// Lv5
var associationServiceProvider = wire.NewSet(association.NewService, association.NewRepository, SetupEntityService, SetupTimelineService, SetupMessageService, SetupKeyService, SetupSchemaService)

// Lv6
var storeServiceProvider = wire.NewSet(store.NewService, store.NewRepository, SetupKeyService,
	SetupMessageService,
	SetupAssociationService,
	SetupProfileService,
	SetupEntityService,
	SetupTimelineService,
	SetupAckService,
	SetupSubscriptionService,
)

// Lv7
var agentServiceProvider = wire.NewSet(agent.NewAgent, SetupStoreService)
