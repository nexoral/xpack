# Installation Guide

This document explains how to build and install xpack locally or on CI.

Requirements
- A Linux system (Debian, Ubuntu, Fedora, CentOS, etc.)
- Go 1.18+ for building from source
- Basic build tools (bash, tar, dpkg-buildpackage if creating .deb locally)

Build from source (recommended for developers)

1. Build the CLI using the included script:

   ./Scripts/BinBuilder.sh

   On success the `bin/` directory will contain the `xpack` executable.

2. Optionally move the binary to a system path:

   sudo mv bin/xpack /usr/local/bin/xpack
   sudo chmod +x /usr/local/bin/xpack

One-line installer (recommended for end users)

Use the repository installer script to install xpack with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/nexoral/xpack/main/Scripts/installer.sh | sudo bash -
```

Create packages using the project scripts

Use `PackageBuilder.sh` and `BinBuilder.sh` under `Scripts/` to produce distribution packages. The scripts are opinionated but designed to be easy to inspect and adapt for CI.

Example: produce .deb and .rpm from a binary

1. Place your input binary (for example `myapp.bin`) into a working folder.
2. Run xpack (or the scripts) to create packages:

   ./bin/xpack build --input ./myapp.bin --name myapp --version 1.0.0 --formats deb,rpm

3. Output packages will be written to `dist/` or the path printed by the tool.

Package installation (end user)

Debian/Ubuntu (.deb)

   sudo dpkg -i myapp_1.0.0_amd64.deb
   sudo apt-get install -f  # fix missing deps if necessary

Fedora/RPM-based (.rpm)

   sudo rpm -Uvh myapp-1.0.0.x86_64.rpm

Tarball

   tar -xzf myapp-1.0.0-linux-amd64.tar.gz
   sudo mv myapp /usr/local/bin/

Using in CI

The scripts are written to be CI-friendly: set inputs via environment variables and run the scripts in a container or runner to produce artifacts. Keep artifacts under `dist/` and upload as CI job artifacts.

Uninstall

Remove the installed binary or package via the platform package manager or by removing `/usr/local/bin/xpack` if installed manually.

Troubleshooting

- If a script fails due to permissions, ensure executable bits are set: `chmod +x Scripts/*.sh`.
- If package metadata is incorrect, edit the template files under `Scripts/` before running the builder.
