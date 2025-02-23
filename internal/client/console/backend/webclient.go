package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang-jwt/jwt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type WebClient struct {
	jwt            []byte
	isValid        bool
	reAuthCallback func()
	baseUrl        string
	jwtFilePath    string
	hasFilePath    bool
	storeInFile    bool
	exp            float64
}

type WebClientOption func(*WebClient) error

// region Constructor

// NewWebClient returns a new WebClient struct, optionally initialized with one or more WebClientOption parameters.
// Some WebClientOption may return an error, the first error caught is returned through the error return var.
func NewWebClient(options ...WebClientOption) (*WebClient, error) {
	client := &WebClient{
		baseUrl: ServerUrl,
	}
	for _, option := range options {
		err := option(client)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return client, nil
}

// WithJwt specifies a pre-retrieved JWT as an array of bytes. If this is specified and WithStoreInFile is specified,
// then the file will not be read upon startup.
func WithJwt(jwt []byte) WebClientOption {
	return func(client *WebClient) error {
		client.jwt = jwt
		_, _, client.exp = client.Identify()
		return nil
	}
}

// WithStoreInFile specifies that the `path` should be used to read/write JWT data from/to when `enabled` is set to true.
func WithStoreInFile(enabled bool, path string) WebClientOption {
	log.Printf("Location for storing JWT files = %s. Enabled = %t", path, enabled)
	return func(client *WebClient) error {
		client.jwtFilePath = path
		client.hasFilePath = true
		client.storeInFile = enabled
		if enabled && client.jwt == nil {
			return client.loadFromFile(path)
		}
		return nil
	}
}

// WithReAuthCallback specifies a func that should be called when the webclient determines that re-authentication
// is needed.
func WithReAuthCallback(callback func()) WebClientOption {
	return func(client *WebClient) error {
		client.reAuthCallback = callback
		return nil
	}
}

// WithBaseUrl specifies the base url to use for all http calls. The `url` argument in the `Call...` methods with be
// appended to this `baseUrl` before the call is made.
func WithBaseUrl(baseUrl string) WebClientOption {
	return func(client *WebClient) error {
		client.baseUrl = baseUrl
		return nil
	}
}

// endregion

// region Local storage

func (wc *WebClient) loadFromFile(path string) error {
	log.Printf("Reading jwt from file from location '%s'\n", path)
	file, err := os.Open(path)
	if err != nil {
		log.Printf("ERROR: Could not open jwt file: %v", err)
		return err
	}
	wc.jwt, err = io.ReadAll(file)
	if err != nil {
		log.Printf("ERROR: Could not read jwt from file: %v", err)
		return err
	}

	var name, email string
	name, email, wc.exp = wc.Identify()
	if strings.TrimSpace(name) != "" && strings.TrimSpace(email) != "" {
		wc.isValid = true
		log.Printf("Logged in as %s", email)
	}

	return nil
}

func (wc *WebClient) writeJwtToFile() {
	log.Printf("Writing jwt to file for user to location '%s'\n", JwtFileName())
	err := os.WriteFile(JwtFileName(), wc.jwt, 0666)
	if err != nil {
		log.Printf("ERROR: Could not write jwt to file: %v", err)
		tea.Println("The login credentials could not be written to file.")
	}
}

func (wc *WebClient) removeJwtFile() {
	log.Printf("Removing jwt file for user to location '%s'\n", JwtFileName())
	err := os.Remove(JwtFileName())
	if err != nil {
		log.Printf("ERROR: Could not remove jwt file: %v", err)
		log.Println("Disabling storage of the JWT file for now.")
		wc.isValid = false
		wc.storeInFile = false
	}
}

// endregion

// region Calls

func (wc *WebClient) Url(parts ...string) string {
	out, _ := url.JoinPath(wc.baseUrl, parts...)
	return out
}

// Call makes a request to the specified url, using the specified method (use e.g. http.MethodGet or http.MethodPost)
// and then tries to decode the returned JSON into the 'output' var.
func (wc *WebClient) Call(method string, url string, output any) error {
	return wc.CallWithBody(method, url, nil, output)
}

func (wc *WebClient) CallLogin(url string, body any) error {

	bodyJson, _ := json.Marshal(body)
	buf := bytes.NewBuffer(bodyJson)
	wc.isValid = false

	response, err := http.Post(url, "application/json", buf)
	if err != nil {
		log.Printf("There was an error making a request to the api: %v\n", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("The api responded with an error: %d - %s\n", response.StatusCode, response.Status)
		if response.StatusCode == http.StatusUnauthorized {
			return errors.New("invalid credentials")
		} else {
			return errors.New(fmt.Sprintf("Login failed: %d - %s", response.StatusCode, response.Status))
		}
	}

	defer response.Body.Close()
	wc.jwt, err = io.ReadAll(response.Body)
	_, _, wc.exp = wc.Identify()

	if err != nil {
		log.Printf("Could not read the JWT returned from the login api: %v\n", err)
		return err
	}

	if len(wc.jwt) == 0 {
		log.Println("The JWT returned from the login api is empty")
		return errors.New("invalid JWT returned from the login api")
	}

	wc.isValid = true
	if wc.hasFilePath && wc.storeInFile {
		wc.writeJwtToFile()
	}

	return nil

}

func (wc *WebClient) CallWithBody(method string, url string, body any, output any) error {

	if wc.IsExpired() {
		wc.reAuth()
	}

	var req *http.Request
	var buf *bytes.Buffer
	if body != nil {
		bodyJson, _ := json.Marshal(body)
		buf = bytes.NewBuffer(bodyJson)
		req, _ = http.NewRequest(method, url, buf)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Add("Authorization", "Bearer "+string(wc.jwt))
	client := &http.Client{}
	// Make the actual request
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v\n", err)
		return fmt.Errorf("making the request to the server failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			// indicate that we need to (re)authenticate
			wc.reAuth()
			return errors.New("invalid credentials - please authenticate")
		}
		log.Printf("The api responded with an error: %d - %s\n", response.StatusCode, response.Status)
		return errors.New(fmt.Sprintf("server responded with error: %d %s", response.StatusCode, response.Status))
	}

	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	err = dec.Decode(output)

	if err != nil {
		log.Printf("Decoding the response failed: %v\n", err)
		return fmt.Errorf("decoding the response failed %w", err)
	}

	return nil
}

// endregion

// region Validity

// reAuth removes the current JWT information (also from disk) and calls the 'reAuthCallback' to indicate that the
// app should ask the user for their credentials.
func (wc *WebClient) reAuth() {
	log.Printf("The login credentials were invalid and re-authentication is required.")
	wc.jwt = nil
	wc.isValid = false
	wc.exp = 0
	if wc.storeInFile {
		wc.removeJwtFile()
	}
	wc.reAuthCallback()
}

// IsValid returns whether the current JWT is considered valid.
func (wc *WebClient) IsValid() bool {
	return wc.isValid
}

// Identify reads the JWT token to extract the name, email and expiry date.
func (wc *WebClient) Identify() (name string, email string, exp float64) {
	token, _, err := new(jwt.Parser).ParseUnverified(string(wc.jwt), jwt.MapClaims{})
	if err != nil {
		log.Printf("ERROR: Could not parse jwt token: %v", err)
		return "", "", 0
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email, ok = claims["email"].(string)
		name, ok = claims["name"].(string)
		exp, ok = claims["exp"].(float64)
		return name, email, exp
	}
	log.Printf("ERROR: Could not get identity from jwt token")
	return "", "", 0
}

// IsExpired returns whether the expiration date of the known JTW is in the past.
func (wc *WebClient) IsExpired() bool {
	return wc.exp == 0 ||
		time.Unix(int64(wc.exp), 0).Before(time.Now())
}

// endregion
