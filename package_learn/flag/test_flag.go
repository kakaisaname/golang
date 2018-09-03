package main

import (
	"flag"
	"fmt"
)

var Input_Name = flag.String("name", "dabaojian", "input your name")
var Input_Age = flag.Int("age", 20, "input your age")
var Input_Gender = flag.String("gender", "male", "input your age")
var Input_flagvar int

func Init() {
	flag.IntVar(&Input_flagvar, "flagname", 1234, "help message for flagname")
}

func main() {
	Init()
	flag.Parse()

	fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())
	fmt.Printf("flag num=%d\n", flag.NFlag())

	for i := 0; i != flag.NArg(); i++ {
		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}

	fmt.Println("name=", *Input_Name)
	fmt.Println("age=", *Input_Age)
	fmt.Println("gender=", *Input_Gender)
	fmt.Println("flagname=", Input_flagvar)
}

//go run  –name “11” -age=22 -flagname=0x22 fuck hit ds

//输出：
//args=[fuck hit ds], num=3
//flag num=3
//arg[0]=fuck
//arg[1]=hit
//arg[2]=ds
//name= 11
//age= 22
//gender= male
//flagname= 34