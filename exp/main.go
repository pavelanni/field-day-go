package main

import (
	"field-day/models"
	"fmt"

	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "fieldday_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Age   int
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(dsn)
	if err != nil {
		panic(err)
	}
	us.DestructiveReset()

	user := models.User{
		Name:  "Michael Zhabitsky",
		Email: "jabitsky@mail.ru",
	}

	if err = us.Create(&user); err != nil {
		panic(err)
	}

	user.Name = "Updated Name"
	if err := us.Update(&user); err != nil {
		panic(err)
	}

	foundUser, err := us.ByEmail("jabitsky@mail.ru")
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser.Name)

	if err := us.Delete(foundUser.ID); err != nil {
		panic(err)
	}

	// Verify the user is deleted
	_, err = us.ByID(foundUser.ID)
	if err != models.ErrNotFound {
		panic("user was not deleted!")
	}
}
