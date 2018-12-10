package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             int64     `json:"id" form:"id"`
	UserName       string    `json: "username" form:"username"`
	PasswordEncryp []byte    `json:"-" form:"-"`
	Password       string    `json:"-" form:"-"`
	Email          string    `json:"email" form:"email"`
	CreatedAT      time.Time `json:"create_at" form:"create_at"`
	Admin          bool
}

//verify if a ID is valid
func (u *User) isValid() bool {
	return u.ID > 0
}

//encryt a password
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

//validy if the password is correct
func ValidatePassword(userPassword string, passwordEncrypt []byte) error {
	return bcrypt.CompareHashAndPassword(passwordEncrypt, []byte(userPassword))
}
func GenerateUser() []User {
	t := time.Now()
	pwd := "12345"
	pwdEncrypt, _ := GeneratePassword(pwd)
	users := []User{{ID: 1, UserName: "Osiris", PasswordEncryp: pwdEncrypt, Password: pwd, Email: "osiris@gmail.com", CreatedAT: t, Admin: true},
		{ID: 2, UserName: "Cleopatra", PasswordEncryp: pwdEncrypt, Password: pwd, Email: "cleopatra@gmail.com", CreatedAT: t, Admin: false},
		{ID: 3, UserName: "Seth", PasswordEncryp: pwdEncrypt, Password: pwd, Email: "seth@gmail.com", CreatedAT: t, Admin: false},
		{ID: 4, UserName: "Minerva", PasswordEncryp: pwdEncrypt, Password: pwd, Email: "minerva@gmail.com", CreatedAT: t, Admin: false}}
	return users
}
func (u *User) GetEmailsToSend(users []User) (emailTo []string) {
	for _, user := range users {
		if user.ID != u.ID {
			emailTo = append(emailTo, user.Email)
		}
	}
	return
}
