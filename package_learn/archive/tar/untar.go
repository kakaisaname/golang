//单文件解压缩
package main

import (
	"log"
	"os"
	"archive/tar"
	"io"
)

func main() {
	var srcFile = "D:/go/goprojects/src/archive/tar/password.tar"

	fr, err := os.Open(srcFile)
	ErrPrint(err)
	defer fr.Close()

	//从fr中创建一个新的reader
	reader := tar.NewReader(fr)

	for hdr, err := reader.Next();err != io.EOF; hdr,err = reader.Next() {
		ErrPrint(err)

		//获取文件信息
		fi := hdr.FileInfo()
		//创建一个空文件，用来写入解包后的数据
		fw, err := os.Create("D:/go/goprojects/src/archive/tar/"+fi.Name())
		ErrPrint(err)
		//将reader 写入到fw中
		written, err := io.Copy(fw, reader)
		ErrPrint(err)
		log.Printf("解包：%s 到 %s,共处理了%d个字符的数据",srcFile,fi.Name(),written)

		//设置文件权限
		os.Chmod(fi.Name(),fi.Mode().Perm())
		fw.Close()
	}
}

func ErrPrint(err error)  {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
