//rand 包实现了一个密码安全的伪随机数生成器。
//Reader是一个密码强大的伪随机生成器的全球共享实例。------------------------
//
//Read 是一个使用 io.ReadFull 调用 Reader.Read 的辅助函数。返回时，n == len(b) 当且仅当 err == nil。
//
//示例
//本例从 rand.Reader 中读取10个密码安全的伪随机数，并将它们写入字节片。

package main

import (
"bytes"
"crypto/rand"
"fmt"
)

func main() {
c := 10
b := make([]byte, c)
_, err := rand.Read(b)
if err != nil {
fmt.Println("error:", err)
return
}
// 切片现在应该包含随机字节而不是仅包含零。
fmt.Println(bytes.Equal(b, make([]byte, c)))

}