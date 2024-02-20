package main

import (
	"strings"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"net/http"
)

type AuthData struct {
	AccessToken    string `json:"access_token"`
	Bans           []any  `json:"bans,omitempty"`
	DisplayName    string `json:"display_name,omitempty"`
	ExpiresIn      int    `json:"expires_in,omitempty"`
	IsComply       bool   `json:"is_comply,omitempty"`
	Jflgs          int    `json:"jflgs,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	NamespaceRoles []struct {
		RoleID    string `json:"roleId,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	} `json:"namespace_roles,omitempty"`
	Permissions      []any    `json:"permissions,omitempty"`
	PlatformID       string   `json:"platform_id,omitempty"`
	PlatformUserID   string   `json:"platform_user_id,omitempty"`
	RefreshExpiresIn int      `json:"refresh_expires_in,omitempty"`
	RefreshToken     string   `json:"refresh_token"`
	Roles            []string `json:"roles,omitempty"`
	Scope            string   `json:"scope,omitempty"`
	TokenType        string   `json:"token_type,omitempty"`
	UserID           string   `json:"user_id,omitempty"`
	Xuid             string   `json:"xuid,omitempty"`
}

type Instance struct {
	username string
	password string

	token string
	refreshToken string

	elo float64

	inMatchmaking bool
	inGame bool
}

func NewInstance(Username, Password string) *Instance {
	i := new(Instance)
	i.username = Username
	i.password = Password

	return i
}

func (usr *Instance) Authorize() error {
	params := url.Values{}
	params.Add("grant_type", `password`)
	params.Add("username", usr.username)
	params.Add("password", usr.password)

	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://api.voxies.io/iam/oauth/token", body)
	if err != nil {
		return err
	}
	req.Host = "api.voxies.io"
	req.Header.Set("User-Agent", "curl/7.80.0-DEV")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Basic YzcyMTQzNjhlMTIxNDAxMTljYTNmN2ExOTdmZjViMzE6ZCk9ZVdXVVUhbF1KamY5Q3VVb3VWZEQzZHVQTD0mZG4=")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	var data AuthData

	if err := json.Unmarshal(respBody, &data); err != nil {
		return err
	}

	usr.token = data.AccessToken
	usr.refreshToken = data.RefreshToken

	return nil
}
