package module

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Attempts int = iota							// 0
	Retry										// 1
)

//定义结构体，保存后端服务器信息
type Backend struct {
	URL 			*url.URL
	Alive 			bool
	mux 			*sync.RWMutex
	ReverseProxy 	*httputil.ReverseProxy 		//反向代理  ReverseProxy 是一种 HTTP 处理器，它接收入向请求，
												// 将请求发送给另一个服务器，然后将响应发送回客户端。
}

//设置后端服务器的存活状态   设置服务器的存活状态，就用普通的锁
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

//判断后端服务器是否存活	判断存活状态，用读写锁
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

//跟踪所有后端，以及一个计算器变量 	后端服务池
type ServerPool struct {
	backends 		[]*Backend
	current 		uint64
}

//添加后端服务器到服务器池
func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends,backend)
}

//自动增加计数器并返回下一个服务器索引
//由于会有多客户端连接到负债均衡器，发生竞态条件，加锁不好，使用原子操作
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current,uint64(1)) % uint64(len(s.backends)))
}

//找到可用服务器，标记为当前可用服务器   返回下一个可用服务器
func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.backends) + next 				// start from next and move a full cycle
	for i := next;i<l;i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {			//判断是否存活，如果存活
			if i != next {
				atomic.StoreUint64(&s.current,uint64(idx))    //将 idx 存储为current
			}
			return s.backends[idx]				//返回存活的服务器
		}
	}
	return nil
}

//ping 服务器是否可以接通
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable,error: ",err)
		return false
	}
	_ = conn.Close()
	return true
}

//HealthCheck pings the backends and update the status   ping
func (s *ServerPool) HealthCheck() {
	for _,b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n",b.URL,status)
	}
}

//MarkBackendStatus changes a status of a backend
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _,b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

//Value 返回与此上下文相关联的键值，或为零  **
//默认为 1
//context维护重试次数，将重试次数传回lb
func GetAttemptsFromContext(r *http.Request) int {
	if attempts,ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

//默认为0
func GetRetryFromContext(r *http.Request) int {
	if retry,ok := r.Context().Value(Retry).(int);ok {
		return retry
	}
	return 0
}

// 对传入的请求进行负债均衡 ,如果已经达到最大上限，就终结这个请求。
func lb(w http.ResponseWriter,r *http.Request)  {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) max attempts reached,terminating\n",r.RemoteAddr,r.URL.Path)
		http.Error(w,"Service not available", http.StatusServiceUnavailable)
		return
	}
	peer := serverPool.GetNextPeer()																//如果下一个服务器是存活状态
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w,r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

//每20s进行一次健康状态检测
func healthCheck()  {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

var serverPool ServerPool

func main() {
	var serverList string
	var port int
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	//如果没有serverlist
	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}
	tokens := strings.Split(serverList, ",")
	for _,tok := range tokens {
		//http://localhost:8080/time?aaa=111&b=1212424   url.Parse 会被解析为  time?a=111&b=1212424
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}
		//返回了一个ReverseProxy对象，在ReverseProxy中的ServeHTTP方法实现了这个具体的过程，主要是对源http包头进行重新封装，而后发送到后端服务器。
		proxy := httputil.NewSingleHostReverseProxy(serverUrl)

		//主动模式，检查服务器是否运行正常  ReverseProxy 会触发 ErrorHandler 回调函数，我们可以利用它来检查故障。
		//闭包来实现错误处理器，它可以捕获外部变量错误。它会检查重试次数，如果小于 3，就把同一个请求发送给同一个后端服务器
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := GetRetryFromContext(request)
			//检查重试次数，如果小于 3，就把同一个请求发送给同一个后端服务器
			if retries < 3 {
				select {
				//重试时间间隔设定在 10 毫秒
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(),Retry,retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}
			// after 3 retries, mark this backend as down            设置这台服务器挂了
			serverPool.MarkBackendStatus(serverUrl, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			lb(writer, request.WithContext(ctx))
		}

		//添加服务器 ***   最开始都为可用的，默认
		serverPool.AddBackend(&Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}
	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),						//端口
		Handler: http.HandlerFunc(lb),
	}

	// start health checking		异步去检测		**
	go healthCheck()

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}