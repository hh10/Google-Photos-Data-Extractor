# Google Photos Data Extractor
A very primitive go module to extract metadata of media in Google Photo albums. The module uses [Google Photos Library API](https://developers.google.com/photos/library/reference/rest) and needs one of [these Authorization Scopes](https://developers.google.com/photos/library/reference/rest/v1/albums/list#authorization-scopes).

## Usage
- Setting up OAuth 2.0 to access Google Photos:
    Go to [API & services menu in Google developers console](https://console.cloud.google.com/apis), create a project, enable Photos Library API for the project, create OAuth credentials for the project (set project type to Desktop), download the client_secret JSON file for the credentials.
- Install Golang as per [installation instructions](https://go.dev/doc/install).
- Clone the repo and run
```
cd google_photos_data_extractor
go get
go run main.go --client_secret_file=<path to the client_secret.json file you downloaded above>
```
- Follow the instructions at the prompt, i.e., open the suggested link, login with your gmail.
- Select desired album from the available albums in your Google Photos and click on `Fetch Media Metadata` button to download a CSV file with filename, creation times and URL of the media items in the album. Open any URL in browser to view the media file.
