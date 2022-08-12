package utils

import (
	"io/ioutil"
	"os"
)

func IsFile(f string) (bool, error) {
	fi, err := os.Stat(f)
	if err != nil {
		return false, err
	}
	return !fi.IsDir(), nil
}

func IsDir(f string) (bool, error) {
	b, err := IsFile(f)
	if err != nil {
		return false, err
	}
	return !b, nil
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	//如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	//如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
	return false, err
}

// 获取根目录下直属所有文件（不包括文件夹及其中的文件）
func GetAllFile(pathname string) ([]string, error) {
	s := make([]string, 0)
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return nil, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func GetFiles(folder string) ([]string, error) {
	s := make([]string, 0)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			cs, err := GetFiles(folder + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			s = append(cs)
		} else {
			s = append(s, folder+"/"+file.Name())
		}
	}
	return s, nil
}
