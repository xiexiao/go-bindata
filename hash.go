// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// writeHash writes the table of file hash.
func writeHash(w io.Writer, toc []Asset) error {
	err := writeHashHeader(w)
	if err != nil {
		return err
	}

	var maxlen = 0
	for i := range toc {
		l := len(toc[i].Name)
		if l > maxlen {
			maxlen = l
		}
	}

	for i := range toc {
		err = writeHashAsset(w, &toc[i], maxlen)
		if err != nil {
			return err
		}
	}

	return writeTOCFooter(w)
}

// writeHashHeader writes the table of file hash header.
func writeHashHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, `// AssetHash get the file hash 
func AssetHash(name string) (string, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _binhash[canonicalName]; ok {
		return f, nil
	}
	return "", fmt.Errorf("Asset %%s not found", name)
}

// _binhash is a table, holding each asset hash, mapped to its name.
var _binhash = map[string]string{
`)
	return err
}

// writeHashAsset write a Hash entry for the given asset.
func writeHashAsset(w io.Writer, asset *Asset, maxlen int) error {
	hash, _ := getHash(asset.Path)
	_, err := fmt.Fprintf(w, "\t%q: \"%s\",\n", asset.Name, hash)
	return err
}

// getHash get the file hash
func getHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	hash := h.Sum(nil)
	return hex.EncodeToString(hash), nil
}

// writeHashFooter writes the table of contents file footer.
func writeHashFooter(w io.Writer) error {
	_, err := fmt.Fprintf(w, `}

`)
	return err
}
