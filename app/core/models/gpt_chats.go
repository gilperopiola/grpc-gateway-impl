package models

import "time"

type GPTChat struct {
	ID        int          `gorm:"primaryKey" bson:"id"`
	Title     string       `gorm:"not null" bson:"title"`
	Messages  []GPTMessage `gorm:"foreignKey:ChatID" bson:"messages"`
	CreatedAt time.Time    `gorm:"autoCreateTime" bson:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" bson:"updated_at"`
}

func (GPTChat) TableName() string {
	return "gpt_chats"
}

type GPTMessage struct {
	ID        int       `gorm:"primaryKey" bson:"id"`
	ChatID    int       `gorm:"index;not null" bson:"chat_id"`
	Title     string    `gorm:"not null" bson:"title"`
	From      string    `gorm:"not null" bson:"from"`
	Content   string    `gorm:"type:text;not null" bson:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" bson:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" bson:"updated_at"`
}

func (GPTMessage) TableName() string {
	return "gpt_messages"
}
