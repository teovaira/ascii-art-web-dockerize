# Class Diagram

Package relationships, exported types, and function signatures.

```mermaid
classDiagram
    class climain["main (cli)"] {
        +main()
        +ParseArgs(args []string) (string, string, error)
        +GetBannerPath(banner string) (string, error)
        +GetBannerFS() fs.FS
        -runColorMode(args []string)
        -hasColorFlag(args []string) bool
        -extractColorArgs(args []string) (string, string, string, string, error)
    }

    class webmain["main (web)"] {
        +main()
    }

    class handlers {
        <<package>>
        +NewTemplateCache() (map[string]*template.Template, error)
    }

    class Application {
        <<struct>>
        +TemplateCache map[string]*template.Template
        +Home(w, r)
        +HandleASCIIArt(w, r)
    }

    class PageData {
        <<struct>>
        +Result string
        +Title string
        +Error string
    }

    class GenerateASCII {
        <<function>>
        +GenerateASCII(text, banner string) (string, int, error)
    }

    class validation {
        <<package>>
        +ValidateText(text string) error
        +ValidateBanner(banner string) error
        +MaxTextLength int
    }

    class banners {
        <<package>>
        +FS embed.FS
    }

    class parser {
        <<package>>
        +LoadBanner(fsys fs.FS, path string) (Banner, error)
        +CharWidths(text string, banner Banner) []int
    }

    class Banner {
        <<type alias>>
        map~rune, []string~
    }

    class renderer {
        <<package>>
        +ASCII(input string, banner map~rune, []string~) (string, error)
    }

    class color {
        <<package>>
        +Parse(colorSpec string) (RGB, error)
        +ANSI(rgb RGB) string
    }

    class RGB {
        <<struct>>
        +R uint8
        +G uint8
        +B uint8
    }

    class coloring {
        <<package>>
        +ApplyColor(asciiArt []string, text string, substring string, colorCode string, charWidths []int) []string
        +Reset string
    }

    class flagparser {
        <<package>>
        +ParseArgs(args []string) error
    }

    webmain --> handlers : initializes
    handlers --> Application : creates
    Application --> PageData : renders with
    Application --> GenerateASCII : calls
    GenerateASCII --> validation : validates input
    GenerateASCII --> parser : loads banner
    GenerateASCII --> renderer : renders text
    parser --> banners : reads FS
    climain --> parser : loads banners
    climain --> renderer : renders text
    climain --> color : parses colors
    climain --> coloring : applies colors
    climain --> flagparser : validates args
    parser --> Banner : returns
    color --> RGB : returns
```

## Dependency Rules

- Neither `parser`, `renderer`, `validation`, `coloring`, `color`, nor `flagparser` imports any other internal package
- `handlers` imports `parser`, `renderer`, `validation`, and `banners` — no cycles
- Both `main` packages are the only entry points that wire everything together
- All packages depend only on the Go standard library
