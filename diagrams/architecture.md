# Architecture Overview

High-level view of the system architecture. Packages are grouped by responsibility layer.

```mermaid
flowchart LR
    subgraph CLILayer["CLI Layer"]
        main["main\n(cmd/ascii-art)"]
    end

    subgraph WebLayer["Web Layer"]
        webmain["main\n(cmd/ascii-art-web)"]
        handlers["handlers\nHTTP handlers"]
    end

    subgraph Shared["Shared Infrastructure"]
        banners["banners\nEmbedded .txt files"]
        validation["validation\nInput validation"]
    end

    subgraph Input["Input Processing"]
        flagparser["flagparser\nCLI validation"]
        color["color\nColor parsing"]
    end

    subgraph Core["Core Engine"]
        parser["parser\nBanner loading"]
        renderer["renderer\nASCII rendering"]
    end

    subgraph Output["Output Processing"]
        coloring["coloring\nANSI color application"]
    end

    main -->|"validates args"| flagparser
    main -->|"parses color spec"| color
    main -->|"loads banner"| parser
    main -->|"renders text"| renderer
    main -->|"applies color"| coloring

    webmain -->|"initializes"| handlers
    handlers -->|"validates input"| validation
    handlers -->|"loads banner"| parser
    handlers -->|"renders text"| renderer
    handlers -->|"reads embedded FS"| banners
    parser -->|"reads embedded FS"| banners

    style CLILayer fill:#4a90d9,color:#fff
    style WebLayer fill:#9b59b6,color:#fff
    style Shared fill:#e74c3c,color:#fff
    style Input fill:#7b68ee,color:#fff
    style Core fill:#2ecc71,color:#fff
    style Output fill:#e67e22,color:#fff
```

## Package Responsibilities

| Layer | Package | Responsibility |
|-------|---------|---------------|
| CLI | `main` (cmd/ascii-art) | Orchestrates CLI packages, handles I/O and exit codes |
| Web | `main` (cmd/ascii-art-web) | Initializes template cache, registers routes, starts HTTP server |
| Web | `handlers` | HTTP handlers, ASCII generation, template rendering |
| Shared | `banners` | Banner `.txt` files embedded into binary at compile time |
| Shared | `validation` | Validates web form input (text length, banner name) |
| Input | `flagparser` | Validates CLI argument structure |
| Input | `color` | Parses color specs (named, hex, RGB) into RGB values |
| Core | `parser` | Reads banner files from `fs.FS`, builds character maps |
| Core | `renderer` | Converts text to ASCII art using banner maps |
| Output | `coloring` | Applies ANSI color codes to rendered ASCII art |

## Key Design Decisions

- **Shared core engine** — both CLI and web server use the same `parser` and `renderer` packages
- **Embedded banners** — `internal/banners` embeds `.txt` files at compile time; both binaries are self-contained
- **No import cycles** — `handlers` imports `parser`, `renderer`, `validation`, `banners`; nothing imports back
- **Stateless packages** — all functions are pure transformations with no global state
- **Web input validation** — `validation` package is web-only; CLI uses `flagparser` and `bannerPaths` map
