# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-02-17

### Added
- Embedded filesystem for banner files using Go's `embed` package
  - Banner files bundled into binary at compile time
  - Binary is now fully relocatable and self-contained
  - No external testdata files required at runtime
- `GetBannerFS()` function in main package to expose embedded filesystem
- Binary integration test (`TestBuiltBinary_FromRepoRoot`) to validate relocatability
- Color package (`internal/color`) for parsing color specifications
  - Named colors: red, green, blue, yellow, cyan, magenta, white, black, orange, purple, pink, brown, gray
  - Hex format: `#RRGGBB`
  - RGB format: `rgb(R,G,B)`
  - `Parse()` function to convert color specs to RGB values
  - `ANSI()` function to generate 24-bit ANSI escape sequences
- Coloring package (`internal/coloring`) for applying ANSI colors to ASCII art
  - `ApplyColor()` function for full-text and substring coloring
  - Accurate column mapping using character widths
  - Support for non-contiguous and overlapping substring matches
- Flagparser package (`internal/flagparser`) for CLI argument validation
  - `ParseArgs()` function to validate `--color=` flag syntax
  - Argument count, flag position, and empty value checks
- `--color` CLI flag for colored ASCII art output
  - Full text coloring: `--color=red "text"`
  - Substring coloring: `--color=red substring "text"`
  - Works with all banner styles
- `CharWidths()` function in parser package for per-character column widths
- Color mode routing in main (`hasColorFlag`, `runColorMode`, `extractColorArgs`)
- Exit code 4 for color parse errors
- Comprehensive integration tests for color mode (18 test cases)
- Unit tests for all new functions and packages
- `run-color` Makefile target
- CI workflow (`.github/workflows/ci.yml`) with test, lint, and build jobs
  - Test matrix: Go 1.21/1.22 on Ubuntu, macOS, and Windows
  - Lint: golangci-lint v2.1.6
  - Build verification on Ubuntu
- Release workflow (`.github/workflows/release.yml`) for automated binary distribution
  - Triggered on `v*` tags
  - Builds cross-platform binaries (Linux, macOS, Windows)
  - Creates GitHub Release with all binaries attached
- CI status badge in README.md
- Mermaid architecture diagrams in `diagrams/` folder
  - Architecture overview, flowchart, class diagram, and sequence diagram
  - Renders natively on GitHub

### Changed
- `parser.LoadBanner()` signature changed from `LoadBanner(path string)` to `LoadBanner(fsys fs.FS, path string)`
  - Accepts any `fs.FS` implementation (embedded, disk, testing, etc.)
  - Enables dependency injection for filesystem operations
- `parser.readLines()` now uses `fs.ReadFile()` instead of `os.Open()`
- Test files updated to use `os.DirFS()` for reading real files from disk
- Architecture diagrams updated to show new `fs.FS` parameter
- Project restructured to `cmd/internal` layout
  - Main package moved to `cmd/ascii-art/`
  - Internal packages moved to `internal/`
- Renamed renderer exported API from `RendererASCII()` to `ASCII()` and updated all call sites/tests/docs
- Updated all documentation for v1.1.0
  - README.md with color usage examples and correct project structure
  - AGENTS.md with new packages, exit codes, and commands
  - CONTRIBUTING.md with updated structure and scopes
  - Makefile paths updated for `cmd/ascii-art/` layout
- Function documentation standardized with Parameters/Returns sections
- Test files given package-level doc comments
- Fixed goimports ordering in color_test.go and coloring_test.go

---

## [1.0.0] - 2026-01-10

**Final Release** - Complete implementation.

### Added
- Parser package for reading and parsing ASCII art banner files
  - `LoadBanner()` function to read and parse banner files
  - `readLines()` helper for file reading
  - `buildBanner()` helper to construct Banner map
  - Comprehensive error handling with wrapped errors
  - Security annotations for file operations
- Renderer package for converting text to ASCII art
  - `ASCII()` main rendering function with newline support
  - `validateBannerCharacters()` for character validation
  - `validateInput()` for input validation
  - Efficient string building with `strings.Builder`
- Comprehensive test suite
  - 12 parser unit tests with 93.9% coverage
  - 14 renderer unit tests with 97.1% coverage
  - 8 main package unit tests for CLI argument parsing
  - End-to-end integration tests for full application stack
  - Table-driven tests for multiple scenarios
- Professional Makefile with build automation
  - Quality control targets (fmt, vet, lint, check)
  - Development targets (run, build, install)
  - Testing targets (test, coverage)
  - Cross-compilation for Linux, macOS, Windows
  - Utility targets (tidy, version, help, clean)
- Comprehensive documentation
  - Professional README.md with usage examples
  - AGENTS.md for AI coding agents
  - CONTRIBUTING.md with development guidelines
  - PERMISSIONS.md for team workflow
- golangci-lint v2 configuration
  - 10 enabled linters (dupl, goconst, gocyclo, gosec, misspell, prealloc, revive, staticcheck, unconvert, unparam)
  - Comprehensive code quality checks with errcheck, govet, gocritic settings
  - Test-specific exclusions for cleaner test code
  - Formatters (gofmt, goimports) for consistent code style
- Support for three banner styles (standard, shadow, thinkertoy)
- Command-line interface for text to ASCII art conversion
- Cross-platform compatibility (Linux, macOS, Windows)
- Zero external dependencies

### Changed
- Applied Go best practices throughout codebase
  - Package-level documentation
  - Constants for magic numbers
  - Error wrapping with `fmt.Errorf` and `%w` verb
  - Lowercase error messages (Go style guide)

### Fixed
- All linting errors and warnings
- Test coverage for edge cases
- Code formatting consistency
- Newline handling in renderer (now properly renders `\n` in input text)

### Performance
- Parser optimized for fast banner loading
- Renderer optimized with efficient string building
- Linear scaling O(n) with text length

---

## Release Guidelines

### Version Format
This project uses [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for added functionality (backwards compatible)
- **PATCH** version for backwards compatible bug fixes

### Release Types

#### Major Release (X.0.0)
- Breaking changes to CLI interface
- Incompatible changes to banner file format
- Removal of features or flags

#### Minor Release (0.X.0)
- New banner styles
- New CLI flags or options (backwards compatible)
- Performance improvements
- New features

#### Patch Release (0.0.X)
- Bug fixes
- Documentation updates
- Security patches
- Minor performance improvements

### Unreleased Section
The `[Unreleased]` section tracks changes that are committed but not yet released:
- Use for active development
- Move to appropriate version section on release
- Keep organized by change type

### Change Categories

Changes should be grouped into these categories:

- **Added** - New features
- **Changed** - Changes to existing functionality
- **Deprecated** - Soon-to-be removed features
- **Removed** - Removed features
- **Fixed** - Bug fixes
- **Security** - Security vulnerability fixes
- **Performance** - Performance improvements

### Example Entry Format

```markdown
## [1.2.0] - 2024-01-15

### Added
- New "graffiti" banner style (#42)
- Support for color output with --color flag (#38)
- Verbose mode with -v flag for debugging (#45)

### Changed
- Improved error messages for invalid characters (#40)
- Updated help text with more examples (#43)

### Fixed
- Fixed rendering issue with consecutive newlines (#41)
- Corrected alignment for wide characters (#44)

### Performance
- 20% faster rendering for long texts (#39)
- Reduced memory allocations in parser (#46)
```

---

## Links

- [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
- [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
- [Conventional Commits](https://www.conventionalcommits.org/)
