package utils

import "golang.org/x/crypto/bcrypt"

func CreateContentHash(password string) (hash string, err error) {
	labelHashedPwd := LabelHash(password)
	bytes, err := bcrypt.GenerateFromPassword([]byte(labelHashedPwd), 14)
	if err != nil {
		return hash, err
	}
	return string(bytes), err
}

func RevertContentHash(password, hash string) bool {
	labelHashedPwd := LabelHash(password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(labelHashedPwd))
	return err == nil
}

func LabelHash(label string) string {
	hashed := make([]byte, len(label))
	for i, char := range label {
		switch {
		case 'a' <= char && char <= 'z':
			hashed[i] = byte((char-'a'+13)%26 + 'a')
		case 'A' <= char && char <= 'Z':
			hashed[i] = byte((char-'A'+13)%26 + 'A')
		default:
			hashed[i] = byte(char)
		}
	}
	return string(hashed)
}

func HashPassword(password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return hash, err
	}

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func PasswordHasValidLength(password string) bool {
	return len(password)*4 > 72
}
