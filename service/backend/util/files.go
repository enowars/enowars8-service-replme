package util

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	cp "github.com/otiai10/copy"
)

func MakeDirIfNotExists(path string) error {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	if err != nil {
		return err
	}

	if stat.IsDir() {
		return nil
	} else {
		return errors.New("Target path is file")
	}
}

func TouchIfNotExists(dir string, name string) error {
	err := MakeDirIfNotExists(dir)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name)

	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		_, err := os.Create(path)
		return err
	}

	if stat.IsDir() {
		return errors.New("Target path is dir")
	} else {
		return nil
	}
}

func SetFileContent(dir string, name string, content string) error {
	err := TouchIfNotExists(dir, name)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name)

	return os.WriteFile(path, []byte(content), 0600)
}

func GetFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

func CopyRecurse(src string, target string, perm fs.FileMode) error {
	return cp.Copy(src, target, cp.Options{
		OnSymlink:     func(string) cp.SymlinkAction { return cp.Skip },
		AddPermission: perm,
	})
}

func GetFileModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	return info.ModTime(), nil
}

func DeleteFilesOlderThan(path string, t time.Time) {
	dir, err := os.Open(path)
	if err != nil {
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			child := filepath.Join(path, file.Name())
			modTime, err := GetFileModTime(child)
			if err != nil {
				continue
			}
			if modTime.Before(t) {
				SLogger.Debugf("Deleting file %s", child)
				err = DeleteFile(child)
				if err != nil {
					continue
				}
			}
		}
	}
}

func DeleteDirsOlderThan(path string, t time.Time) {
	dir, err := os.Open(path)
	if err != nil {
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			child := filepath.Join(path, file.Name())
			modTime, err := GetFileModTime(child)
			if err != nil {
				continue
			}
			if modTime.Before(t) {
				SLogger.Debugf("Deleting directory %s", child)
				err = DeleteDir(child)
				if err != nil {
					continue
				}
			}
		}
	}
}
