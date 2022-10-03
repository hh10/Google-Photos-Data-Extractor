# Google Photos Data Extractor
A very primitive go module to extract metadata of media in Google Photo albums. The module uses capabilities from [Google Photos Library API](https://developers.google.com/photos/library/reference/rest) and needs one of [these Authorization Scopes](https://developers.google.com/photos/library/reference/rest/v1/albums/list#authorization-scopes).

## Usage
- Setting up OAuth 2.0 to access 
    Go to [API & services menu in Google developers console](https://console.cloud.google.com/apis), create a project, create credentials for the project, enable Google Photos API for the project, download the client_secret JSON file for the credential.
- Install Golang as per [installation instructions](https://go.dev/doc/install).
- Clone the repo and run
```
cd google_photos_data_extractor
go get
go run main.go --client_secret_file=<path to the client_secret.json file you downloaded above>
```
- Follow the instructions at the prompt, i.e., open the suggested link, login with your gmail.
- Select from the desired albums from the available albums in your Google Photos and click on `fetch Media Metadata` to doenload a CSV file with filename, creation times and URL. Open any URL in browser to view the file.