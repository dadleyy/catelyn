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
	SearchSpaces(*ConfluenceSpaceSearchInput) ([]ConfluenceSpace, *ConfluencePaging, error)
}

// NewConfluenceClient returns a client implementing the client interface.
func NewConfluenceClient(user *url.Userinfo, hostname string) (ConfluenceClient, error) {
	logger := log.New(os.Stdout, constants.ConfluenceClientLoggerPrefix, log.LstdFlags)
	api, e := url.Parse(hostname)

	if e != nil {
		return nil, e
	}

	if api.Scheme == "" {
		api.Scheme = "https"
	}

	client := &confluenceClient{
		credentials: user,
		logger:      logger,
		apiHome:     api,
	}

	return client, nil
}

type confluenceClient struct {
	logger      *log.Logger
	credentials *url.Userinfo
	apiHome     *url.URL
}

func (c *confluenceClient) SearchSpaces(i *ConfluenceSpaceSearchInput) ([]ConfluenceSpace, *ConfluencePaging, error) {
	destination, e := url.Parse(fmt.Sprintf("%s/%s", c.apiHome, constants.SpacesAPIEndpoint))

	if e != nil {
		return nil, nil, e
	}

	if i != nil {
		query := make(url.Values)
		query.Set("limit", fmt.Sprintf("%d", i.Limit))
		query.Set("type", i.Type)
		query.Set("start", fmt.Sprintf("%d", i.Start))
		destination.RawQuery = query.Encode()
	}

	r, e := c.send("GET", destination.String(), nil)

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

	request.Header.Set("Authorization", c.authorization())
	return client.Do(request)
}

func (c *confluenceClient) authorization() string {
	password, _ := c.credentials.Password()
	joined := fmt.Sprintf("%s:%s", c.credentials.Username(), password)
	encoded := base64.StdEncoding.EncodeToString([]byte(joined))
	return fmt.Sprintf("Basic %s", encoded)
}
