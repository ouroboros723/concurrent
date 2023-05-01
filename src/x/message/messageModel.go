package message

import (
    "time"
)

type Message struct {
    ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Author string `json:"author" gorm:"type:varchar(64)"`
    Schema string `json:"schema" gorm:"type:varchar(1024)"`
    Payload string `json:"payload" gorm:"type:json"`
    Signature string `json:"signature" gorm:"type:char(130)"`
    CDate time.Time `json:"cdate" gorm:"type:timestamp with time zone;not null;default:clock_timestamp()"`
    Associations []string `json:"associations" gorm:"type:uuid[]"`
    Streams string `json:"streams" gorm:"type:text"`
}

type MessagesResponse struct {
    Messages []Message `json:"messages"`
}

type MessageResponse struct {
    Message Message `json:"message"`
}

