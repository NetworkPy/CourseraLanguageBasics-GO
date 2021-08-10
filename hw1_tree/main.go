package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	out := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {
	prefix := ""
	dirTreeRec(out, path, prefix, printFiles)
	return nil
}

func dirTreeRec(out *bytes.Buffer, paths, prefix string, printFiles bool) {
	files, err := ioutil.ReadDir(paths)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	if !printFiles {
		for f := range files {
			if files[f].IsDir() {
				count++
			}
		}
	} else {
		count++
	}

	for idx, file := range files {
		if !printFiles {
			if count == 0 {
				break
			}
			if !file.IsDir() {
				continue
			} else {
				count--
			}
		}
		if idx == len(files)-1 || count == 0 {
			if file.IsDir() {
				// fmt.Print(prefix+"└───", file.Name(), "\n")
				fmt.Fprint(out, prefix+"└───", file.Name(), "\n")
				dirTreeRec(out, path.Join(paths, file.Name()), prefix+"\t", printFiles)
				continue
			} else {
				if file.Size() == 0 {
					// fmt.Printf("%s└───%s (empty)\n", prefix, file.Name())
					fmt.Fprintf(out, "%s└───%s (empty)\n", prefix, file.Name())
				} else {
					// fmt.Printf("%s└───%s (%db)\n", prefix, file.Name(), file.Size())
					fmt.Fprintf(out, "%s└───%s (%db)\n", prefix, file.Name(), file.Size())
				}
			}
		} else {
			if file.IsDir() {
				// fmt.Print(prefix+"├───", file.Name(), "\n")
				fmt.Fprint(out, prefix+"├───", file.Name(), "\n")
				dirTreeRec(out, path.Join(paths, file.Name()), prefix+"│\t", printFiles)
				continue
			} else {
				if file.Size() == 0 {
					// fmt.Printf("%s├───%s (empty)\n", prefix, file.Name())
					fmt.Fprintf(out, "%s├───%s (empty)\n", prefix, file.Name())
				} else {
					// fmt.Printf("%s├───%s (%db)\n", prefix, file.Name(), file.Size())
					fmt.Fprintf(out, "%s├───%s (%db)\n", prefix, file.Name(), file.Size())
				}
			}
		}

	}
}
