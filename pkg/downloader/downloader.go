package downloader

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// DownloadImage 下载图片
func DownloadImage(imageURL, fpath string) error {
	// 自动创建文件夹
	if err := CheckDir(path.Dir(fpath)); err != nil {
		return err
	}
	out, err := os.Create(fpath)
	if err != nil {
		return err
	}
	// 获取图片
	res, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// 写出文件
	writer := bufio.NewWriter(out)
	_, err = io.Copy(writer, bytes.NewReader(b))
	return err
}

// CheckDir 检查文件夹是否存在
func CheckDir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else {
		err := os.MkdirAll(path, 0711)
		if err != nil {
			return err
		}
	}
	// check again
	_, err := os.Stat(path)
	return err
}
