package utils

import (
	"context"
	"reflect"

	"cape-project.eu/sdk-generator/provider/pulumi/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

func GetTenantFromImputs(ctx context.Context, inputs any) string {
	if _, ok := reflect.TypeOf(inputs).FieldByName("Tenant"); ok {
		return reflect.ValueOf(inputs).FieldByName("Tenant").String()
	}
	config := infer.GetConfig[config.Config](ctx)
	return config.Tenant
}

func GetWorkspaceFromInputs(ctx context.Context, inputs any) (string, bool) {
	if _, ok := reflect.TypeOf(inputs).FieldByName("Workspace"); ok {
		return reflect.ValueOf(inputs).FieldByName("Workspace").String(), true
	}
	config := infer.GetConfig[config.Config](ctx)
	if config.Workspace != nil {
		return *config.Workspace, true
	}

	return "", false
}
