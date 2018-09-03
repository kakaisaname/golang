package encoding
//encoding/json 是Go语言自带的JSON转换库
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"io"
	"log"
)

//转json   不需要流的场景
func a() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := json.Marshal(group)  //对象序列化成字符串
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)  //{“ID”:1,”Name”:”Reds”,”Colors”:[“Crimson”,”Red”,”Ruby”,”Maroon”]}
}

type Student struct {
	Name string `json:"userName"`
	Age  int
}

func d() {
	s := &Student{"张三", 19} //将 s 编码为 json
	buf, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(buf))//将 json 字符串转换成 Student 对像
	var s1 Student
	json.Unmarshal(buf, &s1)
	fmt.Println(s1)
}


//json解析  不需要流的场景
func b() {
	var jsonBlob = []byte(`[
        {"Name": "Platypus", "Order": "Monotremata"},
        {"Name": "Quoll",    "Order": "Dasyuromorphia"}
    ]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals) //字符串反序列号对象  Unmarshal 最大的特点就是,可以把 json 解析到一个 **map[string]interface{}**里
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals) //[{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}

func c() {
	str := `{"userName":"张三","Age":19}`
	var m map[string]interface{}
	json.Unmarshal([]byte(str), &m)
	for k, v := range m {
		switch v.(type) {
		case float64:
			fmt.Println(k, " 是 int 类型,值为:", v)
		case string:
			fmt.Println(k, " 是 string 类型,值为:", v)
		default:
			fmt.Println(k, "无法误用别的类型")
		}
	}
}

//

type St struct {
	Name string
	Age  int
}

//用于需要使用流的场景  Decode 字符串序列号成对象  Encode 对象序列化成字符串
func e() {
	f, err := os.Create("data.dat")
	if err != nil {
		fmt.Println(err)
	}
	s := &Student{"张三", 19} //创建 encode 对像
	encoder := json.NewEncoder(f) //将 s 序列化到文件中       --------------------
	encoder.Encode(s)

	//重置文件指针到开始位置
	f.Seek(0, os.SEEK_SET)
	decoder := json.NewDecoder(f)
	var s1 St //从文件中反序列化成对像
	decoder.Decode(&s1)
	fmt.Println(s1)
}

func f() {
	const jsonStream = `
        {"Name": "Ed", "Text": "Knock knock."}
        {"Name": "Sam", "Text": "Who's there?"}
        {"Name": "Ed", "Text": "Go fmt."}
        {"Name": "Sam", "Text": "Go fmt who?"}
        {"Name": "Ed", "Text": "Go fmt yourself!"}
    `
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}
	//输出--------------------------
	//Ed: Knock knock.
	//Sam: Who's there?
	//Ed: Go fmt.
	//		Sam: Go fmt who?
	//Ed: Go fmt yourself!
}

//json 中参数是数组
//当struct中需要数组[]interface{}形如这样的参数时，json传的格式是**"params":["hello world!"] **
//
//json中用数组用[]表示，前后不加"" 引号，内部不需要\ 转义。

//RPCParam rpc的data json数据
type RPCParam struct {
	Version   string        `json:"version"`
	User      string        `json:"user"`
	Password  string        `json:"password"`
	Timestamp int64         `json:"-"`
	Class     string        `json:"class"`
	Method    string        `json:"method"`
	Params    []interface{} `json:"params,omitempty"`
}

func g() {
	str := `{"version":"2.0","user":"","password":"","timestamp":1482723555,"class":"RpcClient_acsdispatcher","method":"serveraddr","params":["hello world!"]}`
	var n RPCParam
	err := json.Unmarshal([]byte(str), &n)
	fmt.Println(err, n)

}


