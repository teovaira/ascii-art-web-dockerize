# Program Flowchart

Execution flows for both the **CLI** and **Web** interfaces.

## CLI Flow

The CLI has two modes: **Normal** (text only) and **Color** (with ANSI coloring).

```mermaid
flowchart TD
    A["CLI Arguments<br/>os.Args"] --> B{"hasColorFlag?<br/>--color= prefix"}

    B -->|No| C["ParseArgs()<br/>text, banner"]
    B -->|Yes| D["flagparser.ParseArgs()<br/>validate syntax"]

    C --> E["GetBannerPath()<br/>banner file path"]
    D --> F["extractColorArgs()<br/>colorSpec, substring,<br/>text, banner"]

    E --> E2["GetBannerFS()<br/>embedded filesystem"]
    E2 --> G["parser.LoadBanner(fsys, path)<br/>Banner map"]
    F --> H["color.Parse()<br/>RGB struct"]

    H --> I["GetBannerPath()<br/>banner file path"]
    I --> I2["GetBannerFS()<br/>embedded filesystem"]
    I2 --> J["parser.LoadBanner(fsys, path)<br/>Banner map"]
    J --> K["color.ANSI()<br/>ANSI escape code"]

    G --> L["renderer.ASCII()<br/>ASCII art string"]
    L --> M["fmt.Print()<br/>stdout"]

    K --> N["For each line in text"]
    N --> O["renderer.ASCII()<br/>ASCII art lines"]
    O --> P["parser.CharWidths()<br/>character widths"]
    P --> Q["coloring.ApplyColor()<br/>colored ASCII art"]
    Q --> R{"More lines?"}
    R -->|Yes| N
    R -->|No| S["fmt.Print()<br/>stdout"]

    style B fill:#f39c12,color:#fff
    style R fill:#f39c12,color:#fff
    style M fill:#2ecc71,color:#fff
    style S fill:#2ecc71,color:#fff
```

## Web Flow

Browser request through the HTTP server to rendered response.

```mermaid
flowchart TD
    W["Browser<br/>GET /"] --> H1["app.Home()<br/>render index.html"]
    H1 --> R1["200 OK<br/>HTML page with form"]

    W2["Browser<br/>POST /ascii-art<br/>text + banner"] --> H2["app.HandleASCIIArt()"]

    H2 --> V1["validation.ValidateText()<br/>check empty / too long"]
    V1 -->|invalid| E1["re-render page<br/>with error message<br/>400 Bad Request"]

    V1 -->|valid| V2["validation.ValidateBanner()<br/>check known banner"]
    V2 -->|invalid| E2["re-render page<br/>with error message<br/>404 Not Found"]

    V2 -->|valid| P["parser.LoadBanner(banners.FS, banner)<br/>Banner map"]
    P --> RND["renderer.ASCII(text, banner)<br/>ASCII art string"]
    RND -->|error| E3["re-render page<br/>with error message<br/>500 Internal Server Error"]
    RND -->|success| R2["re-render page<br/>with ASCII art in pre<br/>200 OK"]

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
