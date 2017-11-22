package constants

const (
	// APIUserAgentString is the user agent send when communicating w/ confluence.
	APIUserAgentString = "catelyn-cli/0.0.1 (command line tool written in golang)"

	// ContentAPIEndpoint see https://docs.atlassian.com/atlassian-confluence/REST/latest-server/#content-getContent
	ContentAPIEndpoint = "/wiki/rest/api/content"

	// SpacesAPIEndpoint see https://developer.atlassian.com/cloud/confluence/rest/#api-space-get
	SpacesAPIEndpoint = "/wiki/rest/api/space"
)
