package main

import (
	"log"
	"os"
)

var (
	Wrn *log.Logger
	Inf *log.Logger
	Err *log.Logger
)

func init() {
	Inf = log.New(os.Stdout, "\033[36m[INF]\033[0m ", log.Lmsgprefix|log.Ldate|log.Ltime)
	Wrn = log.New(os.Stdout, "\033[33m[WRN]\033[0m ", log.Lmsgprefix|log.Ldate|log.Ltime)
	Err = log.New(os.Stdout, "\033[31m[ERR]\033[0m ", log.Lmsgprefix|log.Ldate|log.Ltime)
}
