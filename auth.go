package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	Expiration  int    `json:"expires_in"`
	Type        string `json:"token_type"`
}

func GetAccessToken(cid, csec string) (tok AccessToken, err error) {
	sbody := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", cid, csec)
	resp, err := http.Post("https://id.twitch.tv/oauth2/token", "application/x-www-form-urlencoded", strings.NewReader(sbody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&tok)
	return
}
