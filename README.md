# xpack

xpack — a universal Linux package builder written in Go. Convert standalone binaries (.bin) into native Linux packages such as .deb, .rpm, and other distribution formats using simple, repeatable workflows.

Version: 1.1.1-stable

Overview
--------

xpack automates the repetitive parts of packaging Linux software: metadata, file layout, control scripts, signing hooks, and producing multiple output formats from a single source binary. It's intended for maintainers who want a fast, reproducible way to produce .deb, .rpm, and tarball installers from a compiled artifact.

Key features
- Convert a single compiled binary into multiple package formats (.deb, .rpm, tar.gz)
- Produce deterministic and reproducible packages with configurable metadata
- Support for pre/post install/remove scripts, service files, and file permissions
- Build scripts and helpers under `Scripts/` to create packages or installers
- Written in Go and designed to be scriptable and CI-friendly

Quick start
-----------

Build from source and create a package:

1. Build the xpack CLI:

	./Scripts/BinBuilder.sh

One-line installer (recommended for end users)

```bash
curl -fsSL https://raw.githubusercontent.com/nexoral/xpack/main/Scripts/installer.sh | sudo bash -
```

2. Package a binary (example):

	./bin/xpack build --input ./myapp.bin --name myapp --version 1.0.0 --formats deb,rpm

Or run xpack directly (non-interactive) specifying input, arch and version:

```bash
./bin/xpack -i ./bin -arch amd64 -v 1.1.1
```

3. Find output packages under `dist/` (or the path shown by the CLI).

See `INSTALLATION.md` and `LEARN.md` for more details and advanced usage examples.

Where to get help
- Open an issue on the repository for bugs or feature requests
- Read `CONTRIBUTING.md` for contribution guidelines

License
-------
This project is licensed under the terms in the `LICENSE` file.
