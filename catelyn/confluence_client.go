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
	SearchPages(*ConfluencePageSearchInput) ([]ConfluenceContent, *ConfluencePaging, error)
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
		Logger:      logger,
		credentials: user,
		apiHome:     api,
	}

	return client, nil
}

type confluenceClient struct {
	*log.Logger
	credentials *url.Userinfo
	apiHome     *url.URL
}

func (c *confluenceClient) SearchPages(i *ConfluencePageSearchInput) ([]ConfluenceContent, *ConfluencePaging, error) {
	destination, e := url.Parse(fmt.Sprintf("%s/%s", c.apiHome, constants.ContentAPIEndpoint))

	if e != nil {
		return nil, nil, e
	}

	if i != nil {
		query := make(url.Values)
		query.Set("limit", fmt.Sprintf("%d", i.Limit))
		query.Set("spaceKey", i.SpaceKey)
		query.Set("title", i.Title)
		query.Set("start", fmt.Sprintf("%d", i.Start))
		destination.RawQuery = query.Encode()
	}

	response := struct {
		Results []ConfluenceContent `json:"results"`
	}{}

	if e := c.get(destination, &response); e != nil {
		return nil, nil, e
	}

	return response.Results, nil, nil
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

	response := struct {
		ConfluencePaging
		Results []ConfluenceSpace `json:"results"`
	}{}

	if e := c.get(destination, &response); e != nil {
		return nil, nil, e
	}

	return response.Results, &response.ConfluencePaging, nil
}

func (c *confluenceClient) get(url *url.URL, out interface{}) error {
	r, e := c.send("GET", fmt.Sprintf("%s", url), nil)

	if e != nil {
		return e
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	return decoder.Decode(out)
}

func (c *confluenceClient) send(method string, url string, body io.Reader) (*http.Response, error) {
	client := http.Client{}
	request, e := http.NewRequest("GET", url, nil)

	if e != nil {
		return nil, e
	}

	request.Header.Set("Authorization", c.authorization())

	r, e := client.Do(request)

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("invalid-response")
	}

	return r, nil
}

func (c *confluenceClient) authorization() string {
	password, _ := c.credentials.Password()
	joined := fmt.Sprintf("%s:%s", c.credentials.Username(), password)
	encoded := base64.StdEncoding.EncodeToString([]byte(joined))
	return fmt.Sprintf("Basic %s", encoded)
}
