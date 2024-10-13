package memory

import (
	"cloud-logic/database"
	"encoding/base64"
	"log"

	"github.com/google/uuid"
)

type Collection map[string]interface{}

type User struct {
	ID          string     `json:"id"`
	PW          string     `json:"pw"`
	Email       string     `json:"email"`
	AllowIP     []string   `json:"allow-ip"`
	AllowID     []string   `json:"allow-id"`
	API_TOKEN   string     `json:"api-token"`
	Collection  Collection `json:"collect"`
	Certi       bool
	Certi_Token string
}

type Storage struct {
	Path string           `json:"path"`
	User map[string]*User `json:"user"`
}

func load(storage string) map[string]*User {
	user := make(map[string]*User)
	files, err := database.ListFiles(storage)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range files {
		var parse *User
		err := database.ParseJSON(file, &parse)
		if err != nil {
			log.Fatalln(err)
		}
		user[parse.ID] = parse
	}
	return user
}

func CheckKey[T comparable, V any](m map[T]V, k T) bool {
	if _, exists := m[k]; exists {
		return true
	} else {
		return false
	}
}

func NewStorage(path string) *Storage {
	return &Storage{
		Path: path,
		User: load(path),
	}
}

func NewToken() string {
	return base64.RawURLEncoding.EncodeToString([]byte(uuid.New().String()))
}

func (s *Storage) NewUser(token, email, id, pw string) *User {
	password, err := Hash(pw)
	if err != nil || s.CheckUser(id) {
		return nil
	}
	user := &User{
		ID:          id,
		PW:          password,
		Email:       email,
		AllowIP:     make([]string, 0),
		AllowID:     make([]string, 0),
		API_TOKEN:   NewToken(),
		Collection:  make(Collection),
		Certi:       false,
		Certi_Token: token,
	}
	s.User[id] = user
	data, err := database.ByteJSON(user)
	if err != nil {
		return nil
	}
	err = database.WriteFile(database.FileName(s.Path, id+".json"), data)
	if err != nil {
		return nil
	}
	return user
}

func (s *Storage) DeleteUser(id string) {
	if !s.CheckUser(id) {
		return
	}
	delete(s.User, id)
	database.DeleteFile(database.FileName(s.Path, id+".json"))
}

func (s *Storage) CheckUser(id string) bool {
	return CheckKey(s.User, id)
}

func (s *Storage) GetUser(id string) *User {
	if s.CheckUser(id) {
		return s.User["id"]
	}
	return nil
}

func (s *Storage) CheckAPIToken(token string) *User {
	for _, v := range s.User {
		if v.API_TOKEN == token {
			return v
		}
	}
	return nil
}

func (s *Storage) CheckEmail(email string) *User {
	for _, v := range s.User {
		if v.Email == email {
			return v
		}
	}
	return nil
}

func (s *Storage) Backup() {
	for _, user := range s.User {
		data, err := database.ByteJSON(user)
		if err != nil {
			continue
		}
		err = database.WriteFile(database.FileName(s.Path, user.ID+".json"), data)
		if err != nil {
			continue
		}
	}
}
