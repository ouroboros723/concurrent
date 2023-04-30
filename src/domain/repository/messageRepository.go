package repository

import (
    "gorm.io/gorm"
    "concurrent/domain/model"
)

type IMessageRepository interface {
    Create(message model.Message)
    GetAll() []model.Message
    GetFollowee() []model.Message
}

type MessageRepository struct {
    db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
    return MessageRepository{db: db}
}

func (r *MessageRepository) Create(message model.Message) {
    r.db.Create(&message)
}

func (r *MessageRepository) GetAll() []model.Message {
    var messages []model.Message
    r.db.Find(&messages)
    return messages
}

func (r *MessageRepository) GetFollowee(followeeList []string) []model.Message {
    var messages []model.Message
    r.db.Where("author = ANY($1)", followeeList).Find(&messages)
    return messages
}
