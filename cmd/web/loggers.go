package main

import (
	"flag"
	"go/build"
	"log"
	"os"
)

var (
	infoLogger *log.Logger
	errLogger  *log.Logger
)

func init() {
	// set location of log file
	var logpath = build.Default.GOPATH + "/info.log"

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}

	infoLogger = log.New(file, "[INFO] ", log.LstdFlags)
	errLogger = log.New(file, "[ERROR] ", log.LstdFlags)
}
