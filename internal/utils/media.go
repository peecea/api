package utils

import "github.com/gabriel-vasile/mimetype"

func IsValidImage(mType string) bool {
	allowed := []string{"image/png", "image/jpeg"}
	return mimetype.EqualsAny(mType, allowed...)
}

func IsValidDocument(mType string) bool {
	allowed := []string{"application/msword", "application/pdf"}
	return mimetype.EqualsAny(mType, allowed...)

}

func IsValidVideo(mType string) bool {
	allowed := []string{"video/mpeg", "video/mp4"}
	return mimetype.EqualsAny(mType, allowed...)
}

func IsValidFile(mType string) bool {
	allowed := []string{"text/plain", "application/png", "application/jpg", "application/word", "application/pdf"}
	return mimetype.EqualsAny(mType, allowed...)
}
