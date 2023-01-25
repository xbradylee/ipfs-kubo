//go:generate npm run build --prefix ./dir-index-html/
package assets

import (
	"embed"
	"io"
	"io/fs"
	"strconv"

	"github.com/cespare/xxhash"
)

//go:embed dir-index-html/dir-index.html dir-index-html/knownIcons.txt
var Asset embed.FS

// AssetHash a non-cryptographic hash of all embedded assets
var AssetHash string

func init() {
	sum := xxhash.New()
	err := fs.WalkDir(Asset, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := Asset.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(sum, file)
		return err
	})
	if err != nil {
		panic("error creating asset sum: " + err.Error())
	}

	AssetHash = strconv.FormatUint(sum.Sum64(), 32)
}
