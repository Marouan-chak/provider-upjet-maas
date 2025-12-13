package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path-to-provider-metadata.yaml>\n", os.Args[0])
		os.Exit(2)
	}

	path := os.Args[1]
	// #nosec G304 -- path is a fixed repo file passed by go:generate / Makefile
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
		os.Exit(1)
	}

	in := string(b)
	out := in

	// Replace PEM blocks embedded in JSON example strings.
	// These are placeholders from upstream docs but can be flagged by secret scanners.
	reCert := regexp.MustCompile(`("certificate"\s*:\s*)"[^\"]*?-----BEGIN CERTIFICATE-----\\n[^\"]*?-----END CERTIFICATE-----\\n"`)
	out = reCert.ReplaceAllString(out, `${1}"<certificate PEM data>"`)

	reKey := regexp.MustCompile(`("key"\s*:\s*)"[^\"]*?-----BEGIN [A-Z ]*PRIVATE KEY-----\\n[^\"]*?-----END [A-Z ]*PRIVATE KEY-----\\n"`)
	out = reKey.ReplaceAllString(out, `${1}"<private key PEM data>"`)

	if out == in {
		return
	}

	// #nosec G306 -- generated documentation file should be world-readable in the repo
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
		os.Exit(1)
	}
}
