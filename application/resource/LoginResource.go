package resource

type Resource struct {
	ErrorMessage  string
	UserName      string
	Domains       []string
	Redirect      string
	AccessToken   string
	DisableLoader bool
	FromURL       string
}

func (rsc *Resource) IsSetCookies() bool {
	return len(rsc.Domains) > 0
}
