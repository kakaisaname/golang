//单文件压缩
package main

import (
	"fmt"
	"os"
	"log"
	"archive/tar"
	"io"
)

func main() {
	var srcFile = "D:/go/goprojects/src/archive/tar/password"

	var desFile = fmt.Sprintf("%s.tar",srcFile)

	//os.Create()这个函数是创见一个文件  创建的文件要关闭
	fw, err := os.Create(desFile)
	ErrPrintln(err)
	defer fw.Close()

	// 通过 fw 创建一个 tar.Writer  创建一个写入 参数为刚创建的文件
	//func NewWriter(w io.Writer) *Writer  //创建一个新的writer，向w中写入
	tw := tar.NewWriter(fw)
	//关闭
	defer func() {
		if err := tw.Close();err !=nil {
			ErrPrintln(err)
		}
	}()

	//os.Stat获取文件信息  -------
	//fileinfo, err := os.Stat(`C:\Users\Administrator\Desktop\UninstallTool.zip`)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(fileinfo.Name())    //获取文件名
	//fmt.Println(fileinfo.IsDir())   //判断是否是目录，返回bool类型
	//fmt.Println(fileinfo.ModTime()) //获取文件修改时间
	//fmt.Println(fileinfo.Mode())
	//fmt.Println(fileinfo.Size()) //获取文件大小
	//fmt.Println(fileinfo.Sys())
	fi, err := os.Stat(srcFile)
	ErrPrintln(err)

	//文件头信息------------
	// FileInfoHeader(fi os.FileInfo, link string) (*Header, error)//该函数通过os.fileInfo便可创建一个Header，Header中大部分内容自动填充，一些内容需要自己设定。
	//如自己写 header.qiao="wife" ...
	header, err := tar.FileInfoHeader(fi, "")

	//写入文件头----------------
	//func (tw *Writer) WriteHeader(hdr *Header) error//该函数将hdr写入tar文件中，如果hdr不是第一个header，该函数调用flush。在调用close之后在调用该函数就会报错ErrWriteAfterClose。
	err = tw.WriteHeader(header)
	ErrPrintln(err)

	//打开要写入的文件--------------- 之后要关闭
	//将文件数据写入
	//打开准备写入的文件
	fr, err := os.Open(srcFile)
	ErrPrintln(err)
	defer fr.Close()
	//copy 从fr中复制数据到tw中，知道所有数据复制完毕,返回复制的字节数和错误
	written, err := io.Copy(tw, fr)
	ErrPrintln(err)

	log.Printf("共写入了 %d 个字符的数据\n",written)
}

func ErrPrintln(err error)  {
	if err != nil {
		//Golang的标准库提供了log的机制，但是该模块的功能较为简单（看似简单，其实他有他的设计思路）。不过比手写fmt. Printxxx还是强很多的。至少在输出的位置做了线程安全的保护
		log.Println(err)
		//调用系统的 os.exit(1) 接口，退出程序返回状态为 “1”
		os.Exit(1)
	}
}
