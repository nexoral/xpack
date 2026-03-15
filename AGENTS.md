# AGENTS.md

OpenAI Codex CLI Instructions for xpack

## Project Overview

**xpack** - Universal Linux Package Builder

- **Language**: Go ≥1.18
- **Type**: CLI for binary-to-package conversion
- **Formats**: .deb, .rpm, .tar.gz
- **Purpose**: Automate Linux package creation

## Build Commands

```bash
# Development
go build -o xpack src/Core/main.go
./bin/xpack -i ./bin -arch amd64 -v 1.1.1

# Test
go test ./...
go fmt ./...
go vet ./...

# Distribution
./Scripts/BinBuilder.sh    # Multi-platform build
```

## Core Principles

### 1. Reproducible Builds
- Same input must produce identical output
- No timestamps (use SOURCE_DATE_EPOCH)
- Deterministic file ordering
- Consistent compression

### 2. Package Standards
- **Debian**: Follow Debian Policy Manual
- **RPM**: Follow RPM Packaging Guidelines
- **Tarball**: Standard Unix layout

### 3. Validation
- Use `lintian` for .deb packages
- Use `rpmlint` for .rpm packages
- Test installation on target systems

## Architecture

### Structure
```
src/
├── Core/main.go     # CLI entry point
├── base/            # Utilities
└── packager/        # Package builders

Scripts/
├── BinBuilder.sh    # Multi-platform builder
└── installer.sh     # One-line installer

dist/                # Output packages
```

### Flow
```
Binary → CLI Args → Metadata → Package Structure → Build → Validate → Output
```

## Go Standards

### Error Handling
```go
// ✅ GOOD
if err != nil {
    return fmt.Errorf("failed to create package '%s': %w", name, err)
}

// ❌ BAD
if err != nil {
    return err
}
```

### Structs for Options
```go
// ✅ GOOD
type PackageMetadata struct {
    Name         string
    Version      string
    Architecture string
    Maintainer   string
    Description  string
}

func BuildDeb(meta PackageMetadata, binPath string) error

// ❌ BAD
func BuildDeb(name, version, arch, maintainer, desc, binPath string) error
```

## Key Patterns

### Deterministic Sorting
```go
// ✅ Sort for reproducibility
sort.Strings(files)
for _, file := range files {
    addToPackage(file)
}
```

### Metadata Validation
```go
// ✅ Validate package names
if !regexp.MustCompile(`^[a-z0-9][a-z0-9+.-]+$`).MatchString(name) {
    return fmt.Errorf("invalid package name: %s", name)
}
```

## Documentation

Update when features change:
- README.md - Usage, features
- INSTALLATION.md - Install methods
- LEARN.md - Advanced usage, guides
- CHANGELOG.md - Version changes

## Testing

- Unit tests: Metadata, file operations
- Integration tests: Full package creation
- Validation: lintian, rpmlint
- Installation: Test on Debian, RHEL/Fedora

## Package Formats

### Debian (.deb)
- Control file with metadata
- Maintainer scripts (preinst, postinst, prerm, postrm)
- File layout: /usr/bin, /usr/share/doc

### RPM (.rpm)
- Spec file with metadata
- Scriptlets (%pre, %post, %preun, %postun)
- File layout: /usr/bin, /usr/share/doc

### Tarball (.tar.gz)
- Standard Unix layout
- Install script
- Checksum file

## Anti-Patterns

❌ Non-deterministic builds
❌ Invalid package metadata
❌ Missing validation
❌ Breaking package standards
❌ Generic error messages

## Success Criteria

- ✅ `go build` passes
- ✅ `go test ./...` passes
- ✅ Packages validate (lintian, rpmlint)
- ✅ Reproducible builds
- ✅ Documentation updated
- ✅ Standards compliant
