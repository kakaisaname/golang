package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"compress/gzip"
)

func main() {
	//exec.Command ---------------
	//Command(name string, arg ...string) *Cmd　　　　//command返回cmd结构来执行带有相关参数的命令，它仅仅设定cmd结构中的Path和Args参数，
	// 如果name参数中不包含路径分隔符，command使用LookPath来解决路径问题，否则的话就直接使用name；Args直接跟在command命令之后
	cmd := exec.Command("php", "D:/baiduyun/command/application/cli", "read")
	var out bytes.Buffer
	cmd.Stdout = &out    ////标准输出

	//使某个命令开始执行，但是并不等到他执行结束，这点和Run命令有区别．然后使用Wait方法等待命令执行完毕并且释放响应的资源
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	fmt.Println(cmd.Args)
	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	fmt.Println(out.String())
}
