package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	es "gopkg.in/olivere/elastic.v6"
	"log"
	"os"
	"reflect"
	"sync"
)

var (
	client *es.Client
	host 	= "http://127.0.0.1:9200/"
	once 	sync.Once
)

type Employee struct {
	FirstName 	string 		`json:"first_name"`
	LastName 	string 		`json:"last_name"`
	Age 		int			`json:"age"`
	About 		string		`json:"about"`
	Interests 	[]string 	`json:"interests"`
}

func init() {
	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	once.Do(func() {
		client, err = es.NewClient(es.SetErrorLog(errorlog), es.SetURL(host), es.SetSniff(false))
		if err != nil {
			checkError(err)
		}
		info, code, err := client.Ping(host).Do(context.Background())
		if err != nil {
			checkError(err)
		}
		fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
		esversion, err := client.ElasticsearchVersion(host)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Elasticsearch version %s\n", esversion)
	})
}

func main() {
	//create()
	//insertMany()
	//gets()
	query()
	//delete()
	//list(2,2)
	//update()
}

//创建数据
//http://localhost:9200/megacorp/employee/1     1、2、3 进行查询
func create()  {
	//使用结构体
	e1 := Employee{"kobe","bryant",40,"play basketball",[]string{"music","ball"}}
	put1, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("1").
		BodyJson(e1).
		Do(context.Background())
	if err != nil {
		checkError(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

	//使用字符串
	e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
	put2, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("2").
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		checkError(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)

	e3 := `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
	put3, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("3").
		BodyJson(e3).
		Do(context.Background())
	if err != nil {
		checkError(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put3.Id, put3.Index, put3.Type)
}

//批量增加
func insertMany()  {
	bulk := client.Bulk()
	defer client.Stop()
	docs := []Employee{
		Employee{"kobe1","bryant",40,"play basketball",[]string{"music","ball"}},
		Employee{"kobe2","bryant",40,"play basketball",[]string{"music","ball"}},
		Employee{"kobe3","bryant",40,"play basketball",[]string{"music","ball"}},
	}
	for _,doc := range docs {
		request := es.NewBulkIndexRequest().Index("megacorp").Type("employee").Doc(doc)
		bulk = bulk.Add(request)
	}
	bulk.Do(context.Background())
}

func gets()  {
	result, err := client.Get().Index("megacorp").Type("employee").Id("2").Do(context.Background())
	if err != nil {
		checkError(err)
	}
	if result.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", result.Id, result.Version, result.Index, result.Type)
		//var doc Employee								反序列化结构体没字段
		//json.Unmarshal(*result.Source,&doc)
		var doc map[string]interface{}					//反序列化为map  如果要反序列化为slice，就应该为  []map[string]interface{}
		json.Unmarshal(*result.Source,&doc)
		fmt.Println(doc)
	}
}

func query()  {
	var res *es.SearchResult
	var err error
	//取所有
	res, err = client.Search("megacorp").Type("employee").Do(context.Background())
	printEmployee(res, err)
	if res.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d Employee \n", res.Hits.TotalHits)
		var t Employee
		for _,hit := range res.Hits.Hits {
			fmt.Println(hit.Id)
			err := json.Unmarshal(*hit.Source,&t)
			if err != nil {
				fmt.Println("Deserialization failed")
			}
			fmt.Printf("Employee name %s : %s\n", t.FirstName, t.LastName)
			fmt.Println(t)
		}
	} else {
		fmt.Printf("Found no Employee \n")
	}

	//条件查询
	//年龄大于30岁的
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
	boolQ.Filter(elastic.NewRangeQuery("age").Gt(20))
	res, err = client.Search("megacorp").Type("employee").Query(boolQ).Do(context.Background())
	printEmployee(res, err)

	//短语搜索 搜索about字段中有 rock climbing
	matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	res, err = client.Search("megacorp").Type("employee").Query(matchPhraseQuery).Do(context.Background())
	printEmployee(res, err)

	//分析 interests
	//当使用到term 查询的时候，由于是精准匹配，所以查询的关键字在es上的类型，必须是keyword而不能是text，
	//比如你的搜索条件是 “name”:”jack”,那么该name 字段的es类型得是keyword，而不能是text											***
	aggs := elastic.NewTermsAggregation().Field("interests")
	res, err = client.Search("megacorp").Type("employee").Aggregation("all_interests", aggs).Do(context.Background())
	printEmployee(res, err)

}

func list(size,page int)  {
	if size < 0 || page < 1 {
		fmt.Printf("param error")
		return
	}
	res,err := client.Search("megacorp").
		Type("employee").
		Size(size).
		From((page-1)*size).
		Do(context.Background())
	printEmployee(res, err)
}

func update()  {
	res, err := client.Update().
		Index("megacorp").
		Type("employee").
		Id("2").
		Doc(map[string]interface{}{"age": 88}).
		Do(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("update age %s\n",res.Result)
}


//删除
func delete() {

	res, err := client.Delete().Index("megacorp").
		Type("employee").
		Id("1").
		Do(context.Background())
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("delete result %s\n", res.Result)
}

func printEmployee(res *es.SearchResult,err error)  {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var result Employee
	for _,item := range res.Each(reflect.TypeOf(result)) {
		t := item.(Employee)
		fmt.Printf("%#v\n",t)
	}
}

func checkError(err error)  {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
