# Google Photos Data Extractor
A basic go module to extract metadata of media in Google Photo albums. The module uses [Google Photos Library API](https://developers.google.com/photos/library/reference/rest) and needs one of [these Authorization Scopes](https://developers.google.com/photos/library/reference/rest/v1/albums/list#authorization-scopes).

## Usage
- Setting up OAuth2.0 to access Google Photos:
    - Go to [API & services menu in Google developers console](https://console.cloud.google.com/apis)
    - Create a project (or use an existing project), enable Photos Library API for the project.
    - Create OAuth credentials for the project (set project type to Desktop, *note that currently [main.go](main.go) assumes that the credentials are for a Desktop app client. For another client type, change the json tag accordingly in the Credential struct*), and download the client_secret JSON file for the credentials.
    These credentials include the client id and secret from the AuthServer, which this application uses to identify itself to it and get authenticated.
- Install Golang as per [installation instructions](https://go.dev/doc/install).
- Clone the repo and run
```
cd google_photos_data_extractor
go get
go run main.go --client_secret_file=<path to the client_secret.json file you downloaded above>
```
- Open the link printed at the prompt, login with your gmail and wait to be redirected to a webpage showing a list of photo albums owned by your gmail account. Basically, the opening of the suggested link makes the AuthServer:
    - authenticate this application and verify that its requested scopes (readonly here) are permitted by gphotoslibrary, 
    - communicate with gphotoslibrary to get an authorization code,
    - send this code to a redirect URL that exchanges it for Access and Refresh Tokens (latter not used here) to be used by the application for making further requests to the gphotoslibrary.
- Select desired album from the available albums in your Google Photos and click on `Fetch Media Metadata` button to download a CSV file with filename, creation times and URL of the media items in the album. Open any URL in browser to view the media file.
