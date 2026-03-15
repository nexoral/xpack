# GitHub Copilot Instructions for xpack

## Project Overview

**xpack** - Universal Linux Package Builder
- **Language**: Go ≥1.18
- **Type**: CLI for binary-to-package conversion
- **Formats**: .deb, .rpm, .tar.gz

## Core Rules

### 1. Reproducible Builds
Same input must produce identical output

### 2. Package Standards
- Debian: Follow Debian Policy
- RPM: Follow RPM Guidelines

### 3. Validation
- Use `lintian` for .deb
- Use `rpmlint` for .rpm

## Commands

```bash
go build -o xpack src/Core/main.go
go test ./...
./bin/xpack -i ./bin -arch amd64 -v 1.1.1
```

## Go Standards

### Error Handling
```go
// ✅ GOOD
if err != nil {
    return fmt.Errorf("failed to create deb '%s': %w", name, err)
}
```

### Type Safety
```go
// ✅ Use structs
type PackageMetadata struct {
    Name    string
    Version string
}
```

## Key Patterns

### Deterministic Builds
```go
// ✅ Sort files
sort.Strings(files)
```

### Validation
```go
// ✅ Validate names
if !regexp.MustCompile(`^[a-z0-9][a-z0-9+.-]+$`).MatchString(name) {
    return fmt.Errorf("invalid name")
}
```

## Testing

- Unit: Metadata, file ops
- Integration: Package creation
- Validation: lintian, rpmlint

## Success Criteria

- ✅ `go build` passes
- ✅ Tests pass
- ✅ Packages validate
- ✅ Reproducible
- ✅ Docs updated
