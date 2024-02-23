package os_ext

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DeleteLocalFiles 删除本地文件
func DeleteLocalFiles(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("DeleteLocalFiles error（path：）：", err)
		return err
	} else {
		fmt.Println("DeleteLocalFiles success（path：）")
	}

	return nil
}

// CreateDir 创建文件夹
func CreateDir(dir string) error {
	if dir == "" {
		return errors.New("dir is empty")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		} else {
			fmt.Println("Directory created successfully")
		}
	} else {
		fmt.Println("Directory already exists")
	}

	return nil
}

// Zip zip 压缩
func Zip(top string, srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	var top1 = top + "\\"
	var top2 = top + "/"

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		var p = strings.Replace(path, top1, "", -1)
		p = strings.Replace(p, top2, "", -1)
		header.Name = p
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

// DeCompressZip zip 解压
func DeCompressZip(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := filepath.Join(dest, file.Name)
		dir, err := getDir(filename)
		if err != nil {
			return err
		}

		err = CreateDir(dir)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

// 传入文件路径，返回目录
func getDir(path string) (string, error) {
	sep := runtime.GOOS
	if sep == "windows" {
		sep = "\\"
	} else {
		sep = "/"
	}
	return SubString(path, 0, strings.LastIndex(path, sep))
}

// SubString 截取字符串
func SubString(str string, start, end int) (string, error) {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return "", errors.New("start is wrong")
	}

	if end < start || end > length {
		return "", errors.New("end is wrong")
	}

	return string(rs[start:end]), nil
}
