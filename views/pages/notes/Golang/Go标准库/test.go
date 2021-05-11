package main

import (
	"flag"
	"fmt"
)

func main() {
	a := flag.Int("a", 123, "a是光")
	b := flag.Int("b", 456, "b是电")
	c := flag.Int("c", 789, "c是唯一的神话")
	d := flag.Int("d", 101112, "想不出来了")
	flag.Parse()
	fmt.Println(*a, *b, *c, *d)
}
