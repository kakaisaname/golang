package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"sync"
	"time"
)

//type  BaseEntity interface {
//	GetId() 	string
//	SetId(id string)
//}

//type PageFilter struct {
//	SortBy     string
//	SortMode   int8
//	Limit      *int64
//	Skip       *int64
//	Filter     map[string]interface{}
//	RegexFiler map[string]string
//}

//要插入的数据
type Howie struct {
	HowieId 	primitive.ObjectID `bson:"_id"`
	Name 		string
	Pwd 		string
	Age 		int64
	CreateTime 	int64
	ExpiredTime time.Time
}

//mongo 的客户端
type MongoClient struct {
	Client *mongo.Client
	Ctx    context.Context
}

var Mongo *MongoClient

//初始化
func init() {
	//保证只执行一次 	**
	var once sync.Once
	once.Do(func() {
		//配置文件的路径 	***
		//conf, err := config.NewConfig("ini", "D:/project/src/go_dev/day11/config/logagent.conf")
		//if err != nil {
		//	fmt.Println("config failed,err: ",err)
		//	return
		//}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancel()

		//want, err := readpref.New(readpref.SecondaryMode) //表示只使用辅助节点
		//if err != nil {
		//	//checkErr(err)
		//}
		//wc := writeconcern.New(writeconcern.WMajority())
		//readconcern.Majority()

		//opt := options.Client().ApplyURI(conf.String("mongo::uri"))
		opt := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
		opt.SetLocalThreshold(3 * time.Second)     //只使用与mongo操作耗时小于3秒的
		opt.SetMaxConnIdleTime(5 * time.Second)    //指定连接可以保持空闲的最大毫秒数
		opt.SetMaxPoolSize(200)                    //使用最大的连接数
		//opt.SetReadPreference(want)                //表示只使用辅助节点
		//opt.SetReadConcern(readconcern.Majority()) //指定查询应返回实例的最新数据确认为，已写入副本集中的大多数成员
		//opt.SetWriteConcern(wc)                    //请求确认写操作传播到大多数mongod实例
		if client, err := mongo.Connect(ctx, opt);err == nil {
			Mongo = &MongoClient{}
			Mongo.Ctx = ctx
			Mongo.Client = client
		}
		if err = Mongo.Client.Ping(ctx,readpref.Primary());err != nil {
			fmt.Println("连接mongo错误")
			panic(err)
		} else {
			fmt.Println("连接成功")
		}
	})
}

var (
	howieArray = GetHowieArray()
	insertOneRes    *mongo.InsertOneResult
	insertManyRes   *mongo.InsertManyResult
	err 			error
	howie			Howie
	cursor          *mongo.Cursor
	howieArrayEmpty []Howie
	size            int64
	delRes          *mongo.DeleteResult
	updateRes       *mongo.UpdateResult
)

//删除文档
func drop()  {
	collection := Mongo.Client.Database("test").Collection("user")
	collection.Drop(Mongo.Ctx)
}

//插入数据
func (m *MongoClient) insertOne() {
	collection := Mongo.Client.Database("test").Collection("user")
	if insertOneRes, err = collection.InsertOne(Mongo.Ctx, howieArray[2]);err != nil {
		fmt.Println("insert err:",err)
	}
	fmt.Printf("InsertOne插入的消息ID:%v\n", insertOneRes.InsertedID)
}

//批量插入数据
func (m *MongoClient) insertMany() {
	collection := Mongo.Client.Database("test").Collection("user")
	if insertManyRes, err = collection.InsertMany(Mongo.Ctx, howieArray[2:]);err != nil {
		fmt.Println("insert err:",err)
	}
	fmt.Printf("InsertOne插入的消息ID:%v\n", insertManyRes.InsertedIDs)
}

func (m *MongoClient) find() {
	collection := Mongo.Client.Database("test").Collection("user")
	var Dinfo = make(map[string]interface{})
	err = collection.FindOne(Mongo.Ctx, bson.D{{"name", "howie_2"}, {"age", 11}}).Decode(&Dinfo)
	fmt.Println(Dinfo)
	fmt.Println("_id", Dinfo["_id"])
	//map[_id:ObjectID("5ddddd4482784a40d5e4c059") age:11 createtime:1574821188 expiredtime:1574821188505 name:howie_2 pwd:pwd_2]
	//_id ObjectID("5ddddd4482784a40d5e4c059")
}

func (m *MongoClient) findOne()  {
	collection := Mongo.Client.Database("test").Collection("user")
	if err = collection.FindOne(Mongo.Ctx, bson.D{{"name", "howie_2"}, {"age", 11}}).Decode(&howie);err != nil {
		fmt.Println(err)
	}
	fmt.Printf("FindOne查询到的数据:%v\n", howie)
	//FindOne查询到的数据:{ObjectID("5ddddd4482784a40d5e4c059") howie_2 pwd_2 11 1574821188 2019-11-27 02:19:48.505 +0000 UTC}
}

//查询单条数据后删除
func (m *MongoClient) findOneAndDelete()  {
	collection := Mongo.Client.Database("test").Collection("user")
	if err = collection.FindOneAndDelete(Mongo.Ctx, bson.D{{"name", "howie_3"}}).Decode(&howie); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("FindOneAndDelete查询到的数据:%v\n", howie)
	//FindOne查询到的数据:{ObjectID("5ddddd4482784a40d5e4c059") howie_2 pwd_2 11 1574821188 2019-11-27 02:19:48.505 +0000 UTC}
}

//查找单条数据后修改
func (m *MongoClient) findOneAndUpdate()  {
	collection := Mongo.Client.Database("test").Collection("user")
	if err = collection.FindOneAndUpdate(Mongo.Ctx, bson.D{{"name", "howie_4"}}, bson.M{"$set": bson.M{"name": "这条数据我需要修改了"}}).Decode(&howie); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("FindOneAndDelete查询到的数据:%v\n", howie)
	//FindOneAndDelete查询到的数据:{ObjectID("5dddde59e88439086d5f6e0e") howie_4 pwd_4 13 1574821465 2019-11-27 02:24:25.416 +0000 UTC}
}

//查询单条数据后替换该数据(以前的数据全部清空)
func (m *MongoClient) FindOneAndReplace()  {
	collection := Mongo.Client.Database("test").Collection("user")
	if err = collection.FindOneAndReplace(Mongo.Ctx, bson.D{{"name", "howie_5"}}, bson.M{"hero": "这条数据我替换了"}).Decode(&howie); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("FindOneAndReplace查询到的数据:%v\n", howie)
	//FindOneAndReplace查询到的数据:{ObjectID("5dddde59e88439086d5f6e0f") howie_5 pwd_5 14 1574821465 2019-11-27 02:24:25.416 +0000 UTC}
}

//一次查询多条数据
// 查询createtime>=3
// 限制取2条
// createtime从大到小排序的数据
func (m *MongoClient) findManyByCondition()  {
	collection := Mongo.Client.Database("test").Collection("user")
	if cursor, err = collection.Find(Mongo.Ctx, bson.M{"createtime": bson.M{"$gte": 2}}, options.Find().SetLimit(2), options.Find().SetSort(bson.M{"createtime": -1})); err != nil {
		fmt.Println(err)
	}
	if err = cursor.Err(); err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&howie); err != nil {
			fmt.Println(err)
		}
		howieArrayEmpty = append(howieArrayEmpty, howie)
	}
	for _, v := range howieArrayEmpty {
		fmt.Printf("Find查询到的数据ObejectId值%s 值:%v\n", v.HowieId.Hex(), v)
	}
	//Find查询到的数据ObejectId值5dddde59e88439086d5f6e0e 值:{ObjectID("5dddde59e88439086d5f6e0e") 这条数据我需要修改了 pwd_4 13 1574821465 2019-11-27 02:24:25.416 +0000 UTC}
	//Find查询到的数据ObejectId值5dddde59e88439086d5f6e10 值:{ObjectID("5dddde59e88439086d5f6e10") howie_6 pwd_6 15 1574821465 2019-11-27 02:24:25.416 +0000 UTC}
}

//查询集合中有多少条数据
func (m *MongoClient) findAll() {
	collection := Mongo.Client.Database("test").Collection("user")
	if size, err = collection.CountDocuments(Mongo.Ctx, bson.D{}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Count里面有多少条数据:%d\n", size)
}

//查询集合里面有多少数据(查询createtime>=3的数据)
func (m *MongoClient) findAllByCondition() {
	collection := Mongo.Client.Database("test").Collection("user")
	if size, err = collection.CountDocuments(Mongo.Ctx, bson.M{"createtime": bson.M{"$gte": 3}}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Count里面有多少条数据:%d\n", size)
}

//修改一条数据
func (m *MongoClient) updateOne() {
	collection := Mongo.Client.Database("test").Collection("user")
	if updateRes, err = collection.UpdateOne(Mongo.Ctx, bson.M{"name": "howie_2"}, bson.M{"$set": bson.M{"name": "我要改了他的名字"}}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("UpdateOne的数据:%d\n", updateRes)
	//UpdateOne的数据:&{1 1 0 <nil>}
}

//修改多条数据
func (m *MongoClient) updateMany() {
	collection := Mongo.Client.Database("test").Collection("user")

	if updateRes, err = collection.UpdateMany(Mongo.Ctx, bson.M{"createtime": bson.M{"$gte": 3}}, bson.M{"$set": bson.M{"name": "我要批量改了他的名字"}}); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("UpdateMany的数据:%d\n", updateRes)
}

//删除单条数据
func (m *MongoClient) deleteOne() {
	collection := Mongo.Client.Database("test").Collection("user")
	if delRes, err = collection.DeleteOne(Mongo.Ctx, bson.M{"name": "howie_1"}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("DeleteOne删除了多少条数据:%d\n", delRes.DeletedCount)

}

//删除多条数据
func (m *MongoClient) deleteMany() {
	collection := Mongo.Client.Database("test").Collection("user")
	if delRes, err = collection.DeleteOne(Mongo.Ctx, bson.M{"name": "howie_1"}); err != nil {
		fmt.Println(err)
	}
	if delRes, err = collection.DeleteMany(Mongo.Ctx, bson.M{"createtime": bson.M{"$gte": 7}}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("DeleteMany删除了多少条数据:%d\n", delRes.DeletedCount)
}

func main() {
	Mongo.deleteMany()
}

func GetHowieArray() (data []interface{}) {
	var i int64
	for i = 0;i<10;i++ {
		data = append(data,Howie{
			HowieId:     primitive.NewObjectID(),
			Name:        fmt.Sprintf("howie_%d",i+1),
			Pwd:         fmt.Sprintf("pwd_%d",i+1),
			Age:         i+10,
			CreateTime:  time.Now().Unix(),							// 1574820501    时间戳
			ExpiredTime: time.Now(), 								//2019-11-27 10:07:53.5185863 +0800 CST m=+0.017999801			具体的时间
		})
	}
	return
}

func checkErr(err error)  {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("没有查到数据")
			os.Exit(0)							//正常退出
		} else {
			fmt.Println("没有查到数据")
			os.Exit(0)
		}
	}
}
