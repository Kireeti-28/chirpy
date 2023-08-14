package database

import (
	"errors"

	"github.com/kireeti-28/chirpy/internal/auth"
)

type UserDB struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserReq struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(userReq UserReq) (UserDB, error) {
	if db.userExists(userReq.Email) == nil {
		return UserDB{}, errors.New("email already exists")
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return UserDB{}, err
	}

	hash, err := auth.HashPassword(userReq.Password)
	if err != nil {
		return UserDB{}, err
	}

	id := len(dbStructure.Users) + 1
	newUser := UserDB{
		ID:       id,
		Email:    userReq.Email,
		Password: hash,
		IsChirpyRed: false,
	}

	dbStructure.Users[id] = newUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserDB{}, err
	}

	return newUser, nil
}

func (db *DB) GetUserById(id int) (UserDB, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserDB{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return UserDB{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) userExists(email string) error {
	_, err := db.GetUserByEmail(email)

	return err
}

func (db *DB) GetUserByEmail(email string) (UserDB, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserDB{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return UserDB{}, ErrNotExist
}

func (db *DB) UpdateUser(userId int, updateUser UserReq) (UserDB, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserDB{}, err
	}

	user, ok := dbStructure.Users[userId]
	if !ok {
		return UserDB{}, errors.New("User not found")
	}

	user.Email = updateUser.Email
	hash, err := auth.HashPassword(updateUser.Password)
	if err != nil {
		return UserDB{}, err
	}
	user.Password = hash

	dbStructure.Users[userId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserDB{}, err
	}

	return user, nil
}

func (db *DB) UpdateMemberShipRed(userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[userId]
	if !ok {
		return errors.New("invalid user id")
	}

	user.IsChirpyRed = true
	dbStructure.Users[userId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
