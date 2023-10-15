package auth

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/arpanrec/secureserver/internal/common"
)

var (
	users map[string]interface{}
	mo    = &sync.Once{}
)

func getUsers() map[string]interface{} {
	mo.Do(func() {
		users = common.GetConfig()["users"].(map[string]interface{})
	})
	return users
}

func GetUserDetails(authHeader string) (string, bool) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", false
	}
	usersDB := getUsers()
	credentialsBase64 := authHeader[6:]
	credentialsBytes, err := base64.StdEncoding.DecodeString(credentialsBase64)
	if err != nil {
		return "", false
	}
	credentials := string(credentialsBytes)
	username := strings.Split(credentials, ":")[0]
	password := strings.Split(credentials, ":")[1]
	_, ok := usersDB[username]
	if ok {
		githubLogin, githubOKStatus := getGithubUserDetails(password)
		if !githubOKStatus {
			return "", false
		}
		if githubLogin == username {
			return username, true
		}
	}
	return "", false
}

func getGithubUserDetails(gPAT string) (string, bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Fatalln("Error creating request: ", err)
		return "", false
	}
	req.Header.Add("Authorization", "token "+gPAT)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error sending request: ", err)
		return "", false
	}
	githubResBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response: ", err)
		return "", false
	}
	var githubResp map[string]interface{}
	ok := json.Unmarshal(githubResBody, &githubResp)
	if ok != nil {
		log.Fatalln("Error unmarshalling response: ", ok)
		return "", false
	}
	log.Println("Github User Details: ", githubResp["login"].(string))
	return githubResp["login"].(string), true
}
