package main

import (
	"os/exec"
	"strings"
)

func main() {
  cmd := exec.Command("pandoc", "-w", "html", "-f", "latex")
  cmd.Stdin = strings.NewReader(`
\documentclass[12pt, letterpaper]{article}
\begin{document}
\begin{math}
E=mc^2
\end{math} is typeset in a paragraph using inline math mode---as is $E=mc^2$, and so too is \(E=mc^2\).
\end{document}
  `)
  out, err := cmd.CombinedOutput()
  println(string(out))
  if err != nil {
    panic(err)
  }
}
