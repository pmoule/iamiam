# IamIam (yumyum)
A very simple oauth provider for testing purposes.

## How to use
Configure hostname, port and the array of valid userinfos in `iam_config.json`.
```
{
  "hostname": "localhost",
  "port": 8081,
  "validUserInfos": [
    {
      "email": "foo@bar.com",
      "firstName": "foo",
      "lastName": "bar"
    }
  ]
}
```

This is everything required. Now run iamiam.

For initiating authentication from your client application with iamiam, be inspired by the following code snippet.

```
// just example values for hostname and port where iamiam service is running.
hostname = "localhost"
port = 8081

// create config for iamiam
oauthConfig := &oauth2.Config{
            //set your callback URL
			RedirectURL:  "http://localhost:8080/login/callback",   
            //currently not required
			ClientID:     "any",       
            //currently not required                                
			ClientSecret: "any",
            // choose from iamiam.Email, iamiam.SimpleProfile
			Scopes:       []string{iamiam.Email},
			Endpoint:     iamiam.CreateEndpoint(hostname, port),
		}

// create some state for confirmation
b := make([]byte, 16)
rand.Read(b)
state := base64.URLEncoding.EncodeToString(b)


// start authentication process
u := oauthConfig.AuthCodeURL(state)
http.Redirect(w, r, u, http.StatusTemporaryRedirect)
```
In the shown form just enter an email you like to use for authentication. (Email must be provided in `iam_config.json`). After this do the following steps in your login callback handler.
```
func loginCallback(w http.ResponseWriter, r *http.Request) {
	// do stuff like checking oauth state 
	...
	// read code from request body
	code := r.FormValue("code")
	token, _ := oauthConfig.Exchange(context.Background(), code)
	//use iamiam's hostname and port for requesting user info
	userDataURI = iamiam.CreateInfoURL(hostname, port, token.AccessToken)
	response, _ := http.Get(userDataURI)
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	var userInfo iamiam.UserInfo
	json.Unmarshal(data, &userInfo)
}
```
## Documentation
See package documentation:

[![GoDoc](https://godoc.org/github.com/pmoule/iamiam?status.svg)](https://godoc.org/github.com/pmoule/iamiam)

## License
`iamiam` is released under Apache License 2.0. See [LICENSE](LICENSE).