//文件夹解压缩
package main

import (
"archive/tar"
"compress/gzip"
"fmt"
"io"
"os"
"path/filepath"
)

func main() {
	var dst = "" // 不写就是解压到当前目录
	var src = "log.tar.gz"

	UnTar(dst, src)
}

func UnTar(dst, src string) (err error) {
	// 打开准备解压的 tar 包
	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer fr.Close()

	// 将打开的文件先解压
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return
	}
	defer gr.Close()  //关闭文件

	// 通过 gr 创建 tar.Reader
	tr := tar.NewReader(gr)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:    //错误
			return err
		case hdr == nil:
			continue
		}

		// 处理下保存路径，将要保存的目录加上 header 中的 Name
		// 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
		//s := filepath.Join("a", "b", "", ":::", "  ", `//c////d///`)
		//fmt.Println(s) // a/b/:::/  /c/d   多个元素合并为一个路径
		dstFileDir := filepath.Join(dst, hdr.Name)

		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeDir: // 如果是目录时候，创建目录   目录
			// 判断下目录是否存在，不存在就创建
			if b := ExistDir(dstFileDir); !b {
				// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
				// 可以递归创建每一级目录
				if err := os.MkdirAll(dstFileDir, 0775); err != nil {
					return err
				}
			}
		case tar.TypeReg: // 如果是文件就写入到磁盘    文件
			// 创建一个可以读写的文件，权限就使用 header 中记录的权限
			// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下   os.FileMode 权限
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			n, err := io.Copy(file, tr)
			if err != nil {
				return err
			}
			// 将解压结果输出显示
			fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}

	return nil
}

// 判断目录是否存在
func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	//os.IsExist 返回一个布尔值，它指明err错误是否报告了一个文件或者目录已经存在
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}