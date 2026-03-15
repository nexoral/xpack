# GEMINI.md

This file provides guidance to Gemini Code Assist when working with code in this repository.

## Project Overview

**xpack** - Universal Linux Package Builder

- **Language**: Go ≥1.18
- **Type**: CLI for converting binaries to .deb, .rpm, .tar.gz packages
- **Platform**: Linux (target)
- **Purpose**: Automate Linux package creation

## Commands

```bash
# Build & Run
go build -o xpack src/Core/main.go
./bin/xpack -i ./bin -arch amd64 -v 1.1.1
./Scripts/BinBuilder.sh

# Test
go test ./...
go fmt ./...
go vet ./...
```

## Core Rules (NON-NEGOTIABLE)

1. **Reproducible builds**: Same input = same output
2. **Package standards**: Follow Debian Policy, RPM guidelines
3. **ALWAYS test**: Validate with lintian (deb), rpmlint (rpm)
4. **ALWAYS build**: Run `go build` after changes
5. **Update docs**: README, INSTALLATION, LEARN

## Architecture

### Structure
```
src/
├── Core/main.go     # CLI entry
├── base/            # Utilities
└── packager/        # Package builders (deb, rpm, tar.gz)
```

### Flow
```
Binary Input → Parse Args → Generate Metadata → Build Package → Output to dist/
```

## Go Standards

### Error Handling
```go
// ✅ GOOD
if err != nil {
    return fmt.Errorf("failed to create deb for '%s': %w", name, err)
}

// ❌ BAD
if err != nil {
    return err
}
```

### Type Safety
```go
// ✅ Use structs
type PackageMetadata struct {
    Name    string
    Version string
    Arch    string
}

// ❌ Too many params
func BuildDeb(name, version, arch string) error
```

## Key Patterns

### Deterministic Builds
```go
// ✅ Sort for reproducibility
sort.Strings(files)

// ❌ Non-deterministic
for file := range fileMap { }
```

### Validation
```go
// ✅ Validate metadata
if !regexp.MustCompile(`^[a-z0-9][a-z0-9+.-]+$`).MatchString(name) {
    return fmt.Errorf("invalid package name")
}
```

## Documentation

Update when features change:
- README.md - Usage, quick start
- INSTALLATION.md - Install methods
- LEARN.md - Package creation guide

## Testing

- Unit tests for metadata, file ops
- Integration tests for package creation
- Validate with lintian, rpmlint
- Test on Debian, RPM-based systems

## Package Standards

- **Debian**: Control file, maintainer scripts, /usr/bin layout
- **RPM**: Spec file, scriptlets, /usr/bin layout
- **Tarball**: Standard Unix layout, install script

## Definition of "Done"

- ✅ `go build` passes
- ✅ `go test ./...` passes
- ✅ Packages validate (lintian, rpmlint)
- ✅ Reproducible builds
- ✅ Documentation updated
