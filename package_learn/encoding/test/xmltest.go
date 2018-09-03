package main

import (
	"encoding/xml"
	"fmt"
)

type SS struct {
	XMLName xml.Name `xml:"student"`
	Name    string   `xml:"name"`
	Age     int      `xml:"age"`
}

func main() {
	str := `<?xml version="1.0" encoding="utf-8"?>
           <student>
<name>张三</name> <age>19</age> </student>`
	var s SS
	xml.Unmarshal([]byte(str), &s)
	fmt.Println(s)  //{{ student} 张三 19}
}
