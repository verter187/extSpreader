package main

import (
	"archive/zip"
	"io"

	"log"
	"os"
	"path/filepath"
	"strings"
)

func zipit(source, target string, needBaseDir bool) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			if needBaseDir {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			} else {
				path := strings.TrimPrefix(path, source)
				if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
					path = path[1:]
				}
				if len(path) == 0 {
					return nil
				}
				header.Name = path
			}
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

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
				err = os.Rename(path, filepath.Join(ctlgExt, fname))
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func main() {

	curP := "/home/wurtow977/Development/others/testfiles1/"
	ncurP := filepath.Base(curP)
	tmpP, err := os.MkdirTemp("", "tmpP")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer os.RemoveAll(tmpP) // clean up

	err = zipit(curP, filepath.Join(tmpP, "old_"+ncurP+".zip"), true)
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(curP, moveFiles(tmpP))
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(curP)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(tmpP, curP)
	if err != nil {
		log.Fatal(err)
	}

}
