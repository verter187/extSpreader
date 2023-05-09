package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func moveFiles(newP string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// Walk will not walk into path directory
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fname := filepath.Base(path)
			cExt := filepath.Ext(fname)
			cExt = strings.TrimLeft(cExt, ".")

			ctlgExt := filepath.Join(newP, cExt)

			err = os.MkdirAll(ctlgExt, 0766)

			if err == nil {
				err = os.Rename(filepath.Join(path), filepath.Join(ctlgExt, fname))
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func main() {

	oldP := "/home/wurtow977/Development/others/testfiles1/"
	newP := "/home/wurtow977/Development/others/newfiles1/"

	err := filepath.Walk(oldP, moveFiles(newP))
	if err != nil {
		log.Fatal("888", err)
	}

}
