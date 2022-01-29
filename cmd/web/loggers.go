package main

import (
	"log"
	"os"
)

var (
	infoLogger = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	errLogger  = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
)
