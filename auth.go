package main

import (
	"fmt"
	"io"
	"net/http"
)

const AUTH_URL = "https://id.twitch.tv/oauth2/authorize?response_type=token&client_id=%s&redirect_uri=https://twitchapps.com/tmi/&scope=chat:read+chat:edit+channel:moderate+whispers:read+whispers:edit+channel_editor"

func GetOauthToken(client string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(AUTH_URL, client))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(bytes))
	return "", nil
}
