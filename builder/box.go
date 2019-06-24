package builder

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type box struct {
	Name     string
	Files    []file
	compress bool
}

func (b *box) Walk(root string) error {
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}
	if _, err := os.Stat(root); err != nil {
		// return nil
		return fmt.Errorf("could not find folder for box: %s", root)
	}
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() || strings.HasSuffix(info.Name(), "-packr.go") {
			return nil
		}
		name := strings.Replace(path, root+string(os.PathSeparator), "", 1)
		name = strings.Replace(name, "\\", "/", -1)
		f := file{
			Name: name,
		}

		DebugLog("packing file %s\n", f.Name)

		bb, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if b.compress {
			bb, err = compressFile(bb)
			if err != nil {
				return err
			}
		}
		bb, err = json.Marshal(bb)
		if err != nil {
			return err
		}
		f.Contents = strings.Replace(string(bb), "\"", "\\\"", -1)

		DebugLog("packed file %s\n", f.Name)
		b.Files = append(b.Files, f)
		return nil
	})
}

func compressFile(bb []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(bb)
	if err != nil {
		return bb, err
	}
	err = writer.Close()
	if err != nil {
		return bb, err
	}
	return buf.Bytes(), nil
}
