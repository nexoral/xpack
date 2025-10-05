git clone https://github.com/nexoral/xpack.git
docker logs [container-name]
# xpack - Learning Guide

## Introduction

xpack is a minimal, focused CLI that turns compiled Linux binaries into native distribution packages (.deb, .rpm, and tarballs). This guide explains key concepts, workflows, and examples to help you produce consistent packages for distribution or CI pipelines.

## Core concepts

- Input artifact: a single compiled binary (for example `myapp.bin` or `myapp` built for linux/amd64).
- Metadata: name, version, maintainer, description, license, and target formats.
- Hooks: optional pre/post install or remove scripts and service unit files for systemd.
- Output formats: .deb, .rpm, and plain tar.gz/zip artifacts.

## Basic workflow

1. Build or obtain your Linux binary (statically linked preferred for simplicity).
2. Run xpack to generate packages and installers.
3. Test-install the generated packages on target distributions.
4. Distribute artifacts via releases, package repos, or direct download links.

Example: package a binary into deb and rpm

1. Ensure you have the xpack executable under `bin/`:

   ./Scripts/BinBuilder.sh

2. Run the build command:

   ./bin/xpack build --input ./myapp --name myapp --version 1.0.0 --formats deb,rpm

3. Install and test the package on a VM or container matching the target distro.

## Packaging options

- --name: package name
- --version: version string (semver recommended)
- --maintainer: contact for package metadata
- --formats: comma-separated list of formats (deb,rpm,tar)
- --install-script / --preinst / --postinst: paths to control scripts
- --service-file: path to a systemd unit file to include and enable

## Reproducible builds and CI considerations

- Keep metadata in a simple config or a small Makefile to ensure repeatable package builds.
- Pin builder container images for CI to ensure consistent tooling across runs.
- Use the scripts in `Scripts/` as a baseline and adapt them for your CI provider (GitHub Actions, GitLab CI, etc.).

## Architecture (high level)

- CLI front-end: argument parsing and subcommands
- Packaging engine: assembles package file structure and control metadata
- Format adapters: code paths that write .deb control files, .rpm spec files, and tarball layout
- Script templates: small templates for common hooks and systemd units

## Testing packages

- Use lightweight containers (Debian/Ubuntu/Fedora) to install generated packages and run smoke tests.
- Validate files, permissions, and service activation (if applicable).

## Troubleshooting

- Incorrect file permissions in packages: verify file modes before packaging.
- Missing dependencies for .deb/.rpm: add dependency metadata or prefer static binaries.
- Service files not enabling: ensure unit files are installed to `/lib/systemd/system` and use `systemctl daemon-reload` during postinst.

## Further reading

- Debian packaging basics: https://www.debian.org/doc/manuals/maint-guide/
- RPM Packaging Guide: https://rpm.org/
