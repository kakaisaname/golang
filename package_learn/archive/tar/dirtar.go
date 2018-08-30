//文件夹压缩
package main

import (
	"fmt"
	"os"
	"compress/gzip"
	"archive/tar"
	"path/filepath"
	"io"
	"log"
)

func main() {

	var src = "apt"
	var dst = fmt.Sprintf("%s.tar.gz",src)

	if err := Tar(src,dst);err != nil {
		log.Fatal(err)
	}
}

func Tar(src,dst string) (err error)  {
	//创建文件
	fw, err := os.Create(dst)
	if err != nil {
		return
	}
	defer fw.Close()

	// 将 tar 包使用 gzip 压缩，其实添加压缩功能很简单，
	// 只需要在 fw 和 tw 之前加上一层压缩就行了，和 Linux 的管道的感觉类似
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	//创建writer结构
	tw := tar.NewWriter(fw)
	defer tw.Close()

	//遍历指定目录（包括子目录），对遍历到的项目用 walkFn 函数进行处理
	return filepath.Walk(src, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return  err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return  err
		}

		//写入文件信息
		if err := tw.WriteHeader(hdr);err != nil {
			return  nil
		}

		//打开文件
		fr, err := os.Open(filename)
		defer fr.Close()
		if err != nil {
			return  err
		}
		n, err := io.Copy(tw, fr)
		if err != nil {
			return  err
		}

		log.Printf("成功打包 %s,共写入 %d 字节的数据\n",filename,n)
		return  nil
	})
}