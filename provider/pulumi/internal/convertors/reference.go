package convertors

import (
	"cape-project.eu/provider/pulumi/internal/schemas"
	"cape-project.eu/provider/pulumi/secapi/models"
)

// goverter:variables
// goverter:output:format assign-variable
var (
	ConvertReferenceURNToOpenAPI    func(schemas.ReferenceURN) models.ReferenceURN
	ConvertReferenceObjectToOpenAPI func(schemas.ReferenceObject) models.ReferenceObject
	ConvertReferenceURNToPulumi     func(models.ReferenceURN) schemas.ReferenceURN
	ConvertReferenceObjectToPulumi  func(models.ReferenceObject) schemas.ReferenceObject
)
