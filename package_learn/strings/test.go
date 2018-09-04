package main

import (
	"strings"
	"fmt"
	"os"
	"io/ioutil"
	"unicode"
)
func main() {
	var str = "zhaojunwei,您好"
	fmt.Println(strings.Repeat(str, 2))
	strPrt := fmt.Sprint(str)
	fmt.Println(strPrt)
	fmt.Println(strings.Count("cheese", "e"))                           //3
	fmt.Println(strings.Count("five", ""))                              //5
	fmt.Println(strings.EqualFold("Go", "go"))                          // true
	fmt.Println(strings.Contains(str, "zhao"))                          // true
	fmt.Println(strings.ContainsRune(str, '您'))                         // true
	fmt.Printf("Fields are: %q\n", strings.Fields("  foo bar  baz   ")) //Fields are: ["foo" "bar" "baz"]
	fmt.Println(strings.Replace(str, "zhao", "zhang", -1))              //zhangjunwei,您好
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	//Fields are: ["foo1" "bar2" "baz3"]
	fmt.Printf("Fields are: %q", strings.FieldsFunc("  foo1;bar2,baz3...", f))

	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26                                                              }
		return r
	}
	fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))
	reader := strings.NewReader("widuuv web")
	fmt.Printf("%#v\n",reader)
	fmt.Println(reader.Len())//10
	n, err := reader.Read(make([]byte, 10))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)//10 该值依据定义的字节数组的长度，如果定义数组长度为6，该值也为6
	reader1 := strings.NewReader("hello zhaojunwei")//
	b := make([]byte, 10)
	if n1, err := reader1.ReadAt(b, 2); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b[:n1]))
	}
	reader2 := strings.NewReader("hello shanghai China")
	b1 := make([]byte, 8)
	n2, _ := reader2.Read(b1)
	fmt.Println(string(b1[:n2])) //hello sh
	reader2.Seek(2, 1)
	n3,_ := reader2.Read(b1)
	fmt.Println(string(b1[:n3])) //ghai Chi
	reader3 := strings.NewReader("hello shanghai")
	b2 := make([]byte, 4)
	n4, _ := reader3.Read(b2)
	fmt.Println(string(b2[:n4])) //hell
	reader3.Seek(2, 1)
	reader3.UnreadByte()
	n5, _ := reader3.Read(b2)
	fmt.Println(string(b2[:n5]))// 空格sh

	reader4 := strings.NewReader("hello xuangubao")
	w, _ := os.Create("xugubao.txt")
	defer w.Close()
	n6, err := reader4.WriteTo(w)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n6) //15
	// ---------------Replacer--------------------
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
	fmt.Println(r.Replace("This is <b>HTML</b>")) //This is &lt;b&gt;HTML&lt;/b&gt;
	n7,err := r.WriteString(w, "This is <b>Html</b>!")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n7)//32
	d, _ := ioutil.ReadFile("xugubao.txt")
	fmt.Println(string(d))//hello xuangubaoThis is &lt;b&gt;Html&lt;/b&gt;!
}
