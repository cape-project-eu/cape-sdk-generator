package convertors

import (
	"cape-project.eu/provider/pulumi/internal/schemas"
	"cape-project.eu/provider/pulumi/secapi/models"
)

// goverter:variables
// goverter:output:format assign-variable
var (
	ConvertAnnotationsToOpenAPI func(schemas.Annotations) models.Annotations
	ConvertLabelsToOpenAPI      func(schemas.Labels) models.Labels
	ConvertExtensionsToOpenAPI  func(schemas.Extensions) models.Extensions
	ConvertAnnotationsToPulumi  func(models.Annotations) schemas.Annotations
	ConvertLabelsToPulumi       func(models.Labels) schemas.Labels
	ConvertExtensionsToPulumi   func(models.Extensions) schemas.Extensions
)
