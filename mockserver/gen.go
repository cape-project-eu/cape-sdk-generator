package main

//go:generate find . -name "*.gen.go" -not -name "gen.go" -delete
//go:generate sh -c "cd models && ./gen_models.sh"
//go:generate ./gen_stubs.sh
