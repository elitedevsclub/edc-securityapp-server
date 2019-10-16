package repos

import (
	"edc-security-app/services"
	"edc-security-app/types"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"web"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(opt *types.CreateUserOpts) (*types.User, *types.Otp, error) {
	user, err := types.NewUser(opt.Name, opt.Email, opt.Address, opt.Password)
	if err != nil {
		return nil, nil, err
	}
	user.Password = web.HashPassword(user.Password)
	tx := repo.db.Begin()
	if err := tx.Error; err != nil {
		return nil, nil, err
	}
	if err := tx.Create(user).Error; err != nil {
		return nil, nil, err
	}
	// creates OTP for this user
	otp := types.NewOtp(user.ID)
	if err := tx.Create(otp).Error; err != nil {
		return nil, nil, err
	}
	go func() {
		i := 0
		mr := &types.MailRequest{
			Email: user.Email,
			Title: "Security App OTP",
			User:  user.Name,
			Body:  fmt.Sprintf("Your authentication code is: %s", otp.Code),
		}
		// try send OTP 5 times
		for i == 5 {
			if err := services.SendEmail(mr); err == nil {
				break
			}
			i++
		}
	}()
	return user, otp, nil
}

func (repo *UserRepository) AuthenticateUser(opt *types.AuthenticateUserOpts) (*types.User, error) {
	user, err := repo.GetUserByAttr("email", opt.Email)
	if err != nil {
		return nil, err
	}
	if ok := web.VerifyPassword(user.Password, opt.Password); !ok {
		return nil, errors.New("invalid email and password combination")
	}
	return user, nil
}

func (repo *UserRepository) GetUserByAttr(attr string, value interface{}) (*types.User, error) {
	user := &types.User{}
	err := repo.db.Table("users").Where(attr+" = ?", value).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
