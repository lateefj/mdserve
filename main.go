package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"

  "github.com/russross/blackfriday"
)

func main() {
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
  var file = flag.Arg(0)
  if file == "" {
    file = "README.md"
  }
  htmlFlags := 0
  htmlFlags |= blackfriday.HTML_USE_XHTML
  htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
  htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
  htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
  htmlFlags |= blackfriday.HTML_TOC

  renderer := blackfriday.HtmlRenderer(htmlFlags, "mdserve Markdown http serve", "")
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    input, err := os.Open(file) // For read access.
    if err != nil {
      log.Fatal(err)
    }
    b, err := ioutil.ReadAll(input)
    if err != nil {
      http.Error(w, fmt.Sprintf("Failed to read file %s", file), http.StatusInternalServerError)
      return
    }
    output := blackfriday.Markdown(b, renderer, extensions)
    w.WriteHeader(http.StatusOK)

    w.Write(output)
  })

  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
