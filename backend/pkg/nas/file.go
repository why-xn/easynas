package nas

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileInfo holds the path and size of a file or folder
type FileInfo struct {
	Name string
	Path string
	Size int64 // Size is 0 for directories
}

// ListAndSortFilesFolders lists all folders and files under a path,
// sorting folders first and files second, both in ascending order by name,
// and includes file sizes.
func ListAndSortFilesFolders(path string) ([]FileInfo, error) {
	var folders []FileInfo
	var files []FileInfo

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	// Walk through the directory structure.
	err := filepath.Walk(path, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root path itself
		if fullPath == path {
			return nil
		}

		// Collect file/folder info
		if info.IsDir() {
			folders = append(folders, FileInfo{Name: info.Name(), Path: fullPath, Size: 0})
		} else {
			files = append(files, FileInfo{Name: info.Name(), Path: fullPath, Size: info.Size()})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort folders and files individually by name
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Path < folders[j].Path
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	// Combine folders and files, with folders first
	result := append(folders, files...)

	return result, nil
}
