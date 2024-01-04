package changelog

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/updatecli/releasepost/internal/core/dryrun"
	"github.com/updatecli/releasepost/internal/core/result"
)

func initDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return fmt.Errorf("creating directory %s: %v", dirName, err)
		}
	}
	return nil
}

func dataToFile(data []byte, filename string) error {

	currentChecksum := getChecksumFromFile(filename)
	newChecksum := getChecksumFromByte(data)

	if !dryrun.Enabled {
		f, err := os.Create(filename)
		if err != nil {
			fmt.Printf("creating file %s: %v", filename, err)
			return err
		}
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			fmt.Printf("writing file %s: %v", filename, err)
		}
	}

	if currentChecksum == "" && newChecksum != "" {
		result.ChangelogResult.Created = append(result.ChangelogResult.Created, filename)
		return nil
	}

	if currentChecksum != newChecksum {
		result.ChangelogResult.Modified = append(result.ChangelogResult.Modified, filename)
	}

	result.ChangelogResult.UnModified = append(result.ChangelogResult.UnModified, filename)
	return nil
}

func getChecksumFromFile(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
		}

		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func getChecksumFromByte(raw []byte) string {
	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewBuffer(raw)); err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
