package models

import (
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Callsign  string
	Email     string
	Nfarl     bool
	Contactme bool
	Firsttime bool
	Youth     bool
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(dbfile string) (*UserService, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)
	//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	db, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic(err)
	}
	return &UserService{db: db}, nil
}

func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (us *UserService) DestructiveReset() error {
	if us.db.Migrator().HasTable(&User{}) {
		err := us.db.Migrator().DropTable(&User{})
		if err != nil {
			return err
		}
	}
	return us.db.AutoMigrate(&User{})
}

func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}
