package catelyn

import "os"
import "io"
import "fmt"
import "log"
import "net/url"
import "net/http"
import "encoding/json"
import "encoding/base64"
import "github.com/dadleyy/catelyn/catelyn/constants"

// ConfluenceClient is an interface to the confluence http rest api.
type ConfluenceClient interface {
	SearchSpaces(string) ([]ConfluenceSpace, *ConfluencePaging, error)
}

// NewConfluenceClient returns a client implementing the client interface.
func NewConfluenceClient(user *url.Userinfo, hostname string) (ConfluenceClient, error) {
	logger := log.New(os.Stdout, constants.ConfluenceClientLoggerPrefix, log.LstdFlags)
	password, _ := user.Password()

	creds := confluenceCredentials{
		username: user.Username(),
		password: password,
	}

	api, e := url.Parse(hostname)

	if e != nil {
		return nil, e
	}

	if api.Scheme == "" {
		api.Scheme = "https"
	}

	client := &confluenceClient{
		credentials: creds,
		logger:      logger,
		apiHome:     api,
	}

	return client, nil
}

type confluenceCredentials struct {
	username string
	password string
}

func (c confluenceCredentials) String() string {
	joined := []byte(fmt.Sprintf("%s:%s", c.username, c.password))
	return base64.StdEncoding.EncodeToString(joined)
}

type confluenceClient struct {
	logger      *log.Logger
	credentials confluenceCredentials
	apiHome     *url.URL
}

func (c *confluenceClient) SearchSpaces(query string) ([]ConfluenceSpace, *ConfluencePaging, error) {
	r, e := c.send("GET", fmt.Sprintf("%s/%s", c.apiHome, constants.SpacesAPIEndpoint), nil)

	if e != nil {
		return nil, nil, e
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	response := struct {
		ConfluencePaging
		Results []ConfluenceSpace `json:"results"`
	}{}

	if r.StatusCode != 200 {
		return nil, nil, fmt.Errorf("invalid response from confluence: %d", r.StatusCode)
	}

	if e := decoder.Decode(&response); e != nil {
		return nil, nil, e
	}

	return response.Results, &response.ConfluencePaging, nil
}

func (c *confluenceClient) send(method string, url string, body io.Reader) (*http.Response, error) {
	client := http.Client{}
	request, e := http.NewRequest("GET", url, nil)

	if e != nil {
		return nil, e
	}

	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.credentials))
	return client.Do(request)
}
