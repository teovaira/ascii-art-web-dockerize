# AGENTS.md

Instructions for AI coding agents working on the ascii-art-web-dockerize project.

## Project Overview

Go CLI and web application converting text to ASCII art with three banner styles (standard, shadow, thinkertoy) and optional ANSI 24-bit color support for full text or specific substrings. **Zero external dependencies** — Go standard library only.

## Quick Commands

```bash
# No external dependencies to install (Go standard library only)

# Docker (requires Docker installed)
make docker-build  # Build the Docker image (ascii-art-web-docker)
make docker-run    # Run the container (dockerize) on port 8080
make docker-stop   # Stop and remove the container
make docker-clean  # Stop container + remove image
# Or use the helper script:
./docker-build.sh          # Build image + start container
./docker-build.sh stop     # Stop and remove container
./docker-build.sh clean    # Stop container + remove image

# Build
make build         # CLI binary
make build-web     # Web server binary

# Test (run before every commit)
make test          # All tests
make coverage      # With coverage report

# Quality checks (run before every commit)
make fmt           # Format code (gofmt + goimports)
make vet           # Static analysis
make lint          # Linter checks (golangci-lint)
make check         # All quality checks (fmt + vet + lint)

# Run CLI
cd cmd/ascii-art && go run . "Hello" standard
cd cmd/ascii-art && go run . --color=red "Hello"
cd cmd/ascii-art && go run . --color=red He "Hello"

# Run web server (must run from project root)
make run-web
# or: go run ./cmd/ascii-art-web
```

## Project Structure

```
ascii-art-web-dockerize/
├── .github/
│   └── workflows/
│       ├── ci.yml             # CI workflow (test, lint, build)
│       └── release.yml        # Release workflow (cross-platform binaries)
├── .gitignore                 # Git ignore rules
├── .golangci.yml              # Linter configuration
├── Dockerfile                 # Multi-stage Docker build
├── docker-build.sh            # Helper script: build image + run container
├── LICENSE                    # Project license
├── Makefile                   # Build automation (incl. docker-* targets)
├── go.mod                     # Go module file (no external deps)
├── AGENTS.md                  # This file
├── CHANGELOG.md               # Version history
├── CONTRIBUTING.md            # Contribution guidelines
├── PERMISSIONS.md             # Team permissions
├── README.md                  # User documentation
├── diagrams/                  # Mermaid architecture diagrams
│   ├── architecture.md        # High-level system overview
│   ├── class-diagram.md       # Package types and relationships
│   ├── flowchart.md           # Program execution flow
│   └── sequence-diagram.md    # Color mode call sequence
├── cmd/
│   ├── ascii-art/             # CLI entry point
│   │   ├── main.go
│   │   ├── args.go
│   │   ├── banner.go
│   │   ├── color_mode.go
│   │   ├── main_test.go
│   │   ├── integration_test.go
│   │   └── testdata/          # Banner files and test fixtures
│   │       ├── standard.txt
│   │       ├── shadow.txt
│   │       ├── thinkertoy.txt
│   │       ├── corrupted.txt  # Test fixture
│   │       ├── empty.txt      # Test fixture
│   │       └── oversized.txt  # Test fixture
│   └── ascii-art-web/         # Web server entry point
│       ├── main.go
│       └── integration_test.go
├── static/                    # Static web assets
│   ├── style.css
│   └── favicon files
├── templates/                 # HTML templates
│   ├── base.html
│   └── index.html
└── internal/
    ├── banners/               # Embedded banner files (//go:embed *.txt)
    │   ├── banners.go
    │   ├── standard.txt
    │   ├── shadow.txt
    │   └── thinkertoy.txt
    ├── color/                 # Color specification parsing
    │   ├── color.go
    │   └── color_test.go
    ├── coloring/              # ANSI color application to ASCII art
    │   ├── coloring.go
    │   └── coloring_test.go
    ├── flagparser/            # CLI argument validation
    │   ├── flagparser.go
    │   └── flagparser_test.go
    ├── handlers/              # HTTP handlers, ASCII generation, template cache
    │   ├── handlers.go
    │   ├── handlers_test.go
    │   └── template_cache.go
    ├── parser/                # Banner file parsing (from fs.FS)
    │   ├── banner_parser.go
    │   └── parser_test.go
    ├── renderer/              # ASCII art rendering
    │   ├── renderer.go
    │   └── renderer_test.go
    └── validation/            # Web input validation
        ├── validation.go
        └── validation_test.go
```

## AI Assistant Guidelines

### Educational Approach
When suggesting code:
1. **Explain "Why"** — reasoning behind the solution
2. **Show alternatives** — discuss trade-offs
3. **Connect concepts** — relate to broader patterns
4. **Avoid over-engineering** — prefer simple over clever
5. **Start simple** — ELI5 first, then dive deeper when asked

### Response Style
- Use analogies for complex concepts
- Provide minimal, focused code examples
- Acknowledge what developer did well before suggesting changes
- Prioritize: correctness > readability > performance

---

## Go Idioms & Best Practices

### Naming Conventions
```go
// GOOD: No stuttering
package renderer
func ASCII(text string) {} // Called as: renderer.ASCII()

// BAD: Stuttering
func ASCII(text string) {} // Called as: renderer.ASCII()

// Package: lowercase, single word (parser, renderer, color)
// Exported: PascalCase (BuildCharacterMap, RenderText)
// Unexported: camelCase (parseLines, validateInput)
// Constants: PascalCase or ALL_CAPS for groups
// Test functions: TestFunctionName_Scenario
```

### Error Handling
```go
// GOOD: Wrap with context
if err != nil {
    return fmt.Errorf("failed to load banner %s: %w", filename, err)
}

// BAD: No context
if err != nil {
    return err
}

// BAD: Swallow errors
if err != nil {
    return "" // Lost error information
}

// Error messages: lowercase, no ending punctuation (Go convention ST1005)
```

### Guard Clauses (Early Returns)
```go
// GOOD: Early returns reduce nesting
func Process(input string) error {
    if input == "" {
        return errors.New("empty input")
    }
    if len(input) > MaxLength {
        return errors.New("input too long")
    }
    // Main logic at lowest indentation
    return process(input)
}

// BAD: Deep nesting
func Process(input string) error {
    if input != "" {
        if len(input) <= MaxLength {
            // Main logic deeply nested
        }
    }
}
```

### Interfaces
```go
// GOOD: Accept interfaces, return structs
func ReadData(r io.Reader) (*Data, error)

// GOOD: Small, focused interfaces
type Writer interface {
    Write([]byte) (int, error)
}

// BAD: Return interfaces (limits implementation)
func NewClient() ClientInterface
```

### Best Practices
- Use `strings.Builder` for efficient string concatenation
- Use `bufio.Scanner` for line-by-line file reading
- Constants over magic numbers: define const for all numeric literals

---

## Documentation Standards

### Package Documentation
```go
// Package color provides ANSI color code functionality for terminal output.
//
// Supports multiple color formats:
//   - Named colors: "red", "orange", "blue"
//   - Hex: "#ff0000"
//   - RGB: "rgb(255, 0, 0)"
//
// Example:
//   code, err := color.Parse("rgb(255, 0, 0)")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(color.ANSI(code))
package color
```

### Function Documentation
```go
// Parse validates and converts a color specification to RGB values.
//
// Parameters:
//   - spec: color string (named, hex "#RRGGBB", or "rgb(R,G,B)")
//
// Returns:
//   - [3]int: RGB values (0-255)
//   - error: if format is invalid
//
// Color names are case-insensitive.
func Parse(spec string) ([3]int, error)
```

### Documentation Rules
- First sentence: summary (appears in package lists)
- Start with function/type name
- Exported functions must have Parameters/Returns sections
- Explain non-obvious behavior, not "what" code does
- Use blank lines to separate paragraphs
- Indent code blocks with spaces

### Inline Comments
```go
// GOOD: Explain "why", not "what"
// Use 32 as offset because ASCII printable chars start at 32
offset := 32

// BAD: Restates code
// Set offset to 32
offset := 32

// Only add inline comments for:
// - Non-obvious logic
// - Performance-critical sections
// - Workarounds for bugs
// - Security considerations
```

---

## Testing Standards

### Test-Driven Development (TDD)
1. Write failing test first
2. Write minimal code to pass
3. Refactor if needed
4. Repeat

### Table-Driven Tests
```go
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid RGB",
            input: "rgb(255, 0, 0)",
            want:  "\033[38;2;255;0;0m",
        },
        {
            name:    "invalid RGB range",
            input:   "rgb(256, 0, 0)",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Parse() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Organization
- Test files: `*_test.go` in same package
- Integration tests: `integration_test.go` in main package
- Test fixtures: pre-made files in `testdata/` for testing edge cases

### Coverage Requirements
- Aim for >90% overall coverage
- 100% coverage on parser and renderer packages (critical)
- main() function coverage optional (os.Exit prevents in-process coverage; tested via integration with `exec.Command`)

### Test Priorities
1. **Happy path**: Normal inputs -> expected output
2. **Error cases**: Invalid inputs -> proper errors
3. **Edge cases**: Empty strings, boundaries, special characters
4. **Integration**: End-to-end workflows

---

## Performance & Optimization

### When to Optimize
- **Profile first**: Use `go test -bench`, `pprof`
- **Focus on bottlenecks**: 3% critical code, not 97% non-critical
- **Measure impact**: Benchmark before and after
- **Premature optimization**: Root of all evil — avoid

### Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof

# Clean up profiling artifacts
rm -f cpu.prof mem.prof
```

### Memory Management
```go
// Preallocate slices when size is known
lines := make([]string, 0, expectedCount)

// Use strings.Builder for concatenation
var sb strings.Builder
for _, s := range items {
    sb.WriteString(s)
}
result := sb.String()

// BAD: Avoid repeated string concatenation
result := ""
for _, s := range items {
    result += s // Creates new string each iteration
}
```

### Efficient Patterns
```go
// Pass large structs by pointer
func ProcessData(data *LargeStruct) error

// Pass small structs by value (clarity)
func ValidateConfig(cfg Config) error

// Use buffered I/O
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    process(scanner.Text())
}
```

---

## Architecture Decisions

### Package Responsibilities
- **color**: Parse color formats -> ANSI codes (named, hex, RGB)
- **coloring**: Apply ANSI codes to ASCII art at correct positions
- **flagparser**: Validate CLI argument **structure** only
- **banners**: Expose embedded banner files via `banners.FS` (`embed.FS`)
- **parser**: Load and parse banner files from `fs.FS`
- **renderer**: Convert text to ASCII art
- **validation**: Validate web form input (text length, banner name)
- **handlers**: HTTP handlers, `GenerateASCII`, template cache — used by web server
- **main (cli)**: Orchestrate CLI packages
- **main (web)**: Initialize template cache, register routes, start HTTP server

### Why Separate Packages?
- **Testability**: Each package independently testable
- **Decoupling**: Changes in one don't break others
- **Reusability**: Packages can be used separately
- **Clarity**: Single responsibility per package

### Exit Codes
```go
const (
    ExitSuccess      = 0  // Normal execution
    ExitUsageError   = 1  // Invalid arguments/flags
    ExitBannerError  = 2  // Banner file issue
    ExitRenderError  = 3  // Rendering failed
    ExitColorError   = 4  // Color parsing failed
)
```

---

## Security & Input Validation

### File Operations
```go
// GOOD: Use embedded filesystem with validated paths
func LoadBanner(fsys fs.FS, name string) error {
    allowed := map[string]bool{
        "standard": true, "shadow": true, "thinkertoy": true,
    }
    if !allowed[name] {
        return fmt.Errorf("invalid banner: %s", name)
    }
    path := filepath.Join("testdata", name+".txt")
    data, err := fs.ReadFile(fsys, path)
}

// BAD: User input directly to file path
data, err := os.ReadFile(userInput + ".txt") // Path traversal risk
```

### Input Validation
- Support only ASCII 32-126 (printable)
- Validate early, fail fast
- Return errors, don't silently skip
- Never expose internal paths in errors

---

## Commit Message Format

Use Conventional Commits format:

```
<type>(<scope>): <description>

[optional body]
```

**Types**: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`, `build`, `ci`

**Scopes**: `parser`, `renderer`, `main`, `web`, `handlers`, `banners`, `validation`, `color`, `coloring`, `flagparser`, `templates`, `static`, `docker`, `docs`, `build`, `tests`, `workflows`

**Example**:
```
feat(parser): add validation for banner file format

Added check to ensure banner file has exactly 855 lines before parsing.
Prevents crash on malformed banner files.
```

---

## Build and Release

### Local Build
```bash
make build          # Current platform
make build-all      # All platforms (Linux, macOS, Windows)
```

### CI/CD (GitHub Actions)

Automated checks run on every push and pull request to `main` and `develop`:

- **Test**: Runs `go test ./...` across a matrix of Go 1.21/1.22 on Ubuntu, macOS, and Windows
- **Coverage**: Generates a coverage report and uploads it as a workflow artifact
- **Lint**: Runs `golangci-lint` (v2.1.6) on Ubuntu using `.golangci.yml`
- **Build**: Verifies compilation with `go build ./cmd/ascii-art` on Ubuntu

Workflows are defined in `.github/workflows/`:
- `ci.yml` — test, coverage, lint, and build jobs (triggered on push/PR to main or develop)
- `release.yml` — builds cross-platform binaries and creates a GitHub Release (triggered on `v*` tags)

### Version Management
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Version info injected via Makefile ldflags from git tags
- Tag a release (e.g. `git tag v1.1.0`) to trigger the release workflow
- Update CHANGELOG.md for all releases

### Release Workflow
1. Update CHANGELOG.md
2. Tag version: `git tag v1.x.x`
3. Push tag: `git push origin v1.x.x`
4. GitHub Action builds cross-platform binaries
5. Creates GitHub Release automatically

---

## Common Tasks

### Adding a New Banner Style
1. Add banner file to `internal/banners/<name>.txt` — picked up by `//go:embed *.txt` (web server)
2. Add banner file to `cmd/ascii-art/testdata/<name>.txt` — picked up by `//go:embed testdata/*.txt` (CLI)
3. Update `bannerPaths` map in `cmd/ascii-art/banner.go` to recognize the new name (CLI)
4. Update `ValidateBanner()` in `internal/validation/validation.go` to accept the new name (web)
5. Rebuild binaries (files are embedded at compile time)
6. Add integration tests for both CLI and web
7. Update README.md and CHANGELOG.md

### Adding a Feature
1. Discuss approach (architectural decision)
2. Write tests first (TDD)
3. Implement to pass tests
4. Run `make check`
5. Update documentation (README, inline docs)
6. Update CHANGELOG.md

### Fixing a Bug
1. Write failing test reproducing bug
2. Fix bug
3. Verify test passes
4. Run full test suite: `make test`
5. Run linters: `make lint`
6. Update CHANGELOG.md if user-facing

---

## Code Quality Checklist

### Before Every Commit
```bash
make check    # Runs fmt + vet + lint
make test     # All tests pass
```

### Pre-Pull Request
- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Coverage maintained or improved
- [ ] golangci-lint passes: `golangci-lint run`
- [ ] Documentation updated (if needed)
- [ ] CHANGELOG.md updated (if user-facing change)
- [ ] Conventional commit message used
- [ ] CI passes (test, lint, build) on push/PR

### Code Review Focus
- [ ] No external dependencies added
- [ ] Error handling with proper wrapping
- [ ] Functions <50 lines (suggest refactoring if longer)
- [ ] No deep nesting (>3 levels -> use guard clauses)
- [ ] No magic numbers (use constants)
- [ ] No code duplication
- [ ] Follows naming conventions (no stuttering)

---

## DO NOT

- Add external dependencies (use only Go standard library)
- Modify banner files in `testdata/`
- Skip tests or reduce coverage
- Commit without running `make check`
- Use deprecated Go features
- Use reflection (use interfaces instead)
- Leave TODO/FIXME without GitHub issues
- Commit debug print statements
- Use panic (return errors instead)

---

**Final Reminders**:
- Correctness > Readability > Performance (in that order)
- Simple code beats clever code
- Document the "why", not the "what"
- Profile before optimizing
- Test-driven development always
- Zero external dependencies
- Run `make check` before every commit

*This file follows the [AGENTS.md](https://agents.md/) open standard for guiding AI coding agents.*
