# xpack Roadmap

This roadmap lists planned work for xpack. Items are grouped by near-term, mid-term, and long-term goals.

Near-term (v1.2.x)
- Add support for packaging multiple binaries in one archive
- Improve CLI UX: clearer flags, help text, and examples
- Add automatic generation of systemd unit templates
- Add signing support hooks for .deb/.rpm packages

Mid-term (v2.x)
- Add plugin hooks to extend output formats (snap, AppImage)
- Provide a small HTTP API for remote packaging requests in CI
- Improved templating for package metadata and changelogs
- Official GitHub Actions workflow to produce artifacts

Long-term (v3.x)
- Repository hosting helpers for publishing packages to packagecloud, apt repos, or private RPM repos
- Windows and macOS helper tooling (cross-platform packaging guidance)
- Marketplace for community-contributed packaging templates

How to influence the roadmap
- Open an issue and label it `enhancement`
- Propose a design in a draft PR for discussion

This file is a living plan and will be updated as priorities evolve.
