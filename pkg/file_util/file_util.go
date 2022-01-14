package file_util

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// EnsureDirectoryExists 确保给定的目录存在，如果不存在的话则会创建
func EnsureDirectoryExists(path string) error {
	if Exists(path) {
		if !IsDir(path) {
			return fmt.Errorf("路径%s是个文件", path)
		} else {
			return nil
		}
	}

	// 创建
	return os.MkdirAll(path, os.ModeDir)
}

// Exists 判断给定的路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir 判断给定的路径是否是目录
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// AppendLine 往给定的文件中追加一行
func AppendLine(filepath, lineContent string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(f)

	if _, err = f.WriteString(lineContent + "\n"); err != nil {
		panic(err)
	}
	return nil
}

// ReadLines 将给定路径的文件按行读取返回
func ReadLines(path string) ([]string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(fi)

	lineSlice := make([]string, 0)
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		lineSlice = append(lineSlice, string(a))
	}
	return lineSlice, nil
}
