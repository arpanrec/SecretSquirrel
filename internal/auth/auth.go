package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/arpanrec/secureserver/internal/appconfig"
)

var (
	usersDb map[string]appconfig.UserConfig
	mo      = &sync.Once{}
)

func getUsers() map[string]appconfig.UserConfig {
	mo.Do(func() {
		usersDb = appconfig.GetConfig().UserDb
	})
	return usersDb
}

func GetUserDetails(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", errors.New("invalid Authorization Header")
	}
	usersDB := getUsers()
	credentialsBase64 := authHeader[6:]
	credentialBytes, err := base64.StdEncoding.DecodeString(credentialsBase64)
	if err != nil {
		return "", err
	}
	credentials := string(credentialBytes)
	username := strings.Split(credentials, ":")[0]
	password := strings.Split(credentials, ":")[1]
	_, ok := usersDB[username]
	if ok {
		githubLogin, errGitHub := getGithubUserDetails(password)
		if errGitHub != nil {
			return "", errGitHub
		}
		if githubLogin == username {
			return username, nil
		}
		return "", errors.New("username and Github Username do not match")
	}
	return "", errors.New("user not found")
}

func getGithubUserDetails(gPAT string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Println("Error creating request: ", err)
		return "", err
	}
	req.Header.Add("Authorization", "token "+gPAT)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request: ", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		log.Println("Error response status code: ", resp.StatusCode)
		return "", errors.New(fmt.Sprintf("Invalid Github Token, Error response status code from GitHub: %v", resp.StatusCode))
	}
	githubResBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response: ", err)
		return "", err
	}
	var githubResp map[string]interface{}
	ok := json.Unmarshal(githubResBody, &githubResp)
	if ok != nil {
		log.Println("Error unmarshalling response: ", ok)
		return "", ok
	}
	log.Println("Github User Details: ", githubResp["login"].(string))
	return githubResp["login"].(string), nil
}
