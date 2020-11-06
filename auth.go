package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pmoule/iamiam/iamiam"
)

const (
	loginFormFile string = "login.html"
)

// AuthCodeEntry an entry associated with a code.
type AuthCodeEntry struct {
	RedirectURI string
	State       string
	Scope       string
	email       string
}

// SetEmail email for entry.
func (e *AuthCodeEntry) SetEmail(email string) {
	e.email = email
}

// Email returns entry's email.
func (e *AuthCodeEntry) Email() string {
	return e.email
}

// AuthTokenEntry  an entry associated with a token.
type AuthTokenEntry struct {
	Scope string
	Email string
}

var knownUsers []*iamiam.UserInfo = []*iamiam.UserInfo{}
var authCodeRegistry map[string]*AuthCodeEntry = map[string]*AuthCodeEntry{}
var authTokenRegistry map[string]*AuthTokenEntry = map[string]*AuthTokenEntry{}

func auth(w http.ResponseWriter, r *http.Request) {
	redirectURIValue, ok := r.URL.Query()["redirect_uri"]

	if !ok || len(redirectURIValue[0]) < 1 {
		log.Println("URL param 'redirect_uri' is missing")
		return
	}

	redirectURI := redirectURIValue[0]
	stateValue, ok := r.URL.Query()["state"]

	if !ok || len(stateValue[0]) < 1 {
		log.Println("URL param 'state' is missing")
		return
	}

	state := stateValue[0]
	scopeValue, ok := r.URL.Query()["scope"]

	if !ok || len(scopeValue[0]) < 1 {
		log.Println("URL param 'scope' is missing")
		return
	}

	scope := scopeValue[0]
	entry := &AuthCodeEntry{RedirectURI: redirectURI, State: state, Scope: scope}
	code := createRandomValue()
	authCodeRegistry[code] = entry
	showLoginForm(code, w)

	return
}

func useEmail(w http.ResponseWriter, r *http.Request) {
	codeValue, ok := r.URL.Query()["code"]

	if !ok || len(codeValue[0]) < 1 {
		log.Println("URL param 'code' is missing")
		return
	}

	code := codeValue[0]
	entry, ok := authCodeRegistry[code]

	if !ok {
		log.Println("'code' not found")
		return
	}

	email := r.FormValue("email")
	trimmedEmail := strings.TrimSpace(email)

	if !isKnownUser(trimmedEmail) {
		showLoginForm(code, w)

		return
	}

	entry.SetEmail(trimmedEmail)
	u, _ := url.Parse(entry.RedirectURI)
	queryString := u.Query()
	queryString.Set("state", entry.State)
	queryString.Set("code", code)
	u.RawQuery = queryString.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func token(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")
	codeEntry, ok := authCodeRegistry[code]

	if !ok {
		log.Println("'code' not found")
		return
	}

	tokenEntry := &AuthTokenEntry{Scope: codeEntry.Scope, Email: codeEntry.Email()}
	accessToken := createRandomValue()
	refreshToken := createRandomValue()
	authTokenRegistry[accessToken] = tokenEntry
	createHeaders(w)
	token := iamiam.Token{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "", ExpiresIn: 0}
	json.NewEncoder(w).Encode(token)
}

func info(w http.ResponseWriter, r *http.Request) {
	tokenValue, ok := r.URL.Query()["access_token"]

	if !ok || len(tokenValue[0]) < 1 {
		log.Println("URL param 'access_token' is missing")
		return
	}

	token := tokenValue[0]
	tokenEntry, ok := authTokenRegistry[token]
	userInfo := findUserInfo(tokenEntry.Email)

	switch tokenEntry.Scope {
	case iamiam.EmailProfile:
		userInfo = userInfo.CreateEmailProfile()
	case iamiam.SimpleProfile:
		userInfo = userInfo.CreateSimpleProfile()
	}

	createHeaders(w)
	json.NewEncoder(w).Encode(userInfo)
}

func showLoginForm(code string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	loginForm := createLoginForm(code)
	fmt.Fprint(w, loginForm)
}

func createLoginForm(code string) string {
	html, err := ioutil.ReadFile(loginFormFile)
	form := ""

	if err != nil {
		form = fmt.Sprintf(`
		<form method="POST" action="/use?code=%s">
		<input type="text" name="email" placeholder="emails">
		<input type="submit" value="Use">
		</form>
		`, code)
	} else {
		form = string(html)
	}

	return form
}

func createHeaders(w http.ResponseWriter) {
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
	//CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func createRandomValue() string {
	b := make([]byte, 16)
	rand.Read(b)
	rValue := base64.URLEncoding.EncodeToString(b)

	return rValue
}

func findUserInfo(email string) *iamiam.UserInfo {
	for _, u := range knownUsers {
		if u.Email == email {
			return u
		}
	}

	return nil
}

func isKnownUser(email string) bool {
	for _, u := range knownUsers {
		if u.Email == email {
			return true
		}
	}

	return false
}
