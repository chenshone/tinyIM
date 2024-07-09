package tools

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

const SessionPrefix = "session_"

func GetRandomToken(length int) string {
	r := make([]byte, length)
	_, _ = io.ReadFull(rand.Reader, r)
	return base64.URLEncoding.EncodeToString(r)
}

func CreateSessionId(data string) string {
	return SessionPrefix + data
}

func GetSessionIdByUserId(userId int) string {
	return fmt.Sprintf("session_map_%d", userId)
}

func GetSessionName(sessionId string) string {
	return SessionPrefix + sessionId
}

func HashWithSalt(data string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func CompareHashWithSaltAndPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetNowDateTime() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
}
