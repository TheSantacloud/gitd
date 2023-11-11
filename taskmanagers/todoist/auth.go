package todoist

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"mgtd/adapters"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

const (
	serviceName     = "mgtd"
	accountName     = "AccessToken"
	redirectUriPort = 13371
)

var (
	oauthConfig *oauth2.Config
	state       = "mgtd"
	tokenMu     sync.Mutex
	authToken   *oauth2.Token
	done        = make(chan struct{})
)

func GenerateAccessToken(todoistConfig adapters.TodoistConfig) (string, error) {
	switch todoistConfig.AuthType {
	case adapters.AuthTypeToken:
		return *todoistConfig.AuthToken, nil
	case adapters.AuthTypeOAuth2:
		return PerformOAuthFlow(todoistConfig), nil
	default:
		return "", fmt.Errorf("unknown auth type: %s", todoistConfig.AuthType)
	}
}

func PerformOAuthFlow(todoistConfig adapters.TodoistConfig) string {
	// TODO: sync api is not supported by oauth2, reuse when i have a server
	clientID := *todoistConfig.ClientID
	clientSecret := *todoistConfig.ClientSecret
	scopes := *todoistConfig.Scopes
	redirectURL := fmt.Sprintf("http://localhost:%d/callback", redirectUriPort)

	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://todoist.com/oauth/authorize",
			TokenURL: "https://todoist.com/oauth/access_token",
		},
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}

	token, err := loadTokenFromKeychain()
	if token == nil || err != nil {
		StartAuthServer(done)
		<-done

		if err != nil {
			fmt.Println("Error getting auth token:", err)
			os.Exit(1)
		}

		if err := saveTokenToKeychain(authToken); err != nil {
			fmt.Println("Error saving token to keychain:", err)
			os.Exit(1)
		}
	}

	return token.AccessToken
}

func StartAuthServer(done chan struct{}) {
	http.HandleFunc("/callback", authenticate)
	open(oauthConfig.AuthCodeURL(state))
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", redirectUriPort), nil)
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting HTTP server:", err)
		}
		close(done)
	}()
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func saveTokenToKeychain(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	cmd := exec.Command("security", "add-generic-password", "-U", "-s", serviceName, "-a", accountName, "-w", string(data))
	fmt.Println(cmd.String())
	return cmd.Run()
}

func loadTokenFromKeychain() (*oauth2.Token, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", serviceName, "-a", accountName, "-w")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tokenKeychain *oauth2.Token
	if err := json.Unmarshal(output, &tokenKeychain); err != nil {
		return nil, err
	}
	return tokenKeychain, nil
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not provided", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Token exchange failed: %v", err), http.StatusInternalServerError)
		return
	}

	tokenMu.Lock()
	authToken = token
	tokenMu.Unlock()

	close(done)

	fmt.Fprint(w, "Authentication successful! You can close this window.")
}
