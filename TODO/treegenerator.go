package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// from github.com/floscodes/golang-tools.
// thanks to Florian on stackoverflow.com for providing it.
func CopyDir(src string, dest string) error {

	if len(dest) >= len(src) && dest[:len(src)] == src {
		return fmt.Errorf("Cannot copy a folder into the folder itself!")
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}

	file, err := f.Stat()
	if err != nil {
		return err
	}
	if !file.IsDir() {
		return fmt.Errorf("Source " + file.Name() + " is not a directory!")
	}

	err = os.Mkdir(dest, 0755)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {

		if f.IsDir() {

			err = CopyDir(src+"/"+f.Name(), dest+"/"+f.Name())
			if err != nil {
				return err
			}

		}

		if !f.IsDir() {

			content, err := ioutil.ReadFile(src + "/" + f.Name())
			if err != nil {
				return err

			}

			err = ioutil.WriteFile(dest+"/"+f.Name(), content, 0755)
			if err != nil {
				return err

			}

		}

	}

	return nil
}

func makeNewSubtree(name string) error {
	// replace spaces with dashes
	name = strings.Replace(name, " ", "-", -1)
	// lowercase all the letters
	name = strings.ToLower(name)

	// fix a re-referencing path with one that will work consistently
	name = path.Clean(name)

	fmt.Println(name)

	// make the directory for this entry
	err := CopyDir("../TEMPLATE", "../"+name)
	if err != nil {
		return fmt.Errorf("Unable to create folder: %s %s\n", name, err)
	}

	return nil
}

func main() {
	list, err := os.ReadFile("List.txt")
	if err != nil {
		panic(err)
	}

	var newEntries []string
	var writtenCount int

	entries := strings.Split(string(list), "\n")
	for i := 0; i < len(entries); i++ {
		err = makeNewSubtree(entries[i])
		if err == nil {
			writtenCount++
		} else {
			newEntries = append(newEntries, entries[i])
		}
	}

	// write the new List.txt file back out
	err = os.WriteFile("List.txt", []byte(strings.Join(newEntries, "\n")), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Operation successful. %d Subfolders Created.\n", writtenCount)
}
