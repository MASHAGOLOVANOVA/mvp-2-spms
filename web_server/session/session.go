package session

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var Sessions = map[string]Session{}
var BotToken = ""

const SessionDefaultExpTime = 7 * 24 * time.Hour

type Session struct {
	user   UserInfo
	expiry time.Time
}

func (s Session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s Session) GetUser() UserInfo {
	return s.user
}

func InitSession(user UserInfo, exp time.Time) Session {
	return Session{
		user:   user,
		expiry: exp,
	}
}

type UserInfo struct {
	username  string
	accountId string
	isProf    bool
}

func InitUserInfo(username string, profId string, isProf bool) UserInfo {
	return UserInfo{
		username:  username,
		accountId: profId,
		isProf:    isProf,
	}
}

func (u UserInfo) GetUsername() string {
	return u.username
}
func (u UserInfo) GetAccId() string {
	return u.accountId
}
func (u UserInfo) IsProf() bool { return u.isProf }

func SetBotTokenFromJson(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to find JSON file for bot token")
	}
	bot := struct {
		TelegramToken string `json:"telegram_bot_token"`
	}{}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&bot); err != nil {
		log.Fatal("unable to decode")
	}

	BotToken = bot.TelegramToken
}
