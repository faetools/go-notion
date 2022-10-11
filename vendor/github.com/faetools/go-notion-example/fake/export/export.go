package export

import (
	"archive/zip"
	"bytes"
	_ "embed" // exports
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

//go:embed html/Export-c69b35a3-6d60-4e7c-a0f7-fbbc8d322841.zip
var HTMLZipped []byte

//go:embed md/Export-5b3b25f1-189e-4e13-8234-7adfc4e2132c.zip
var MarkdownCSVZipped []byte

func UnzipInto(zipped []byte, fs afero.Fs) error {
	archive, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	if err != nil {
		return err
	}

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			if err := fs.MkdirAll(f.Name, os.ModePerm); err != nil {
				return err
			}

			continue
		}

		if err := fs.MkdirAll(filepath.Dir(f.Name), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := fs.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
