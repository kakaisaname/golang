//md5 包实现了 RFC 1321 中定义的 MD5 哈希算法。

//MD5 是加密破解的，不应用于安全应用程序。

//func New() hash.Hash
//New 返回一个新的散列.Hash 计算 MD5 校验和。

//package main
//
//import (
//"crypto/md5"
//"fmt"
//"io"
//)
//
//func main() {
//	h := md5.New()
//	io.WriteString(h, "The fog is getting thicker!")
//	io.WriteString(h, "And Leon's getting laaarger!")
//	fmt.Printf("%x", h.Sum(nil))
//}


//package main
//
//import (
//"crypto/md5"
//"fmt"
//"io"
//"log"
//"os"
//)
//
//func main() {
//	f, err := os.Open("file.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	h := md5.New()
//	if _, err := io.Copy(h, f); err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("%x", h.Sum(nil))
//}

//总和返回数据的 MD5 校验和。
//package main
//
//import (
//	"crypto/md5"
//	"fmt"
//)
//
//func main() {
//	data := []byte("These pretzels are making me thirsty.")
//	fmt.Printf("%x", md5.Sum(data))
//}

//package main
//
//import (
//	"crypto/md5"
//	"fmt"
//	"io"
//)
//
//func main() {
//	str := "abc123"
//
//	//方法一
//	data := []byte(str)
//	has := md5.Sum(data)
//	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
//
//	fmt.Println(md5str1)
//
//	//方法二
//
//	w := md5.New()
//	io.WriteString(w, str)   //将str写入到w中
//md5str2: = fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
//
//	fmt.Println(mdtstr2)
//}
//
//
//
//打印结果：  md5.Sum(data)    md5.New()&&md5.Sum(nil)  ------------------------------------------------------------
//
//e99a18c428cb38d5f260853678922e03
//
//e99a18c428cb38d5f260853678922e03