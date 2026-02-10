# CAPE SDK Generator

Generator workspace for CAPE resources derived from SecAPI.

## Repository layout

- `examples/`: provider usage examples (currently Pulumi; Nitric planned later).
- `ext/`: git submodules, currently the SecAPI reference (`ext/secapi`).
- `mockserver/`: in-memory mock application that behaves like SecAPI.
- `provider/pulumi/`: Pulumi provider implementation for CAPE resources from SecAPI.
- `justfile`: task runner entrypoint for generation, build, install, and local runs.

## Preparations for development

Install these tools first:

- Go
- `just`
- Pulumi CLI
- Docker (needed for mockserver container workflows)

Then initialize submodules:

```bash
just update_modules
```

## Development workflows

Generate/build provider artifacts:

```bash
just build_secapi_spec
just build_pulumi_provider
just build_pulumi_sdk
# or all at once:
just build_pulumi
```

Install local Pulumi plugin for examples:

```bash
just setup_examples
```

Generate/run mockserver:

```bash
just build_mockserver
just run_mockserver
```

Mockserver via Docker:

```bash
just build_mockserver_docker
just run_mockserver_docker
```
