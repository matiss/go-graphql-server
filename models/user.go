package models

import (
	"golang.org/x/crypto/bcrypt"
	"net"
	"regexp"
	"time"
	"unicode"
)

type User struct {
	ID         int32     `pg:"id" json:"id"`
	Email      string    `pg:"email" json:"email" binding:"required"`
	Password   string    `pg:"password" json:"-"`
	Number     string    `pg:"number" json:"number" binding:"required" validate:"min=8,max=12"`
	Name       string    `pg:"name" json:"name"`
	LoginTime  time.Time `pg:"login_at" json:"-"`
	LoginIP    net.IP    `pg:"login_ip" json:"-"`
	LoginCount int64     `pg:"login_count" json:"-"`
	Status     int       `pg:"status" json:"status"`
	Role       int       `pg:"role" json:"role,omitempty"`
	CreatedAt  time.Time `pg:"created_at" json:"createdAt,omitempty"`
	UpdatedAt  time.Time `pg:"updated_at" json:"updatedAt,omitempty"`
	DeletedAt  time.Time `pg:"deleted_at" json:"-"`
}

func (u *User) Active() bool {
	return u.DeletedAt.IsZero()
}

func (user *User) HashedPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

func (user *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (user *User) ValidateEmail() bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(user.Email)
}

// Passwords must contain at least eight characters, including at least 1 letter and 1 number.
func (user *User) ValidatePassword() bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)

	if len(user.Password) >= 8 {
		hasMinLen = true
	}

	for _, char := range user.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasMinLen && (hasUpper || hasLower) && hasNumber
}
