package paradise

import "github.com/JohnArtR/cloudftp/db"
import "crypto/sha256"
import (
	"io"
	"log"
	"fmt"
	"strings"
)

type AuthSystem interface {
	CheckUser(userName, pass string, user *db.User) bool
}

type AuthManager struct {
	AuthSystem
}

type DefaultAuthSystem struct {
}

func (das DefaultAuthSystem) CheckUser(userName, pass string, user *db.User) bool {
	u, err := db.UserRepo.FindUser(userName)
	if err != nil {
		log.Printf(" [ERROR] %s", err)
		return false
	}
	enc := sha256.New()
	io.WriteString(enc, pass)
	encPass := fmt.Sprintf("%x", enc.Sum(nil))
	if strings.EqualFold(u.Password, encPass) {
		user = &u
		return true
	}
	return false
}

func NewDefaultAuthSystem() *AuthManager {
	am := AuthManager{}
	am.AuthSystem = DefaultAuthSystem{}
	return &am
}
