package types

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type User struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email"`
	Address string `json:"address"`
	Password string `json:"-"`
	Verified bool `json:"verified"`
}

type Otp struct {
	gorm.Model
	CodeIdentifier string `json:"code_identifier"`
	Code string `json:"code"`
	UserId uint `json:"user_id"`
}

func NewOtp(userId uint) *Otp {
	code := fmt.Sprintf("%06d", rand.Intn(999999))
	return &Otp{
		CodeIdentifier: randString(),
		Code:           code,
		UserId: userId,
	}
}

func NewUser(name, email, address, password string) (*User, error) {
	if !IsLautechEmail(email) {
		return nil, errors.New(fmt.Sprintf("%s is not a LAUTECH email", email))
	}
	if len(password) < 5 {
		return nil, errors.New("password is too soft")
	}
	return &User{
		Name:     name,
		Email:    email,
		Address:  address,
		Password: password,
		Verified: false,
	}, nil
}
type CreateUserOpts struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
	Address string `json:"address"`
}

type AuthenticateUserOpts struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func IsLautechEmail(email string) bool {
	// i can't write a regex. This'll work just fine
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[1] != "student.lautech.edu.ng" {
		return false
	}
	return true
}

func randString() string {
	m5 := md5.New()
	m5.Write([]byte(time.Now().String()))
	return fmt.Sprintf("%x", m5.Sum(nil))
}


type MailRequest struct {
	Email, Title, User, Body string
}

