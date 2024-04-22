package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 8
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params *CreateUserParams) Valide() map[string]string {
	errs := make(map[string]string)
	if len(params.FirstName) < minFirstNameLen {
		errs["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errs["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errs["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errs["email"] = "email is invalid"
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
