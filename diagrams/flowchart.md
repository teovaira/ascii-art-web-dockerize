# Program Flowchart

Execution flows for both the **CLI** and **Web** interfaces.

## CLI Flow

The CLI has two modes: **Normal** (text only) and **Color** (with ANSI coloring).

```mermaid
flowchart TD
    A["CLI Arguments\nos.Args"] --> B{"hasColorFlag?\n--color= prefix"}

    B -->|No| C["ParseArgs()\ntext, banner"]
    B -->|Yes| D["flagparser.ParseArgs()\nvalidate syntax"]

    C --> E["GetBannerPath()\nbanner file path"]
    D --> F["extractColorArgs()\ncolorSpec, substring,\ntext, banner"]

    E --> E2["GetBannerFS()\nembedded filesystem"]
    E2 --> G["parser.LoadBanner(fsys, path)\nBanner map"]
    F --> H["color.Parse()\nRGB struct"]

    H --> I["GetBannerPath()\nbanner file path"]
    I --> I2["GetBannerFS()\nembedded filesystem"]
    I2 --> J["parser.LoadBanner(fsys, path)\nBanner map"]
    J --> K["color.ANSI()\nANSI escape code"]

    G --> L["renderer.ASCII()\nASCII art string"]
    L --> M["fmt.Print()\nstdout"]

    K --> N["For each line in text"]
    N --> O["renderer.ASCII()\nASCII art lines"]
    O --> P["parser.CharWidths()\ncharacter widths"]
    P --> Q["coloring.ApplyColor()\ncolored ASCII art"]
    Q --> R{"More lines?"}
    R -->|Yes| N
    R -->|No| S["fmt.Print()\nstdout"]

    style B fill:#f39c12,color:#fff
    style R fill:#f39c12,color:#fff
    style M fill:#2ecc71,color:#fff
    style S fill:#2ecc71,color:#fff
```

## Web Flow

Browser request through the HTTP server to rendered response.

```mermaid
flowchart TD
    W["Browser\nGET /"] --> H1["app.Home()\nrender index.html"]
    H1 --> R1["200 OK\nHTML page with form"]

    W2["Browser\nPOST /ascii-art\ntext + banner"] --> H2["app.HandleASCIIArt()"]

    H2 --> V1["validation.ValidateText()\ncheck empty / too long"]
    V1 -->|invalid| E1["re-render page\nwith error message\n400 Bad Request"]

    V1 -->|valid| V2["validation.ValidateBanner()\ncheck known banner"]
    V2 -->|invalid| E2["re-render page\nwith error message\n404 Not Found"]

    V2 -->|valid| P["parser.LoadBanner(banners.FS, banner)\nBanner map"]
    P --> RND["renderer.ASCII(text, banner)\nASCII art string"]
    RND -->|error| E3["re-render page\nwith error message\n500 Internal Server Error"]
    RND -->|success| R2["re-render page\nwith ASCII art in pre\n200 OK"]

    style E1 fill:#e74c3c,color:#fff
    style E2 fill:#e74c3c,color:#fff
    style E3 fill:#e74c3c,color:#fff
    style R1 fill:#2ecc71,color:#fff
    style R2 fill:#2ecc71,color:#fff
```

## Mode Comparison

| Aspect | CLI Normal | CLI Color | Web |
|--------|-----------|-----------|-----|
| Entry | `os.Args` | `os.Args` + `--color` | HTTP POST form |
| Validation | `ParseArgs()` | `flagparser.ParseArgs()` | `validation` package |
| Banner FS | `GetBannerFS()` (cmd embed) | `GetBannerFS()` (cmd embed) | `banners.FS` (internal embed) |
| Color | — | `color.Parse()` + `coloring.ApplyColor()` | — |
| Output | `fmt.Print()` stdout | `fmt.Print()` stdout | HTML template `<pre>` |
| Error handling | `os.Exit(code)` | `os.Exit(code)` | HTTP status + inline message |
