package tools

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

const SessionPrefix = "session_"

func GetSnowflakeId() string {
	//default node id eq 1,this can modify to different serverId node
	node, _ := snowflake.NewNode(1)
	// Generate a snowflake ID.
	id := node.Generate().String()
	return id
}

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
