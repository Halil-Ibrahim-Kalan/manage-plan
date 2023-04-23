package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `json:"Username" form:"Username" query:"Username"`
	Password string `json:"Password" form:"Password" query:"Password"`
}

func CreateUser(db *gorm.DB, User *User) (err error) {
	err = db.Table("users").Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserById(db *gorm.DB, User *User, id int) (err error) {
	err = db.Table("users").Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckUserByUsername(db *gorm.DB, username string) (bool, error) {
	var User User
	result := db.Table("users").Where("username = ?", username).First(User).Error
	if result != nil {
		return false, result
	}
	return true, nil
}

func CheckUserPass(db *gorm.DB, User *User, username, password string) error {
	err := db.Table("users").Where("username = ?", username).Where("password = ?", password).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(db *gorm.DB, User *User) (err error) {
	db.Table("users").Save(User)
	return nil
}

func DeleteUser(db *gorm.DB, User *User, id int) (err error) {
	db.Table("users").Where("id = ?", id).Delete(User)
	return nil
}
