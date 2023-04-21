package tailscale

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/oauth2/clientcredentials"
)

func resolveTailscaleAddr(spec string) (string, uint16, error) {
	if match := regexp.MustCompile("(.*?):([0-9]+)").FindStringSubmatch(spec); len(match) < 3 {
		return "", 0, fmt.Errorf("invalid tailscale address (%v)", spec)
	} else if port, err := strconv.ParseUint(match[2], 10, 16); err != nil {
		return "", 0, err
	} else if port > 65535 {
		return "", 0, fmt.Errorf("invalid tailscale port (%v)", spec)
	} else {
		return match[1], uint16(port), nil
	}
}

func getAuthKey(auth string) (string, error) {
	switch {
	case strings.HasPrefix(auth, "authkey:"):
		return strings.TrimSpace(auth[8:]), nil

	case strings.HasPrefix(auth, "env:"):
		return os.Getenv(strings.TrimSpace(auth[4:])), nil

	case strings.HasPrefix(auth, "oauth2:"):
		return oauth2(auth[7:])
	}

	return "", nil
}

func oauth2(file string) (string, error) {
	if bytes, err := os.ReadFile(file); err != nil {
		return "", err
	} else {
		credentials := struct {
			Tailscale struct {
				OAuth2 struct {
					ClientID     string `json:"client-id"`
					ClientSecret string `json:"client-secret"`
					AuthURL      string `json:"auth-url"`
					Tailnet      string `json:"tailnet"`
				} `json:"oauth2"`
			} `json:"tailscale"`
		}{}

		if err := json.Unmarshal(bytes, &credentials); err != nil {
			return "", err
		} else {
			var oauth = &clientcredentials.Config{
				ClientID:     credentials.Tailscale.OAuth2.ClientID,
				ClientSecret: credentials.Tailscale.OAuth2.ClientSecret,
				TokenURL:     credentials.Tailscale.OAuth2.AuthURL,
			}

			client := oauth.Client(context.Background())
			url := fmt.Sprintf("https://api.tailscale.com/api/v2/tailnet/%v/devices", credentials.Tailscale.OAuth2.Tailnet)

			if response, err := client.Get(url); err != nil {
				return "", err
			} else {
				fmt.Printf(">>>>>>> CREDENTIALS %+v\n", response)
			}
		}
	}

	return "", fmt.Errorf("NOT IMPLEMENTED")
}
