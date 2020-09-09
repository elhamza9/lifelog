package auth

import (
	"errors"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// Service provides methods related to authentication
type Service struct {
	// name of the environment variable where the hash of the password is stored
	hashEnvVarName string
}

// NewService returns a new adding service with provided repository
func NewService(hashVarName string) Service {
	return Service{
		hashEnvVarName: hashVarName,
	}
}

// passwordMinLength specifies minimum length given password should be.
const passwordMinLength int = 8
const passwordMaxLength int = 256

// Errors
var (
	ErrPasswordLength       error = errors.New("Password must be " + strconv.Itoa(passwordMinLength) + " ~ " + strconv.Itoa(passwordMaxLength) + " characters")
	ErrIncorrectCredentials error = errors.New("Incorrect credentials")
	ErrHashNotFound         error = errors.New("Could not perform authentication")
)

// ValidatePassword validates password format
func (s Service) validatePassword(pass string) error {
	if len(pass) < passwordMinLength || len(pass) > passwordMaxLength {
		return ErrPasswordLength
	}
	return nil
}

// getPasswordHash returns bcrypt hash of correct password
func (s Service) getPasswordHash() string {
	return os.Getenv(s.hashEnvVarName)
}

// Authenticate hashes the given password and compares it with stored hash
func (s Service) Authenticate(pass string) error {
	if err := s.validatePassword(pass); err != nil {
		return err
	}
	hash := s.getPasswordHash()
	if hash == "" {
		return ErrHashNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)); err != nil {
		return ErrIncorrectCredentials
	}
	return nil
}
