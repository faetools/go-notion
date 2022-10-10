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

//go:embed html/Export-38cc96f0-cf48-4d26-ba29-a33d94cd6b66.zip
var HTMLZipped []byte

//go:embed md/Export-3795d6a1-6152-4a4d-90e3-de5272593bf3.zip
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
