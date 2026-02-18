//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"cape-project.eu/provider/pulumi/internal/codegen"
)

const PulumiControlResourceFile = "pulumi.gen.yaml"
const ProviderTemplatePath = "internal/codegen/provider.tmpl"
const PulumiPluginTemplatePath = "internal/codegen/pulumi_plugin.tmpl"
const ResourceImportBase = "cape-project.eu/provider/pulumi/internal"

var providerTemplate = codegen.ReadTemplate("provider", ProviderTemplatePath)
var pulumiPluginTemplate = codegen.ReadTemplate("pulumiplugin", PulumiPluginTemplatePath)

func main() {
	cwd, _ := os.Getwd()
	controlPath := PulumiControlResourceFile
	if !filepath.IsAbs(controlPath) {
		controlPath = filepath.Join(cwd, controlPath)
	}

	genYaml, err := codegen.GetPulumiGenYaml(controlPath)
	if err != nil {
		fmt.Printf("error reading control resources: %v\n", err)
		return
	}

	writeTemplate(filepath.Join(cwd, "provider.gen.go"), genYaml, providerTemplate)
	writeTemplate(filepath.Join(cwd, "PulumiPlugin.yaml"), genYaml, pulumiPluginTemplate)
}

func writeTemplate(outPath string, data codegen.PulumiGenYaml, tmpl *template.Template) {
	outFile, err := os.Create(outPath)
	if err != nil {
		println(fmt.Errorf("error creating/opening file: %s", err))
		return
	}
	defer func() {
		_ = outFile.Close()
	}()
	if err := tmpl.Execute(outFile, data); err != nil {
		println(fmt.Errorf("error executing template: %s", err))
	}
}
