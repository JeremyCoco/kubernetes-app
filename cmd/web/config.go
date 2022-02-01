package main

import (
	"errors"
	"flag"
	"strings"
)

func getConfigFlags() (*configFlags, error) {
	addr := flag.String("addr", ":8181", "Set server netowrk address")
	dsn := flag.String("dsn", "root:root@tcp(db:3306)/todo", "Set database data source name")
	flag.Parse()

	if strings.TrimSpace(*addr) == "" || strings.TrimSpace(*dsn) == "" {
		return nil, errors.New("all command-line flags are required")
	}

	return &configFlags{
		addr: *addr,
		dsn:  *dsn,
	}, nil
}
