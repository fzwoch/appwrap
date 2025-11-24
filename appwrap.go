package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

var version = "0.2.0"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "AppWrap - Wrap binaries into minimal macOS application bundles. (%s)\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright (C) 2025 Florian Zwoch <fzwoch@gmail.com>\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <path/to/some/binary/>\n", os.Args[0])
	}
	flag.Parse()

	if len(os.Args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	bin, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exe := filepath.Base(bin)
	app := exe + ".app"
	_, err = os.Stat(bin)

	if os.IsNotExist(err) {
		fmt.Println(bin, "does not exist")
		os.Exit(1)
	}

	_, err = os.Stat(app)
	if !os.IsNotExist(err) {
		fmt.Println(app, "already exists")
		os.Exit(1)
	}

	err = os.MkdirAll(path.Join(app, "Contents", "MacOS"), 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r, err := os.Open(bin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	file := path.Join(app, "Contents", "MacOS", exe)

	w, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Chmod(file, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := struct {
		Name string
	}{
		Name: exe,
	}

	tpl, err := template.New("plist").Parse(tpl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	plist, err := os.Create(path.Join(app, "Contents", "Info.plist"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer plist.Close()

	tpl.Execute(plist, data)
}

const tpl = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" 
"http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>CFBundleExecutable</key>
        <string>{{ .Name }}</string>
        <key>CFBundleIdentifier</key>
        <string>com.example.{{ .Name }}</string>
        <key>NSPrincipalClass</key>
        <string>NSApplication</string>
    </dict>
</plist>
`
