# Sequence Diagram

## CLI — Color Mode

Call sequence for **color mode** — the more complex CLI execution path. Shows how `main` orchestrates all packages over time.

```mermaid
sequenceDiagram
    actor User
    participant main
    participant flagparser
    participant color
    participant parser
    participant renderer
    participant coloring

    User->>main: os.Args with --color flag

    Note over main: hasColorFlag() = true

    main->>flagparser: ParseArgs(args)
    flagparser-->>main: nil (valid)

    Note over main: extractColorArgs()

    main->>color: Parse(colorSpec)
    color-->>main: RGB{R, G, B}

    main->>main: GetBannerPath(banner)
    main->>main: GetBannerFS()

    main->>parser: LoadBanner(fsys, path)
    parser-->>main: Banner map[rune][]string

    main->>color: ANSI(rgb)
    color-->>main: ANSI escape code string

    loop For each line in text
        main->>renderer: ASCII(line, banner)
        renderer-->>main: string (ASCII art with newlines)

        main->>main: split rendered ASCII into artLines

        main->>parser: CharWidths(line, banner)
        parser-->>main: []int (character widths)

        main->>coloring: ApplyColor(artLines, line, substring, colorCode, widths)

        Note over coloring: findPositions(line, substring) + colorLine() for each art line

        coloring-->>main: []string (colored lines)
    end

    main->>User: Colored ASCII art to stdout
```

## CLI — Normal Mode

For comparison, normal mode has a much shorter sequence:

```mermaid
sequenceDiagram
    actor User
    participant main
    participant parser
    participant renderer

    User->>main: os.Args without --color

    Note over main: ParseArgs() + GetBannerPath() + GetBannerFS()

    main->>parser: LoadBanner(fsys, path)
    parser-->>main: Banner map[rune][]string

    main->>renderer: ASCII(text, banner)
    renderer-->>main: ASCII art string

    main->>User: Plain ASCII art to stdout
```

## Web — HTTP Request/Response

Call sequence for a browser form submission through the web server.

```mermaid
sequenceDiagram
    actor Browser
    participant main as main (web)
    participant handlers
    participant validation
    participant parser
    participant renderer

    Note over main: startup — NewTemplateCache() + register routes

    Browser->>main: GET /
    main->>handlers: app.Home(w, r)
    handlers-->>Browser: 200 OK — index.html with form

    Browser->>main: POST /ascii-art (text, banner)
    main->>handlers: app.HandleASCIIArt(w, r)

    handlers->>handlers: r.ParseForm()
    handlers->>handlers: GenerateASCII(text, banner)

    handlers->>validation: ValidateText(text)
    alt invalid text
        validation-->>handlers: ErrEmptyText / ErrTextTooLong
        handlers-->>Browser: 400 Bad Request — page with error message
    end

    handlers->>validation: ValidateBanner(banner)
    alt invalid banner
        validation-->>handlers: ErrInvalidBanner
        handlers-->>Browser: 404 Not Found — page with error message
    end

    handlers->>parser: LoadBanner(banners.FS, banner+".txt")
    alt banner file missing
        parser-->>handlers: error
        handlers-->>Browser: 404 Not Found — page with error message
    end
    parser-->>handlers: Banner map[rune][]string

    handlers->>renderer: ASCII(text, banner)
    alt render error
        renderer-->>handlers: error
        handlers-->>Browser: 500 Internal Server Error — page with error message
    end
    renderer-->>handlers: ASCII art string

    handlers-->>Browser: 200 OK — page with ASCII art in pre block
```
