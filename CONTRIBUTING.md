# Contributing to xpack

First off, thank you for considering contributing to xpack! It's people like you that make xpack such a great tool. This document provides guidelines and steps for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) to understand what behaviors will and will not be tolerated.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for xpack.

Before submitting a bug report:

- Check the [documentation](README.md) to see if there's a solution to your problem.
- Check if the issue has already been reported in our [Issues](https://github.com/nexoral/xpack/issues) section.

When submitting a bug report:

- Use our bug report template.
- Use a clear and descriptive title.
- Describe the exact steps to reproduce the problem.
- Explain the behavior you expected and what you actually observed.
- Include details about your environment (OS, xpack version, Docker version).
- Include screenshots or terminal output if possible.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for xpack.

When submitting an enhancement suggestion:

- Use our feature request template.
- Use a clear and descriptive title.
- Provide a step-by-step description of the suggested enhancement.
- Explain why this enhancement would be useful to xpack users.
- Include any relevant examples or mockups if applicable.

### Your First Code Contribution

Unsure where to begin? Look for issues labeled with:

- `good-first-issue`: Issues suitable for newcomers.
# Contributing to xpack

Thank you for your interest in contributing to xpack. Contributions that improve documentation, add packaging templates, or make the builder more CI-friendly are especially welcome.

## Code of Conduct

Please follow the project's `CODE_OF_CONDUCT.md` when interacting in issues and pull requests.

## How to contribute

### Reporting bugs

- Search existing issues before opening a new one.
- Provide a minimal reproduction, steps taken, and the version of xpack (see `VERSION`).

### Suggesting features

- Open an issue describing the use case, the proposed API/flags, and examples.

### Pull requests

1. Fork the repository and create a branch for your work.
2. Make small, focused changes and keep commits logical.
3. Run `gofmt` and ensure code builds (`go build ./...`).
4. Open a pull request describing the change and any migration notes.

## Developer setup

Prerequisites
- Go 1.18+
- git

Quick start

```bash
git clone https://github.com/nexoral/xpack.git
cd xpack
./Scripts/BinBuilder.sh
./bin/xpack --help
```

## Tests

If you add logic, include small unit tests. Keep them focused and fast. Describe how to run tests in your PR description.

## Style

- Use `gofmt` for formatting.
- Keep functions small and well-documented.

## Thank you

Contributors power open source — thank you for improving xpack!
