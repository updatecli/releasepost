package changelog

import (
	"fmt"
	"os"
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

	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("creating file %s: %v", filename, err)
		return err
	}
	defer f.Close()
	f.Write(data)

	fmt.Printf("* %s created\n", filename)

	return nil
}
