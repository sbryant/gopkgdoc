// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// +build ignore

// Command print fetches and prints package documentation.
//
// Usage: go run print.go importPath
package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/gopkgdoc/doc"
	"log"
	"net/http"
	"strings"
)

func indent(s string, n int) string {
	const space = "                       "
	return strings.Replace(strings.TrimSpace(s), "\n", "\n"+space[:n], -1)
}

var (
	etag    = flag.String("etag", "", "Etag")
	local   = flag.Bool("local", false, "Get package from local directory.")
	present = flag.Bool("present", false, "Get presentation.")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Usage: go run print.go importPath")
	}
	if *present {
		printPresentation(flag.Args()[0])
	} else {
		printPackage(flag.Args()[0])
	}
}

func printPackage(path string) {
	var (
		pdoc *doc.Package
		err  error
	)
	if *local {
		pdoc, err = doc.GetDir(path)
	} else {
		pdoc, err = doc.Get(http.DefaultClient, path, *etag)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ImportPath:  ", pdoc.ImportPath)
	fmt.Println("ProjectRoot: ", pdoc.ProjectRoot)
	fmt.Println("ProjectName: ", pdoc.ProjectName)
	fmt.Println("ProjectURL:  ", pdoc.ProjectURL)
	fmt.Println("BrowseURL :  ", pdoc.BrowseURL)
	fmt.Println("Updated:     ", pdoc.Updated)
	fmt.Println("Etag:        ", pdoc.Etag)
	fmt.Println("Name:        ", pdoc.Name)
	fmt.Println("IsCmd:       ", pdoc.IsCmd)
	fmt.Println("Synopsis:    ", pdoc.Synopsis)
	fmt.Println("Doc:         ", indent(pdoc.Doc, 14))
	fmt.Println("Errors:")
	fmt.Println("Examples:", len(pdoc.Examples))
	for _, s := range pdoc.Errors {
		fmt.Println("    ", s)
	}
	fmt.Println("Files:")
	for _, f := range pdoc.Files {
		fmt.Println("    ", f)
	}
	fmt.Println("Imports:")
	for _, i := range pdoc.Imports {
		fmt.Println("    ", i)
	}
	fmt.Println("TestImports:")
	for _, i := range pdoc.TestImports {
		fmt.Println("    ", i)
	}
	fmt.Println("XTestImports:")
	for _, i := range pdoc.XTestImports {
		fmt.Println("    ", i)
	}
	fmt.Println("References:")
	for _, i := range pdoc.References {
		fmt.Println("    ", i)
	}
	for _, c := range pdoc.Consts {
		fmt.Println("Const:")
		fmt.Println("    Decl:  ", indent(c.Decl.Text, 12))
		fmt.Println("    Doc:   ", indent(c.Doc, 12))
		fmt.Println("    URL:   ", c.URL)
	}
	for _, c := range pdoc.Vars {
		fmt.Println("Var:")
		fmt.Println("    Decl:  ", indent(c.Decl.Text, 12))
		fmt.Println("    Doc:   ", indent(c.Doc, 12))
		fmt.Println("    URL:   ", c.URL)
	}
	for _, f := range pdoc.Funcs {
		fmt.Println("Func:")
		fmt.Println("    Decl:  ", indent(f.Decl.Text, 12))
		fmt.Println("    Doc:   ", indent(f.Doc, 12))
		fmt.Println("    URL:   ", f.URL)
		fmt.Println("    Examples:", len(f.Examples))
	}
	for _, t := range pdoc.Types {
		fmt.Println("Type:")
		fmt.Println("    Decl:  ", indent(t.Decl.Text, 12))
		fmt.Println("    Doc:   ", indent(t.Doc, 12))
		fmt.Println("    URL:   ", t.URL)
		fmt.Println("    Examples:", len(t.Examples))
		for _, f := range t.Funcs {
			fmt.Println("    Func:")
			fmt.Println("        Decl:  ", indent(f.Decl.Text, 16))
			fmt.Println("        Doc:   ", indent(f.Doc, 16))
			fmt.Println("        URL:   ", f.URL)
			fmt.Println("        Examples:", len(f.Examples))
		}
		for _, m := range t.Methods {
			fmt.Println("    Method:")
			fmt.Println("        Decl:  ", indent(m.Decl.Text, 16))
			fmt.Println("        Doc:   ", indent(m.Doc, 16))
			fmt.Println("        URL:   ", m.URL)
			fmt.Println("        Examples:", len(m.Examples))
		}
	}
}

func printPresentation(path string) {
	pres, err := doc.GetPresentation(http.DefaultClient, path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", pres.Files[pres.Filename])
	for name, data := range pres.Files {
		if name != pres.Filename {
			fmt.Printf("---------- %s ----------\n%s\n", name, data)
		}
	}
}
