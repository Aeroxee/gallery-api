package models

import (
	"time"

	"github.com/Aeroxee/gallery-api/auth"
)

// UserType implement for user type
type UserType int8

const (
	UserAdmin  UserType = iota // Type for user admin
	UserMember                 // type for user member
)

// User is model for implement user field in database.
type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"size:50" json:"first_name"`
	LastName  string    `gorm:"size:50" json:"last_name"`
	Username  string    `gorm:"size:50;uniqueIndex" json:"username"`
	Email     string    `gorm:"size:50;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:128" json:"-"`
	Avatar    *string   `gorm:"size:255" json:"avatar"`
	Type      UserType  `gorm:"default:1" json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Galleries []Gallery `gorm:"foreignKey:UserID" json:"galleries"`
}

func CreateNewUser(user *User) error {
	user.Password = auth.EncryptionPassword(user.Password)
	return DB().Create(user).Error
}

func GetUserByID(id int) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("id = ?", id).Preload("Galleries").First(&user).Error
	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("username = ?", username).Preload("Galleries").First(&user).Error
	return user, err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("email = ?", email).Preload("Galleries").First(&user).Error
	return user, err
}
