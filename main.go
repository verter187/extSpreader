package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	oldP := "/home/wurtow977/Development/others/testfiles1/"
	newP := "/home/wurtow977/Development/others/newfiles1/"
	err := filepath.Walk(oldP,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				fname := filepath.Base(path)
				cExt := filepath.Ext(fname)
				cExt = strings.Replace(cExt, ".", "", -1)

				ctlgExt := filepath.Join(newP, cExt)

				err = os.MkdirAll(ctlgExt, 0766)
				fmt.Println(ctlgExt)
				if err == nil {
					err = os.Rename(filepath.Join(path), filepath.Join(ctlgExt, fname))
					if err != nil {
						fmt.Println(err)
						return err
					}
				}

				fmt.Println(fname)
				return nil
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

}
