package store

import (
	"fmt"
	"strconv"
	"time"
)

type User struct {
	ID     string
	Author string
	Points int
}

type ServerConfig struct {
	ID      string
	CanPlay bool
	Tries   int
}

var (
	GetUserByID   map[string]*User
	GetServerByID map[string]*ServerConfig
)

func InitTime(ID string) {
	for range time.Tick(1 * time.Minute) {
		CheckMidnight(ID)
	}
}

func CheckMidnight(ID string) {
	currentTime := time.Now()
	if currentTime.Hour() == 0 && currentTime.Minute() == 0 {
		ResetServer(ID)
	}
}

func StoreServer(server ServerConfig) {
	GetServerByID[server.ID] = &server
	fmt.Printf("Server %s added to DB\n", server.ID)
}

func ServerExists(ID string) bool {
	return GetServerByID[ID] != nil
}

func StoreUser(user User) {
	GetUserByID[user.ID] = &user
	fmt.Printf("User %s added to DB\n", user.Author)
}

func UserExists(ID string) bool {
	return GetUserByID[ID] != nil
}

func AbleToPlay(ID string) bool {
	return GetServerByID[ID].CanPlay
}

func SetCanPlay(ID string, value bool) {
	GetServerByID[ID].CanPlay = value
}

func GetServerTries(ID string) int {
	return GetServerByID[ID].Tries
}

func SubtractServerTry(ID string) {
	GetServerByID[ID].Tries -= 1
}

func ResetServer(ID string) {
	GetServerByID[ID].CanPlay = true
	GetServerByID[ID].Tries = 6
}

func GivePoint(ID string) {
	fmt.Printf("User %s Awarded 1 point\n", GetUserByID[ID].Author)
	GetUserByID[ID].Points += 1
}

func GetPoints(ID string) string {
	return strconv.Itoa(GetUserByID[ID].Points)
}
