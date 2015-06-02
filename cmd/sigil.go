package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gliderlabs/sigil"
	_ "github.com/gliderlabs/sigil/builtin"
)

var Version string

var (
	filename = flag.String("f", "", "use template file instead of STDIN")
	posix    = flag.Bool("p", false, "preprocess with POSIX variable expansion")
	version  = flag.Bool("v", false, "prints version")
)

func template() (string, string, error) {
	if *filename != "" {
		data, err := ioutil.ReadFile(*filename)
		if err != nil {
			return "", "", err
		}
		sigil.PushPath(filepath.Dir(*filename))
		return string(data), filepath.Base(*filename), nil
	}
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", "", err
	}
	return string(data), "<stdin>", nil
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *posix {
		sigil.PosixPreprocess = true
	}
	if os.Getenv("SIGIL_PATH") != "" {
		sigil.TemplatePath = strings.Split(os.Getenv("SIGIL_PATH"), ":")
	}
	vars := make(map[string]string)
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	tmpl, name, err := template()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	render, err := sigil.Execute(tmpl, vars, name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(render)
}
