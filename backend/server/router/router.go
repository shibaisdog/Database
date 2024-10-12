package router

import (
	"cloud-logic/database"
	"cloud-logic/database/memory"
	"net"
	"net/http"
	"strings"
)

type Client struct {
	Login     bool
	API_TOKEN string
}

type JSON struct {
	Action string      `json:"action"`
	Value  interface{} `json:"value"`
}

func Router(db *memory.Storage, request_info *Client, request net.Conn, butter []byte) {
	var json JSON
	err := ParseJSON(butter, &json)
	if err != nil {
		msg, _ := database.ByteJSON(&map[string]interface{}{
			"state": http.StatusBadRequest,
			"error": "Data is invalid | Need : {\n\t\"action\":\"\"\n\t\"value\":\"\"\n}",
		})
		request.Write(msg)
		request.Close()
		return
	}
	if (json.Action != "login" && !request_info.Login) || json.Value == "" {
		msg, _ := database.ByteJSON(&map[string]interface{}{
			"state": http.StatusBadRequest,
			"error": "Login required | Need : {\n\t\"action\":\"login\"\n\t\"value\":\"YOUR-API-TOKEN\"\n}",
		})
		request.Write(msg)
		return
	}
	if json.Action == "login" {
		if db.CheckAPIToken(json.Value.(string)) != nil {
			request_info.Login = true
			request_info.API_TOKEN = json.Value.(string)
			msg, _ := database.ByteJSON(&map[string]interface{}{
				"state": http.StatusOK,
			})
			request.Write(msg)
			return
		} else {
			msg, _ := database.ByteJSON(&map[string]interface{}{
				"state": http.StatusBadRequest,
				"error": "There is no such key",
			})
			request.Write(msg)
			return
		}
	} else if strings.HasPrefix(json.Action, "read:") {
		parts := strings.SplitN(json.Action, ":", 2)
		if parts[1] == "" {
			msg, _ := database.ByteJSON(&map[string]interface{}{
				"state": http.StatusBadRequest,
				"error": "Invalid value \"read:<collection>\" must be",
			})
			request.Write(msg)
			return
		}
		user := db.CheckAPIToken(request_info.API_TOKEN)
		msg, _ := database.ByteJSON(&map[string]interface{}{
			"state": http.StatusOK,
			"data":  user.Collection[parts[1]],
		})
		request.Write(msg)
		return
	} else if strings.HasPrefix(json.Action, "write:") {
		parts := strings.SplitN(json.Action, ":", 2)
		if parts[1] == "" {
			msg, _ := database.ByteJSON(&map[string]interface{}{
				"state": http.StatusBadRequest,
				"error": "Invalid value \"write:<collection>\" must be",
			})
			request.Write(msg)
			return
		}
		user := db.CheckAPIToken(request_info.API_TOKEN)
		user.Collection[parts[1]] = json.Value
		db.User[user.ID] = user
		msg, _ := database.ByteJSON(&map[string]interface{}{
			"state": http.StatusOK,
		})
		request.Write(msg)
		db.Backup()
		return
	} else {
		msg, _ := database.ByteJSON(&map[string]interface{}{
			"state": http.StatusBadRequest,
			"error": "Invalid Action value\nlogin\nread:<collection>\nwrite:<collection>",
		})
		request.Write(msg)
		return
	}
}
