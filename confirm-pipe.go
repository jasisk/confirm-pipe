package main

import (
  "bytes"
  "log"
  "fmt"
  "io"
  "os"

  "github.com/pkg/term"
)

func main() {
  var buffer bytes.Buffer
  rc := 0

  defer func () { os.Exit(rc) }()

  fd, err := os.OpenFile("/dev/tty", os.O_RDWR, 0644)
  defer fd.Close()

  w := io.MultiWriter(&buffer, fd)
  io.Copy(w, os.Stdin)

  if err != nil {
    log.Fatal(err)
  }

  fmt.Fprintf(fd, "\ntype c to confirm, anything else to quit\033[?25l")

  t, _ := term.Open("/dev/tty")

  defer fmt.Fprintf(fd, "\n\033[?25h")
  defer t.Close()
  defer t.Restore()

  term.RawMode(t)
  bytes := make([]byte, 1)
  num, _ := t.Read(bytes)

  if num == 1 {
    if bytes[0] == 99 {
      fmt.Print(buffer.String())
      return
    }
  }

  rc = 1
}
