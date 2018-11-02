package gowl

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
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
}

func (c *Config) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "Addr: %s\n", c.Addr)
	fmt.Fprintf(buf, "Enable TLS: %t\n", c.EnableTLS)
	if c.EnableTLS {
		fmt.Fprintf(buf, "Cert file: %s\n", c.CertFile)
		fmt.Fprintf(buf, "Key file: %s\n", c.KeyFile)
	}
	fmt.Fprintf(buf, "Server name: %s\n", c.ServerName)
	fmt.Fprintf(buf, "Handle OPTIONS: %t\n", c.HandleOptions)
	fmt.Fprintf(buf, "Handle method not allowed: %t\n", c.HandleMethodNotAllowed)
	fmt.Fprintf(buf, "Redirect trailing slash: %t\n", c.RedirectTrailingSlash)
	fmt.Fprintf(buf, "Redirect upper case path: %t\n", c.RedirectUpperCasePath)
	fmt.Fprintf(buf, "Template path: %s\n", c.TemplatePath)
	fmt.Fprintf(buf, "Template file extension: %s\n", c.TemplateFileExt)
	return buf.String()
}

func NewConfig() *Config {
	return &Config{
		Addr:                   ":8000",
		ServerName:             ServerName,
		HandleOptions:          true,
		HandleMethodNotAllowed: true,
		RedirectTrailingSlash:  true,
		RedirectUpperCasePath:  true,
		TemplatePath:           filepath.Join(execPath, "templates"),
		TemplateFileExt:        ".html",
	}
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
