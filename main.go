package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"

  "github.com/russross/blackfriday"
)

func main() {
  // Makrdown suffix suffixes
  mdSuffixes := []string{".markdown", ".mdown", ".mkdn", ".md", ".mkd", ".mdwn", ".mdtxt", ".mdtext", ".text"}
  extensions := 0
  extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
  extensions |= blackfriday.EXTENSION_TABLES
  extensions |= blackfriday.EXTENSION_FENCED_CODE
  extensions |= blackfriday.EXTENSION_AUTOLINK
  extensions |= blackfriday.EXTENSION_STRIKETHROUGH
  extensions |= blackfriday.EXTENSION_SPACE_HEADERS

  var port int
  flag.IntVar(&port, "port", 7070, "Port number for http server")
  flag.Parse()
  var defaultFile = flag.Arg(0)
  if defaultFile == "" {
    defaultFile = "README.md"
  }
  htmlFlags := 0
  htmlFlags |= blackfriday.HTML_USE_XHTML
  htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
  htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
  htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
  htmlFlags |= blackfriday.HTML_TOC

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    urlPath := r.URL.Path[1:]
    var path string
    _, err := os.Stat(urlPath)
    if os.IsNotExist(err) {
      _, errDefault := os.Stat(defaultFile)
      if os.IsNotExist(errDefault) {
        _, errIndex := os.Stat("index.html")
        if os.IsNotExist(errIndex) {
          http.Error(w, fmt.Sprintf("Failed to find file to read tried README.md and index.html .... %s", path), http.StatusInternalServerError)
          return
        } else {
          path = "index.html"
        }
      } else {
        path = defaultFile
      }
    } else {
      path = urlPath
    }

    md := false
    for _, s := range mdSuffixes {
      if strings.HasSuffix(path, s) {
        md = true
        break
      }
    }
    if md {
      input, err := os.Open(path) // For read access.
      if err != nil {
        log.Fatal(err)
      }
      b, err := ioutil.ReadAll(input)
      if err != nil {
        http.Error(w, fmt.Sprintf("Failed to read file %s", path), http.StatusInternalServerError)
        return
      }
      w.Header().Set("Content-Type", "text/html; charset=utf-8")
      w.WriteHeader(http.StatusOK)
      renderer := blackfriday.HtmlRenderer(htmlFlags, "mdserve Markdown http serve", "")
      output := blackfriday.Markdown(b, renderer, extensions)
      w.Write(output)
      return
    } else if strings.HasSuffix(r.URL.Path, ".json") {
      w.Header().Set("Content-Type", "application/json; charset=utf-8")
    } else if strings.HasSuffix(r.URL.Path, ".html") {
      w.Header().Set("Content-Type", "text/html; charset=utf-8")
    }
    http.ServeFile(w, r, path)
  })

  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
