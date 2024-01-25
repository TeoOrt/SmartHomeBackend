package encryptors

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(encrypted), err
}

func DecryptPassword(db_pass string, in_pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(db_pass), []byte(in_pass))
	return err == nil
}
