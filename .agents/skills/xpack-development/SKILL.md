---
name: xpack-development
description: Development rules for xpack Linux package builder
version: 1.0.0
tags: [golang, packaging, deb, rpm, linux]
author: xpack Team
---

# xpack Development Skill

## Project Identity

**xpack** - Universal Linux Package Builder
- Go ≥1.18
- CLI tool for converting binaries to .deb, .rpm, .tar.gz
- Reproducible, standards-compliant packages

## Mandatory Workflows

### After EVERY Code Change
```bash
go build -o xpack src/Core/main.go
go test ./...
go fmt ./...
go vet ./...
```

### For ANY Package Format Change
1. Validate with lintian (deb) or rpmlint (rpm)
2. Test installation on target system
3. Verify reproducible builds
4. Update documentation

## Definition of "Done"

- ✅ `go build` passes
- ✅ `go test ./...` passes
- ✅ Packages validate (lintian/rpmlint)
- ✅ Reproducible builds verified
- ✅ Documentation updated
- ✅ Package standards followed

## Architecture

```
src/
├── Core/main.go     # CLI entry
├── base/            # Utilities
└── packager/        # Package builders

Scripts/
├── BinBuilder.sh    # Multi-platform
└── installer.sh     # One-line install

dist/                # Output packages
```

## Go Standards

### Error Handling
```go
// ✅ REQUIRED
if err != nil {
    return fmt.Errorf("failed to create deb for '%s': %w", name, err)
}
```

### Type Safety
```go
// ✅ Use structs
type PackageMetadata struct {
    Name         string
    Version      string
    Architecture string
    Maintainer   string
    Description  string
}
```

## Key Patterns

### Reproducible Builds
```go
// ✅ Sort for determinism
sort.Strings(files)
for _, file := range files {
    addToPackage(file)
}
```

### Package Validation
```go
// ✅ Validate metadata
if !regexp.MustCompile(`^[a-z0-9][a-z0-9+.-]+$`).MatchString(name) {
    return fmt.Errorf("invalid package name")
}
```

## Package Standards

- **Debian**: Control file, maintainer scripts, /usr/bin layout
- **RPM**: Spec file, scriptlets, /usr/bin layout
- **Tarball**: Unix layout, install script, checksums

## Testing

- Unit: Metadata validation, file operations
- Integration: Full package creation
- Validation: lintian (deb), rpmlint (rpm)
- Installation: Test on Debian, RPM-based systems

## Anti-Patterns (FORBIDDEN)

❌ Non-deterministic builds
❌ Invalid metadata
❌ Missing validation
❌ Breaking standards
❌ Generic errors

## Success Criteria

- ✅ Builds successfully
- ✅ Tests pass
- ✅ Packages validate
- ✅ Reproducible
- ✅ Standards compliant
