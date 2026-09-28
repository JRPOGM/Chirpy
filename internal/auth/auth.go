package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	//argon2id.CreateHash(string, params) always returns a hash string and an error
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err = argoin2id.ComparePasswordAndHash(password, hash)
	//argoin2id.ComparePasswordAndHash(string, string) always returns a bool and an error
	if err != nil {
		return false, err
	}
	return match, nil
}

