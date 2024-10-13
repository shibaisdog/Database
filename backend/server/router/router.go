package router

import (
	"cloud-logic/database"
	"cloud-logic/database/memory"
	"cloud-logic/protocol"
	js "encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
)

type Client struct {
	User      *memory.User
	Login     bool
	API_TOKEN string
}

type JSON struct {
	Action string      `json:"action"`
	Value  interface{} `json:"value"`
}

type User_Login struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}

type User_Singup struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Pw    string `json:"pw"`
}

func Router(smtp *protocol.Smtp, db *memory.Storage, request_info *Client, request net.Conn, butter []byte) {
	var json JSON
	err := ParseJSON(butter, &json)
	if err != nil {
		go SendErrorMessage(request, http.StatusBadRequest, "Data is invalid | Need : {\n\t\"action\":\"\",\n\t\"value\":\"\"\n}")
		return
	}
	if (json.Action != "certi" && json.Action != "login-user-email" && json.Action != "login" && json.Action != "singup" && !request_info.Login) || json.Value == "" {
		go SendErrorMessage(request, http.StatusBadRequest, "Login required | Need : {\n\t\"action\":\"login\",\n\t\"value\":\"YOUR-API-TOKEN\"\n}")
		return
	}
	if request_info.Login && json.Action != "certi" && !request_info.User.Certi {
		go SendErrorMessage(request, http.StatusBadRequest, "Email verification required | Need : {\n\t\"action\":\"certi\",\n\t\"value\":{\"email\":\"Your-Email\",\"value\":\"Code\"}}")
		return
	}
	if json.Action == "certi" {
		if !request_info.Login {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "must login first | Need : {\n\t\"action\":\"certi\",\n\t\"value\":{\"email\":\"Your-Email\",\"value\":\"Code\"}}")
			return
		} else {
			if request_info.User.Certi_Token != json.Value.(string) {
				go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Code does not match | Need : {\n\t\"action\":\"certi\",\n\t\"value\":{\"email\":\"Your-Email\",\"value\":\"Code\"}}")
				return
			}
			userdata := db.CheckEmail(request_info.User.Email)
			userdata.Certi = true
			request_info.User = userdata
			go SendMessage(request, http.StatusOK, "Certi successfully")
		}
	} else if json.Action == "singup" {
		var user User_Singup
		valueMap, ok := json.Value.(map[string]interface{})
		if !ok {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Failed to encode JSON data | Need : {\n\t\"action\":\"singup\",\n\t\"value\":{\"email\":\"Your-Email\",\"id\":\"Your-ID\"},\"pw\":\"Your-Password\"}}")
			return
		}
		data, err := js.Marshal(valueMap)
		if err != nil {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Failed to encode JSON data | Need : {\n\t\"action\":\"singup\",\n\t\"value\":{\"email\":\"Your-Email\",\"id\":\"Your-ID\"},\"pw\":\"Your-Password\"}}")
			return
		}
		err = ParseJSON([]byte(data), &user)
		if err != nil {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Data is invalid | Need : {\n\t\"action\":\"singup\",\n\t\"value\":{\"email\":\"Your-Email\",\"id\":\"Your-ID\"},\"pw\":\"Your-Password\"}}")
			return
		}
		if user.Email == "" || user.ID == "" || user.Pw == "" {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Invalid value (email or id or pw) | Need : {\n\t\"action\":\"singup\",\n\t\"value\":{\"email\":\"Your-Email\",\"id\":\"Your-ID\"},\"pw\":\"Your-Password\"}}")
			return
		}
		email := db.CheckEmail(user.Email)
		if email != nil {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "This email already exists | Need : {\n\t\"action\":\"singup\",\n\t\"value\":{\"email\":\"Your-Email\",\"id\":\"Your-ID\"},\"pw\":\"Your-Password\"}}")
			return
		}
		if db.CheckUser(user.ID) {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "This id already exists | Need : {\n\t\"action\":\"singup\",\n\t\"value\":{\"email\":\"Your-Email\",\"id\":\"Your-ID\"},\"pw\":\"Your-Password\"}}")
			return
		}
		token := memory.NewToken()
		err = smtp.Send([]string{user.Email}, "Email Verification", token)
		if err != nil {
			log.Println(err)
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Invalid email | Need : {\n\t\"action\":\"singup\",\n\t\"value\":{\"email\":\"Your-Email\",\"id\":\"Your-ID\"},\"pw\":\"Your-Password\"}}")
			return
		}
		go db.NewUser(token, user.Email, user.ID, user.Pw)
		go SendMessage(request, http.StatusOK, "Singup successfully | You Need Email Verification\n{\n\t\"action\":\"certi\",\n\t\"value\":{\"email\":\"Your-Email\",\"value\":\"Code received by email\"}}")
		return
	} else if json.Action == "login" {
		if db.CheckAPIToken(json.Value.(string)) == nil {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "There is no such key")
			return
		}
		request_info.Login = true
		request_info.API_TOKEN = json.Value.(string)
		request_info.User = db.CheckAPIToken(json.Value.(string))
		go SendMessage(request, http.StatusOK, "Login successfully")
		return
	} else if json.Action == "login-user-email" {
		var user User_Login
		valueMap, ok := json.Value.(map[string]interface{})
		if !ok {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Failed to encode JSON data | Need : {\n\t\"action\":\"login-user-email\",\n\t\"value\":{\"email\":\"Your-Email\",\"pw\":\"Your-Password\"}}")
			return
		}
		data, err := js.Marshal(valueMap)
		if err != nil {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Failed to encode JSON data | Need : {\n\t\"action\":\"login-user-email\",\n\t\"value\":{\"email\":\"Your-Email\",\"pw\":\"Your-Password\"}}")
			return
		}
		err = ParseJSON([]byte(data), &user)
		if err != nil {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Data is invalid | Need : {\n\t\"action\":\"login-user-email\",\n\t\"value\":{\"email\":\"Your-Email\",\"pw\":\"Your-Password\"}}")
			return
		}
		if user.Email == "" || user.Pw == "" {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "Invalid value (email or pw) | Need : {\n\t\"action\":\"login-user-email\",\n\t\"value\":{\"email\":\"Your-Email\",\"pw\":\"Your-Password\"}}")
			return
		}
		userdata := db.CheckEmail(user.Email)
		if userdata == nil {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "The user does not exist or the password is incorrect | Need : {\n\t\"action\":\"login-user-email\",\n\t\"value\":{\"email\":\"Your-Email\",\"pw\":\"Your-Password\"}}")
			return
		}
		is, _ := memory.Compare(userdata.PW, user.Pw)
		if is {
			go SendErrorMessageWithDisconnet(request, http.StatusBadRequest, "The user does not exist or the password is incorrect | Need : {\n\t\"action\":\"login-user-email\",\n\t\"value\":{\"email\":\"Your-Email\",\"pw\":\"Your-Password\"}}")
			return
		}
		request_info.Login = true
		request_info.API_TOKEN = userdata.API_TOKEN
		request_info.User = userdata
		go SendMessage(request, http.StatusOK, "Login successfully")
		return
	} else if strings.HasPrefix(json.Action, "read:") {
		parts := strings.SplitN(json.Action, ":", 2)
		if parts[1] == "" {
			SendErrorMessage(request, http.StatusBadRequest, "Invalid value (Action:collection) \"read:<collection>\" must be")
			return
		}
		var user *memory.User
		if request_info.User != nil {
			user = request_info.User
		} else {
			user = db.CheckAPIToken(request_info.API_TOKEN)
		}
		msg, _ := database.ByteJSON(&map[string]interface{}{
			"state": http.StatusOK,
			"data":  user.Collection[parts[1]],
		})
		request.Write(msg)
		return
	} else if strings.HasPrefix(json.Action, "write:") {
		parts := strings.SplitN(json.Action, ":", 2)
		if parts[1] == "" {
			SendErrorMessage(request, http.StatusBadRequest, "Invalid value (Action:collection) \"write:<collection>\" must be")
			return
		}
		var user *memory.User
		if request_info.User != nil {
			user = request_info.User
		} else {
			user = db.CheckAPIToken(request_info.API_TOKEN)
		}
		user.Collection[parts[1]] = json.Value
		db.User[user.ID] = user
		go SendMessage(request, http.StatusOK, "Write successfully")
		go db.Backup()
		return
	} else {
		go SendErrorMessage(request, http.StatusBadRequest, "Invalid Action value\nlogin\nlogin-user-email\ncerti\nread:<collection>\nwrite:<collection>")
		return
	}
}
