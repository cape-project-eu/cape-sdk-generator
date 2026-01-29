package convertors

import (
	"cape-project.eu/sdk-generator/provider/pulumi/internal/schemas"
	"cape-project.eu/sdk-generator/provider/pulumi/secapi/models"
)

// goverter:variables
// goverter:output:format assign-variable
var (
	ConvertReferenceURNToOpenAPI    func(schemas.ReferenceURN) models.ReferenceURN
	ConvertReferenceObjectToOpenAPI func(schemas.ReferenceObject) models.ReferenceObject
	ConvertReferenceURNToPulumi     func(models.ReferenceURN) schemas.ReferenceURN
	ConvertReferenceObjectToPulumi  func(models.ReferenceObject) schemas.ReferenceObject
)

func ConvertReferenceToOpenAPI(in schemas.Reference) models.Reference {
	var out models.Reference
	if in.ReferenceURN != nil {
		out.FromReferenceURN(ConvertReferenceURNToOpenAPI(*in.ReferenceURN))
	}
	if in.ReferenceObject != nil {
		out.FromReferenceObject(ConvertReferenceObjectToOpenAPI(*in.ReferenceObject))
	}
	return out
}

func ConvertReferenceToPulumi(in models.Reference) schemas.Reference {
	var out schemas.Reference
	if a, err := in.AsReferenceURN(); err != nil {
		r := ConvertReferenceURNToPulumi(a)
		out.ReferenceURN = &r
	}

	if b, err := in.AsReferenceObject(); err != nil {
		r := ConvertReferenceObjectToPulumi(b)
		out.ReferenceObject = &r
	}

	return out
}
