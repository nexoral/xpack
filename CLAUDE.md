# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**xpack** - Universal Linux Package Builder

- **Language**: Go ≥1.18
- **Type**: CLI tool for converting binaries to Linux packages (.deb, .rpm, tar.gz)
- **Platform**: Linux (target), cross-platform build
- **Purpose**: Automate creation of native Linux packages from standalone binaries

## Commands

```bash
# Build
go build -o xpack src/Core/main.go
./Scripts/BinBuilder.sh                     # Multi-platform build

# Run
./bin/xpack -i ./bin -arch amd64 -v 1.1.1
./bin/xpack build --input ./myapp.bin --name myapp --version 1.0.0 --formats deb,rpm

# Test
go test ./...
go fmt ./...
go vet ./...
```

## Core Rules (NON-NEGOTIABLE)

1. **ALWAYS test**: Test package creation on Debian and RPM-based systems
2. **ALWAYS build**: Run `go build` after code changes
3. **Reproducible builds**: Same input must produce identical packages
4. **Package standards**: Follow Debian Policy and RPM standards
5. **Error handling**: Clear, actionable error messages
6. **Update docs**: README.md, INSTALLATION.md, LEARN.md when features change

## Critical Constraints

### Package Format Compliance
- **Debian packages (.deb)**: Follow Debian Policy Manual
- **RPM packages (.rpm)**: Follow RPM Packaging Guidelines
- **Tarball (.tar.gz)**: Standard Unix layout
- **Metadata**: Maintainer, description, dependencies, conflicts

### Deterministic Builds
- **Same input = same output** (reproducible builds)
- **No timestamps** in package metadata (or use SOURCE_DATE_EPOCH)
- **Consistent file ordering**
- **Deterministic compression**

## Architecture Overview

See `LEARN.md` for complete details.

### Structure
```
src/
├── Core/           # main.go - CLI entry point
├── base/           # Utilities (banner, helpers)
└── packager/       # Package generation logic
    ├── deb builder
    ├── rpm builder
    └── tar.gz builder

Scripts/
├── BinBuilder.sh   # Multi-platform binary builder
└── installer.sh    # One-line installer script

dist/               # Output directory for generated packages
```

### Package Creation Flow
```
Input Binary (.bin)
    ↓
Parse CLI Arguments (name, version, arch, formats)
    ↓
Generate Package Metadata (control files, spec files)
    ↓
Create Package Structure (directories, permissions)
    ↓
Add Scripts (pre/post install/remove)
    ↓
Build Package (.deb, .rpm, .tar.gz)
    ↓
Output to dist/
```

## Go Standards

### Error Handling
```go
// ✅ GOOD - Descriptive with context
if err != nil {
    return fmt.Errorf("failed to create deb package for '%s': %w", pkgName, err)
}

// ❌ BAD - Generic
if err != nil {
    return err
}
```

### Package Structure
```go
// ✅ GOOD - Clear struct for package metadata
type PackageMetadata struct {
    Name         string
    Version      string
    Architecture string
    Maintainer   string
    Description  string
    Depends      []string
    Conflicts    []string
}

func BuildDeb(meta PackageMetadata, binPath string) error { }

// ❌ BAD - Too many parameters
func BuildDeb(name, version, arch, maintainer, desc string, binPath string) error { }
```

### File Operations
```go
// ✅ GOOD - Use filepath for cross-platform paths
import "path/filepath"

debPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s_%s.deb", name, version, arch))

// ❌ BAD - Hardcoded separators
debPath := outputDir + "/" + name + "_" + version + ".deb"
```

## Key Patterns

### Deterministic File Ordering
```go
// ✅ GOOD - Sort files for reproducibility
sort.Strings(files)
for _, file := range files {
    addToArchive(file)
}

// ❌ BAD - Non-deterministic order
for file := range fileMap {
    addToArchive(file)
}
```

### Metadata Validation
```go
// ✅ Validate package metadata
func ValidateMetadata(meta PackageMetadata) error {
    if meta.Name == "" {
        return fmt.Errorf("package name is required")
    }
    if !regexp.MustCompile(`^[a-z0-9][a-z0-9+.-]+$`).MatchString(meta.Name) {
        return fmt.Errorf("invalid package name: %s", meta.Name)
    }
    if meta.Version == "" {
        return fmt.Errorf("package version is required")
    }
    return nil
}
```

### Package Building
```go
// ✅ GOOD - Build multiple formats
func BuildPackages(meta PackageMetadata, binPath string, formats []string) error {
    for _, format := range formats {
        switch format {
        case "deb":
            if err := BuildDeb(meta, binPath); err != nil {
                return fmt.Errorf("deb build failed: %w", err)
            }
        case "rpm":
            if err := BuildRPM(meta, binPath); err != nil {
                return fmt.Errorf("rpm build failed: %w", err)
            }
        case "tar":
            if err := BuildTarball(meta, binPath); err != nil {
                return fmt.Errorf("tarball build failed: %w", err)
            }
        default:
            return fmt.Errorf("unsupported format: %s", format)
        }
    }
    return nil
}
```

## Documentation Requirements

**Update when features change**:
1. **README.md** - Usage examples, quick start
2. **INSTALLATION.md** - Installation methods, requirements
3. **LEARN.md** - Package creation guide, advanced usage
4. **CONTRIBUTING.md** - Development guidelines
5. **Code comments** - All exported functions, complex logic

## Security

1. **Input Validation**: Validate all CLI arguments (names, versions, paths)
2. **Path Traversal**: Sanitize file paths to prevent directory traversal
3. **Script Injection**: Validate maintainer scripts (pre/post install)
4. **File Permissions**: Set proper permissions (755 for binaries, 644 for docs)
5. **Signature Support**: Support GPG signing for packages

## Testing

- **Unit tests**: Metadata validation, file operations
- **Integration tests**: Full package creation (.deb, .rpm, .tar.gz)
- **Package validation**: Use `lintian` (Debian), `rpmlint` (RPM)
- **Installation tests**: Install packages on target systems
- **Reproducibility**: Verify identical builds from same input

## Package Standards

### Debian (.deb)
- Control file format
- Maintainer scripts (preinst, postinst, prerm, postrm)
- File layout (/usr/bin, /usr/share/doc, /etc)
- Dependencies, conflicts, provides
- Architecture (amd64, arm64, all)

### RPM (.rpm)
- Spec file format
- Scriptlets (%pre, %post, %preun, %postun)
- File layout (/usr/bin, /usr/share/doc, /etc)
- Requires, conflicts, provides
- Architecture (x86_64, aarch64, noarch)

### Tarball (.tar.gz)
- Standard Unix layout
- Install script
- Manifest file
- Checksum file (SHA256)

## Anti-Patterns (FORBIDDEN)

❌ Non-deterministic builds (timestamps, random order)
❌ Invalid package metadata (missing required fields)
❌ Hardcoded paths without cross-platform support
❌ Missing input validation
❌ Breaking package standards (Debian Policy, RPM guidelines)
❌ Skipping package validation tools
❌ Missing error handling
❌ Generic error messages

## Workflow Guidelines

### When Adding Package Format Support
1. Read official packaging guidelines (Debian Policy, RPM docs)
2. Create builder in `src/packager/`
3. Implement metadata generation
4. Add file layout logic
5. Support maintainer scripts
6. Add tests with validation tools
7. Update documentation (README, LEARN)

### When Fixing Bugs
1. Write failing test that reproduces bug
2. Fix the bug
3. Verify test passes
4. Test on target systems (Debian, RPM-based)
5. Validate with lintian/rpmlint
6. Document in CHANGELOG.md

### When Refactoring
1. Ensure backward compatibility
2. Run all tests before and after
3. Validate packages on target systems
4. Update documentation if behavior changes

## Success Criteria

Every task must meet ALL:
- ✅ Builds successfully (`go build`)
- ✅ Tests pass (`go test ./...`)
- ✅ Lints pass (`go vet ./...`, `go fmt ./...`)
- ✅ Packages validate (`lintian`, `rpmlint`)
- ✅ Reproducible builds (same input = same output)
- ✅ Documentation updated
- ✅ Package standards followed
- ✅ Clear error messages
