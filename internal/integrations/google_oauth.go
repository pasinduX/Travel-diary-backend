package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

func GoogleOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func FetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (GoogleUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return GoogleUserInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return GoogleUserInfo{}, fmt.Errorf("google userinfo request failed: %s", resp.Status)
	}

	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return GoogleUserInfo{}, err
	}
	return info, nil
}
