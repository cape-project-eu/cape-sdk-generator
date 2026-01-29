package config

import "github.com/pulumi/pulumi-go-provider/infer"

type Config struct {
	BaseURL   string  `pulumi:"baseURL"`
	AuthToken *string `pulumi:"authToken,optional" provider:"secret"`
	Tenant    string  `pulumi:"tenant"`
	Workspace *string `pulumi:"workspace,optional"`
}

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.BaseURL, "BaseURL defines the server url for API communication.")
	a.Describe(&c.AuthToken, "AuthToken is the bearer token that is attached to API calls.")
	a.Describe(&c.Tenant, "Tenant defines the default tenant used for all API calls. May be overwritten in specific calls.")
	a.Describe(&c.Workspace, "Workspace defines a default workspace for all API calls. Can be omitted and given to all objects, or specifically overwritten for calls.")
}
