package main

import (
	repl "warehouse/pkg/REPL"
)

func main() {
	repl := &repl.REPL{}
	repl.Exec()
}
