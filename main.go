package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	clientSecretFile = flag.String("client_secret_file", "client_secret.json", "OAuth 2.0 file downloaded from https://console.cloud.google.com/apis/credentials?project=<your_project_name>")
	web              = flag.String("addr", ":9090", "Port to serve on")
)

// Credentials which store client google ids, change if you find the structure of your client secret_file changed.
type Credentials struct {
	Installed struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"installed"`
}

var conf *oauth2.Config
var token *oauth2.Token

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// handler that exchanges the authorization code provided by AuthServer for Access Token
func authHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	token, err = conf.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("http://127.0.0.1%s/", *web), http.StatusSeeOther)
}

// handler that uses a OAuth2 client with Access Token to return list of all albums owned by a gphotos account
func listHandler(w http.ResponseWriter, r *http.Request) {
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://photoslibrary.googleapis.com/v1/albums")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Album struct {
		Id         string `json:"id,omitempty"`
		Title      string `json:"title,omitempty"`
		ProductURL string `json:"productUrl,omitempty"`
		isWritable bool   `json:"isWriteable,omitempty"`
	}
	type AlbumsList struct {
		Albums []Album `json:"albums"`
	}
	data, _ := ioutil.ReadAll(resp.Body)
	albumsList := AlbumsList{}
	if err := json.Unmarshal(data, &albumsList); err != nil {
		log.Println(err)
	}

	templateVars := struct {
		Albums []Album
	}{
		Albums: albumsList.Albums,
	}
	template.Must(template.ParseFiles("data.html")).Execute(w, templateVars)
}

// handler that uses a OAuth2 client with Access Token to return list of all mediaItems in a requested gphotos album
func albumHandler(w http.ResponseWriter, r *http.Request) {
	albumId := strings.Split(r.URL.Path, "/")[2]
	client := conf.Client(oauth2.NoContext, token)
	values := map[string]string{"albumId": albumId}
	json_data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Post("https://photoslibrary.googleapis.com/v1/mediaItems:search", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	flag.Parse()

	// get credentials from the secrets file
	secrets_file, err := ioutil.ReadFile(*clientSecretFile)
	if err != nil {
		log.Printf("Error in reading client secret file: %v\n", err)
		os.Exit(1)
	}
	cred := Credentials{}
	json.Unmarshal(secrets_file, &cred)

	// prepare the oauth2 config with credentials from authorization server
	conf = &oauth2.Config{
		ClientID:     cred.Installed.ClientID,
		ClientSecret: cred.Installed.ClientSecret,
		RedirectURL:  fmt.Sprintf("http://127.0.0.1%s/auth", *web), // url to send the authorization code to
		Scopes: []string{
			"https://www.googleapis.com/auth/photoslibrary.readonly", // list of all scopes here: https://developers.google.com/identity/protocols/oauth2/scopes
		},
		Endpoint: google.Endpoint,
	}
	// NOTE: the implemention here is not secure in that the random state sent to the authorization server when requesting for interactive authentication, is not stored
	// in local or session storage for validating browser requests. It is easy enough to maintain state in storage and compare the set state here and the states in received
	// requests in the handler functions, but this example is deliberately kept minimal.
	state := randToken()
	// opening of the prompted link in the browser makes the AuthServer:
	// - authenticate this application and verify that the scopes requested in OAuth2 config are permitted,
	// - communicate with the resource owner (photoslibrary) to get an authorization code,
	// - send the authorization code to the redirect URL specified in the OAuth2 config.
	fmt.Printf("Open the following link in your browser, login and grant permission to the app to access your Google Photos to proceed.\n\n%s", conf.AuthCodeURL(state))

	http.HandleFunc("/", listHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/album/", albumHandler)
	err = http.ListenAndServe(*web, nil)
	if err != nil {
		panic(err)
	}
}
