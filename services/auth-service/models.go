package main

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore struct {
	users  map[string]User
	emails map[string]string
}

var store *UserStore

func InitStore() {
	store = &UserStore{
		users:  make(map[string]User),
		emails: make(map[string]string),
	}
}

func SaveUser(user User) error {
	if _, exists := store.emails[user.Email]; exists {
		return ErrUserExists
	}
	store.users[user.ID] = user
	store.emails[user.Email] = user.ID
	return nil
}

func FindUserByEmail(email string) (User, error) {
	if id, exists := store.emails[email]; exists {
		return store.users[id], nil
	}
	return User{}, ErrUserNotFound
}

func FindUserByID(id string) (User, error) {
	if user, exists := store.users[id]; exists {
		return user, nil
	}
	return User{}, ErrUserNotFound
}

var (
	ErrUserExists   = &AuthError{"user already exists"}
	ErrUserNotFound = &AuthError{"user not found"}
)

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}
