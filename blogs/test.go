package main

import "time"

func main() {
  t := time.Now()
  println(t.IsDST())
  name, off := t.Zone()
  println(name, off / 3600)
  loc := t.Location()
  println(loc.String())
}
