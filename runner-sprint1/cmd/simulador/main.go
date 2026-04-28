package main

import "fmt"

// version deve permanecer como variável para permitir injeção por -ldflags.
var version = "dev"

func main() {
	fmt.Printf("simulador v%s — em construção\n", version)
}
