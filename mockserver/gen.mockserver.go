//go:build ignore

package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"go.yaml.in/yaml/v4"
)

const (
	specGlobPattern       = "../ext/secapi/dist/specs/*.yaml"
	excludedSpecFileName  = "extensions.wellknown.v1.yaml"
	composeTemplatePath   = "compose.tmpl"
	caddyfileTemplatePath = "caddyfile.tmpl"
	composeOutputPath     = "docker-compose.yaml"
	caddyfileOutputPath   = "Caddyfile"
)

type openAPISpec struct {
	Servers []openAPIServer `yaml:"servers"`
}

type openAPIServer struct {
	URL string `yaml:"url"`
}

type server struct {
	Name   string
	File   string
	Prefix string
}

type templateData struct {
	Servers []server
}

func main() {
	servers, err := collectServers(specGlobPattern)
	if err != nil {
		fatal(err)
	}

	data := templateData{Servers: servers}
	if err := renderTemplateFile(composeTemplatePath, composeOutputPath, data); err != nil {
		fatal(err)
	}
	if err := renderTemplateFile(caddyfileTemplatePath, caddyfileOutputPath, data); err != nil {
		fatal(err)
	}
}

func collectServers(pattern string) ([]server, error) {
	specFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob spec files: %w", err)
	}
	if len(specFiles) == 0 {
		return nil, fmt.Errorf("no spec files matched %q", pattern)
	}

	sort.Strings(specFiles)

	servers := make([]server, 0, len(specFiles))
	for _, specFile := range specFiles {
		if filepath.Base(specFile) == excludedSpecFileName {
			continue
		}
		srv, err := serverFromSpec(specFile)
		if err != nil {
			return nil, err
		}
		servers = append(servers, srv)
	}

	return servers, nil
}

func serverFromSpec(specPath string) (server, error) {
	content, err := os.ReadFile(specPath)
	if err != nil {
		return server{}, fmt.Errorf("failed reading %q: %w", specPath, err)
	}

	var spec openAPISpec
	if err := yaml.Unmarshal(content, &spec); err != nil {
		return server{}, fmt.Errorf("failed parsing %q: %w", specPath, err)
	}
	if len(spec.Servers) == 0 {
		return server{}, fmt.Errorf("%q has no servers section", specPath)
	}
	if spec.Servers[0].URL == "" {
		return server{}, fmt.Errorf("%q has empty first server url", specPath)
	}

	fileName := filepath.Base(specPath)
	name, err := serviceNameFromFileName(fileName)
	if err != nil {
		return server{}, fmt.Errorf("failed deriving service name for %q: %w", specPath, err)
	}

	prefix, err := prefixFromURL(spec.Servers[0].URL)
	if err != nil {
		return server{}, fmt.Errorf("failed deriving prefix for %q: %w", specPath, err)
	}

	return server{
		Name:   name,
		File:   fileName,
		Prefix: prefix,
	}, nil
}

func serviceNameFromFileName(fileName string) (string, error) {
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("file name %q must include at least two dot-separated segments", fileName)
	}
	if parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("file name %q has empty service segments", fileName)
	}
	return parts[0] + "_" + parts[1], nil
}

func prefixFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL %q: %w", rawURL, err)
	}

	prefix := parsedURL.Path
	if prefix == "/" {
		return "", nil
	}
	if prefix != "" && !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return prefix, nil
}

func renderTemplateFile(templatePath string, outputPath string, data templateData) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed parsing template %q: %w", templatePath, err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed creating %q: %w", outputPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed rendering %q: %w", outputPath, err)
	}

	return nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
