package gowl

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type Config struct {
	Addr string `json:"addr"`

	EnableTLS bool   `json:"enable_tls"`
	CertFile  string `json:"cert_file"`
	KeyFile   string `json:"key_file"`

	ServerName string `json:"server_name"`

	NotFoundHandler         Handler `json:"-"`
	MethodNotAllowedHandler Handler `json:"-"`

	HandleOptions          bool `json:"handle_options"`
	HandleMethodNotAllowed bool `json:"handle_method_not_allowed"`
	RedirectTrailingSlash  bool `json:"redirect_trailing_slash"`

	RedirectUpperCasePath bool `json:"redirect_upper_case_path"`

	TemplatePath    string `json:"template_path"`
	TemplateFileExt string `json:"template_file_ext"`

	TemplateFunc template.FuncMap `json:"-"`

	NegotiateDefaultOffer bool `json:"negotiate_default_offer"`
}

func (c *Config) String() string {
	str := "Addr: " + c.Addr + "\n"
	str += fmt.Sprintf("Enable TLS: %t\n", c.EnableTLS)
	if c.EnableTLS {
		str += "Cert file: " + c.CertFile + "\n"
		str += "Key file: " + c.KeyFile + "\n"
	}
	str += "Server name: " + c.ServerName + "\n"
	str += fmt.Sprintf("Handle OPTIONS: %t\n", c.HandleOptions)
	str += fmt.Sprintf("Handle method not allowed: %t\n", c.HandleMethodNotAllowed)
	str += fmt.Sprintf("Redirect trailing slash: %t\n", c.RedirectTrailingSlash)
	str += fmt.Sprintf("Redirect upper case path: %t\n", c.RedirectUpperCasePath)
	str += "Template path: " + c.TemplatePath + "\n"
	str += "Template file extension: " + c.TemplateFileExt + "\n"
	str += fmt.Sprintf("Negotiate default offer: %t", c.NegotiateDefaultOffer)
	return str
}

func NewConfig() *Config {
	c := new(Config)
	c.Addr = ":8000"
	c.ServerName = ServerName
	c.HandleOptions = true
	c.HandleMethodNotAllowed = true
	c.RedirectTrailingSlash = true
	c.RedirectUpperCasePath = true
	c.TemplatePath = filepath.Join(execPath, "templates")
	c.TemplateFileExt = ".html"
	return c
}

func LoadConfig(filename string) (config *Config, err error) {
	var f *os.File

	config = NewConfig()
	f, err = os.Open(filename)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&config)
	return
}
