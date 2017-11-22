package catelyn

import "net/url"

// GlobalCLIOptions are shared flag values across every catelyn command.
type GlobalCLIOptions struct {
	ConfluenceHost     string
	ConfluenceUsername string
	ConfluencePassword string
}

// UserInfo returns the confluence username and password as a url.Userinfo reference.
func (o *GlobalCLIOptions) UserInfo() *url.Userinfo {
	return url.UserPassword(o.ConfluenceUsername, o.ConfluencePassword)
}
