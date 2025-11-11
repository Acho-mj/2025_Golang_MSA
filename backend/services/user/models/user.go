package models

import (
	"time"

	userpb "Acho-mj/2025_Golang_MSA/backend/gen/services/user/api"
)

type User struct {
	UserID    string    `dynamodbav:"user_id"`
	Email     string    `dynamodbav:"email"`
	Name      string    `dynamodbav:"name"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}

// ToProto: DB 모델(User) -> Proto 모델(*userpb.User)로 변환
func (u *User) ToProto() *userpb.User {
	if u == nil {
		return nil
	}
	return &userpb.User{
		UserId:    u.UserID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// UserFromProto: Proto 모델(*userpb.User) -> DB 모델(User)로 변환
func UserFromProto(p *userpb.User) *User {
	if p == nil {
		return nil
	}

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		createdAt = time.Time{}
	}

	return &User{
		UserID:    p.UserId,
		Email:     p.Email,
		Name:      p.Name,
		CreatedAt: createdAt,
	}
}
