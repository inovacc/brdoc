# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Documentation updates across the project:
  - README: corrected API signatures, updated Go version requirement (1.24), added CLI bulk operations (`--from`, `--count`) and stdin examples, updated project structure, and roadmap status for CLI.
  - SETUP: aligned project structure and added CLI quick start with bulk validation examples.
  - Package docs (`doc.go`): corrected CNPJ `Generate()` example to reflect current signature.
  - Clarified expected validation output format in CLI docs (`valid\t<FORMATTED_VALUE>` / `invalid\t<ORIGINAL_INPUT>`).
  - Testing docs (README, SETUP, CONTRIBUTING): documented usage of `testify` (`assert`/`require`) with examples.

### Tests

- Migrated unit tests to use `github.com/stretchr/testify` assertion library (`assert`/`require`) for clearer and more concise assertions.

### Planned

- Support for legacy numeric-only CNPJ validation
- RG (Registro Geral) validation
- CNH (Carteira Nacional de Habilitação) validation
- PIS/PASEP validation
- CEP (postal code) validation and formatting

## [0.1.0] - 2024-11-19

### Added

- Initial release
- CPF validation with check digit verification
- CPF generation with valid check digits
- CPF formatting (XXX.XXX.XXX-XX)
- CPF state/region identification based on 9th digit
- Detection of invalid CPF patterns (all same digits)
- CNPJ alphanumeric validation per SERPRO specification
- CNPJ alphanumeric generation
- CNPJ formatting (XX.XXX.XXX/XXXX-XX)
- Modulo 11 check digit calculation for CNPJ
- Auto-detection of document type (CPF/CNPJ)
- Comprehensive test suite with 95%+ coverage
- Benchmark suite for performance testing
- Thread-safe random number generation
- Zero external dependencies
- Complete API documentation
- Usage examples
- CI/CD pipeline with GitHub Actions

### Technical Details

- Implements official SERPRO algorithm for alphanumeric CNPJ
- Character mapping: 0-9 → 0-9, A-Z → 17-42 (ASCII - 48)
- Weight distribution: 2-9, repeating from right to left
- Special modulo 11 rule: remainder 0 or 1 → check digit 0

## [0.0.1] - 2024-11-15

### Added

- Project initialization
- Basic project structure
- MIT License

---

## Types of Changes

- `Added` for new features
- `Changed` for changes in existing functionality
- `Deprecated` for soon-to-be removed features
- `Removed` for now removed features
- `Fixed` for any bug fixes
- `Security` in case of vulnerabilities

[Unreleased]: https://github.com/inovacc/brdoc/compare/v0.1.0...HEAD

[0.1.0]: https://github.com/inovacc/brdoc/releases/tag/v0.1.0

[0.0.1]: https://github.com/inovacc/brdoc/releases/tag/v0.0.1
