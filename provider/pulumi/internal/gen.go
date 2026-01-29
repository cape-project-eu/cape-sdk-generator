package internal

//go:generate find . -name "*.gen.go" -type f ! -path "./schemas/*" -delete
//go:generate go run gen.controlresources.go
//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen ./...
