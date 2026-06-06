// Command boxctl is the box builder tool: Define + Build + Verify + Resolve over
// the content-addressed Box primitive (pkg/box).
//
// It deliberately does NOT sign. Signing is a human act (realness is rooted in a
// human key holder); boxctl builds and verifies, and can Attach a human-produced
// signature, but never holds a private key.
//
//	boxctl build  -t <mediaType> [-s <schema>] -f <file> [-o out.box] [-ref <digest> ...]
//	boxctl digest <file.box>
//	boxctl verify <sealed.box>
//	boxctl resolve -dir <dir> <root-digest>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/agennext/agent-chat/pkg/box"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "build":
		err = build(os.Args[2:])
	case "digest":
		err = digest(os.Args[2:])
	case "verify":
		err = verify(os.Args[2:])
	case "resolve":
		err = resolve(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "boxctl:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: boxctl <build|digest|verify|resolve> ... (signing is human; boxctl does not sign)")
}

// multiFlag collects repeatable string flags (e.g. -ref).
type multiFlag []string

func (m *multiFlag) String() string { return fmt.Sprint([]string(*m)) }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// build: Define the form, Build with file content, emit the .box and its digest.
func build(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	mt := fs.String("t", "", "media type (form)")
	sc := fs.String("s", "", "schema reference (contract)")
	f := fs.String("f", "", "content file (the ground)")
	out := fs.String("o", "", "output .box file (default stdout)")
	var refs multiFlag
	fs.Var(&refs, "ref", "child box digest (repeatable)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *mt == "" || *f == "" {
		return fmt.Errorf("-t and -f are required")
	}
	content, err := os.ReadFile(*f)
	if err != nil {
		return err
	}
	b := box.Define(*mt, *sc).Build(content, refs...)
	enc, err := b.Encode()
	if err != nil {
		return err
	}
	d, err := b.Digest()
	if err != nil {
		return err
	}
	if *out != "" {
		if err := os.WriteFile(*out, enc, 0o644); err != nil {
			return err
		}
	} else {
		_, _ = os.Stdout.Write(enc)
		fmt.Println()
	}
	fmt.Fprintln(os.Stderr, d) // digest to stderr so stdout stays the box
	return nil
}

func digest(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: boxctl digest <file.box>")
	}
	var b box.Box
	if err := readJSON(args[0], &b); err != nil {
		return err
	}
	d, err := b.Digest()
	if err != nil {
		return err
	}
	fmt.Println(d)
	return nil
}

func verify(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: boxctl verify <sealed.box>")
	}
	var s box.Sealed
	if err := readJSON(args[0], &s); err != nil {
		return err
	}
	ok, err := s.Verify()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("INVALID: seal does not verify")
	}
	fmt.Println("OK:", s.Digest)
	return nil
}

// resolve: load every *.box in dir into a Graph and resolve the DAG at root.
func resolve(args []string) error {
	fs := flag.NewFlagSet("resolve", flag.ExitOnError)
	dir := fs.String("dir", ".", "directory of .box files")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: boxctl resolve -dir <dir> <root-digest>")
	}
	g := box.NewGraph()
	entries, err := filepath.Glob(filepath.Join(*dir, "*.box"))
	if err != nil {
		return err
	}
	for _, p := range entries {
		var b box.Box
		if err := readJSON(p, &b); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		if _, err := g.Add(b); err != nil {
			return err
		}
	}
	order, err := g.Resolve(fs.Arg(0))
	if err != nil {
		return err
	}
	for _, d := range order {
		fmt.Println(d)
	}
	return nil
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
