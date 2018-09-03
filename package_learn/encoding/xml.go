package encoding

import (
	"os"
	"fmt"
	"encoding/xml"
	"strings"
)

//Encode(v interface{}) error可以把一个对像直接序列化到io.Writer对 像中。
//Decode(v interface{}) error 从 io.Reader 中,返序列化 xml


type Studn struct {
	Name string
	Age  int
}

func h() {
	f, err := os.Create("data.dat")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()  //---------------------------------------------------

	s := &Studn{"张三", 19} //创建 encode 对像
	encoder := xml.NewEncoder(f)
	//将 s 序列化到文件中   -----------------------------------------------
	encoder.Encode(s)

	//重置文件指针到开始位置
	f.Seek(0, os.SEEK_SET)  //------------------------------
	decoder := xml.NewDecoder(f) //------------------------
	var s1 Studn //从文件中反序列化成对像
	decoder.Decode(&s1)
	fmt.Println(s1)
}

//data.dat
//<Student><Name>张三</Name><Age>19</Age></Student>


//Marshal(v interface{}) ([]byte, error) 对象直接序列化成字符  -------------
func a1() {
	s := &Studn{"张三", 19}
	result, err := xml.Marshal(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(result))
}

//<Student><Name>张三</Name><Age>19</Age></Student>


type SS struct {
	XMLName xml.Name `xml:"student"`
	Name    string   `xml:"name"`
	Age     int      `xml:"age"`
}

func a2() {
	str := `<?xml version="1.0" encoding="utf-8"?>
           <student>
<name>张三</name> <age>19</age> </student>`
	var s SS
	xml.Unmarshal([]byte(str), &s)
	fmt.Println(s)
}


type Student struct {
	XMLName xml.Name `xml:"student"`
	Name    string   `xml:"name,attr"`
	Age     int      `xml:"age,attr"`
	Phone   []string `xml:"phones>phone",`
}
type ABC string

func main() {
	str := `<?xml version="1.0" encoding="utf-8"?>
			<student name="张三" age="19">
			   <phones>
				 <phone>12345</phone>
				 <phone>67890</phone>
			   </phones>
			</student>`
	var s Student
	xml.Unmarshal([]byte(str), &s)
	fmt.Println(s)
}


// 如果 带命名空间 <data:student name="张三" age="19">
//结果就为 {{data student} 张三 19 []}


//token方式，大文件快速  -----------------------

//在上面这几种方法中 Token 解析是最快的。对于大文件解析,或对性能有要 求时,这种方法是最佳选择。
type Student struct {
	Name string
	Age  int
}

func a4() {
	str :=
		"<?xml version=\"1.0\" encoding=\"utf-8\"?><Student><Name>张三</Name><Age>19</Age></Student>"
	decoder := xml.NewDecoder(strings.NewReader(str))
	var strName string
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			stelm := xml.StartElement(t)
			fmt.Println("Start ", stelm.Name.Local) // 如果带命名空间 stelm.Name.Space
			strName = stelm.Name.Local
		case xml.EndElement:
			endelem := xml.EndElement(t)
			fmt.Println("End ", endelem.Name.Local)
		case xml.CharData:
			data := xml.CharData(t)
			str := string(data)
			switch strName {
			case "Name":
				fmt.Println("姓名:", str)
			case "Age":
				fmt.Println("年龄:", str)
			default:
				fmt.Println("other:", str)
			}
		}
	}
}


//xml to struct  ---------------------------------------------
//将单层xml 属性，转换成struct属性


func a5() {

	str := `<Row SheetId="报告单ID" PatName="患者名称" Age="患者年龄" DeptName=" 科室名称" Sex="患者性别" SheetName="报告名称" SubName="子项目名称" CheckDate="检查日期" VerifyDocName="审核医生" PropDocName="申请医生"ResultStatus="结果状态" RepDocName="报告医生" ReportDate="报告日期"  SheetItemId="报告单ItemID" ItemName="明细项名称" ItemEnName="报告项英文名" Status="高低箭头(↑ ↓)" Indicator="-1：偏低  0：正常 1：偏高 2其他" Value="检查结果值" Unit="检查结果值单位" ValueRange="正常范围" ValueMin="参考最低值" ValueMax="参考最高值" />`

	ChangeXMLToStruct(str)
}

//xml 转换成一个Struct
func ChangeXMLToStruct(str string) {
	decoder := xml.NewDecoder(strings.NewReader(str))
	var strName string
	fmt.Println("strName----", strName)
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			//stelm := xml.StartElement(t)
			//fmt.Println("Start ", stelm.Name.Local) // 如果带命名空间 stelm.Name.Space
			//strName = stelm.Name.Local
			//strName = stelm.Name.Space
			//
			//fmt.Println(t.Attr)
			fmt.Println("-------------start---------------")
			for _, attr := range t.Attr {
				formart := "%s  string  `xml:\"%s,attr\"` //%s"
				fmt.Println(fmt.Sprintf(formart, attr.Name.Local, attr.Name.Local, attr.Value))
			}

			fmt.Println("-------------end---------------")

		case xml.EndElement:
			endelem := xml.EndElement(t)
			fmt.Println("End ", endelem.Name.Local)
		case xml.CharData:
			//先不管

		}
	}
}