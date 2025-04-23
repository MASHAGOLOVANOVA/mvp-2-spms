// yandexapi/api.go
package yandexapi

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
	"net/http"
)

type YandexAPI struct {
	Client  *http.Client
	Context context.Context
	Config  *oauth2.Config
	Token   *oauth2.Token
}

func InitYandexAPI() (*YandexAPI, error) {
	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     "33663435c31443af933742df94b2516f",
		ClientSecret: "2d9259feb4f8400bbc6511b1c0d8d9b3",
		RedirectURL:  "http://localhost:8080/api/v1/auth/integration/access/yandexdisk",
		Scopes:       []string{"login:email", "cloud_api:disk.read", "cloud_api:disk.info", "cloud_api:disk.write", "cloud_api:disk.app_folder"},
		Endpoint:     yandex.Endpoint,
	}

	return &YandexAPI{
		Context: ctx,
		Config:  config,
	}, nil
}

func (y *YandexAPI) GetAuthLink(statestr string) string {
	return y.Config.AuthCodeURL(statestr, oauth2.AccessTypeOffline)
}

func (y *YandexAPI) ExchangeCode(code string) (*oauth2.Token, error) {
	token, err := y.Config.Exchange(y.Context, code)
	if err != nil {
		return nil, err
	}
	y.Token = token
	y.Client = y.Config.Client(y.Context, token)
	return token, nil
}

// SetupClient настраивает HTTP клиент с переданным токеном
func (y *YandexAPI) SetupClient(token *oauth2.Token) error {
	if token == nil {
		return errors.New("token is nil")
	}

	// Если AccessToken пустой, но есть RefreshToken - обновляем токен
	if token.AccessToken == "" && token.RefreshToken != "" {
		newToken, err := y.Config.TokenSource(y.Context, token).Token()
		if err != nil {
			return err
		}
		token = newToken
	}

	y.Token = token
	y.Client = y.Config.Client(y.Context, token)
	return nil
}
