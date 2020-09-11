package files

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileStruct struct {
	RootPath string
}

func InitFileParser(pathRoot string) FileStruct {
	return FileStruct{
		RootPath: pathRoot,
	}
}

func (fs *FileStruct) BuildFilesTree() (files []string) {
	err := filepath.Walk(fs.RootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path == fs.RootPath {
				return nil
			}
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			files = append(files, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return
}
