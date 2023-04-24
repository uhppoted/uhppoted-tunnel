package tailscale

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

type request struct {
	Capabilities capabilities `json:"capabilities"`
	Expiry       uint32       `json:"expirySeconds"`
}

type capabilities struct {
	Devices devices `json:"devices"`
}

type devices struct {
	Create create `json:"create`
}

type create struct {
	Reusable      bool     `json:"reusable"`
	Ephemeral     bool     `json:"ephemeral"`
	Preauthorized bool     `json:"preauthorized"`
	Tags          []string `json:"tags"`
}

type authkey struct {
	ID           string       `json:"id"`
	Key          string       `json:"key"`
	Created      time.Time    `json:"created"`
	Expires      time.Time    `json:"expires"`
	Capabilities capabilities `json:"capabilities"`
}

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
	credentials := struct {
		Tailscale struct {
			OAuth2 struct {
				ClientID     string `json:"client-id"`
				ClientSecret string `json:"client-secret"`
				AuthURL      string `json:"auth-url"`
				Tailnet      string `json:"tailnet"`
				Tag          string `json:"tag"`
				KeyExpiry    uint32 `json:"key-expiry"`
			} `json:"oauth2"`
		} `json:"tailscale"`
	}{}

	if bytes, err := os.ReadFile(file); err != nil {
		return "", err
	} else if err := json.Unmarshal(bytes, &credentials); err != nil {
		return "", err
	}

	var oauth = &clientcredentials.Config{
		ClientID:     credentials.Tailscale.OAuth2.ClientID,
		ClientSecret: credentials.Tailscale.OAuth2.ClientSecret,
		TokenURL:     credentials.Tailscale.OAuth2.AuthURL,
	}

	client := oauth.Client(context.Background())
	url := fmt.Sprintf("https://api.tailscale.com/api/v2/tailnet/%v/keys", credentials.Tailscale.OAuth2.Tailnet)
	tag := fmt.Sprintf("tag:%v", credentials.Tailscale.OAuth2.Tag)
	expiry := credentials.Tailscale.OAuth2.KeyExpiry

	request := request{
		Capabilities: capabilities{
			Devices: devices{
				Create: create{
					Reusable:      true,
					Ephemeral:     true,
					Preauthorized: true,
					Tags:          []string{tag},
				},
			},
		},
		Expiry: expiry,
	}

	var b *bytes.Buffer
	if rq, err := json.Marshal(request); err != nil {
		return "", err
	} else {
		b = bytes.NewBuffer(rq)
	}

	var authKey authkey
	if response, err := client.Post(url, "application/json", b); err != nil {
		return "", err
	} else if body, err := io.ReadAll(response.Body); err != nil {
		return "", err
	} else if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth response:%v %v", response.StatusCode, string(body))
	} else if err := json.Unmarshal(body, &authKey); err != nil {
		return "", err
	} else {
		return authKey.Key, nil
	}
}
