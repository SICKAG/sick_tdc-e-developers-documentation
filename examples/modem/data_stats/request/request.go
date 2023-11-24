/* Package created 04.10.2023. for SICK Mobilisis d.o.o. */
/* Handles HTTP Requests with OAuth2.0 */
package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	oauth2 "golang.org/x/oauth2"
)

type OAuthConf struct {
	ClientId              string `json:"clientId"`
	ClientSecret          string `json:"clientSecret"`
	AuthorizationEndpoint string `json:"authorizationEndpoint"`
	TokenEndpoint         string `json:"tokenEndpoint"`
	RedirectURL           string `json:"redirectURL"`
	Username              string `json:"username"`
	Password              string `json:"password"`
}

type Params struct {
	OAuthconf []OAuthConf `json:"oauthConf"`
}

/* Fetches needed data from .json file */
func setConfValues(jsonFile *os.File) (oauth2.Config, url.Values) {
	byteValue, _ := io.ReadAll(jsonFile)
	var params Params
	json.Unmarshal(byteValue, &params)
	conf := params.OAuthconf[0]

	cfg := oauth2.Config{
		ClientID:     conf.ClientId,
		ClientSecret: conf.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: conf.TokenEndpoint,
		},
	}

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", conf.Username)
	data.Set("password", conf.Password)

	return cfg, data
}

/* Opens JSON file from folder and forwards it to function createAccessTokenFromFile */
func createOAuth2Config() (oauth2.Config, url.Values) {
	jsonFile, err := os.Open("params.json")
	if err != nil {
		fmt.Println("Error opening config file: ", err)
	}
	defer jsonFile.Close()
	return setConfValues(jsonFile)
}

func Authorize() string {
	cfg, data := createOAuth2Config()
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", cfg.Endpoint.TokenURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cfg.ClientID, cfg.ClientSecret)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", resp.Status)
		return ""
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(response, &responseMap); err != nil {
		fmt.Println("Error decoding JSON response: ", err)
		return ""
	}

	accessToken, ok := responseMap["access_token"].(string)
	if !ok {
		fmt.Println("Access token not found in response")
		return ""
	}

	return accessToken
}

/* Function for configuring request and fetching data from url with provided accessToken */
func MakeROPCRequest(urlConn string, accessToken string) []byte {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", urlConn, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return responseBody
}
