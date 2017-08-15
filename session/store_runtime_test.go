package session

import (
	"testing"
	"fmt"
	"log"
	"time"
	"strconv"
)

var conf *StoreConfig
var runtime_store *RuntimeStore
var session_state *SessionState
var session_states []*SessionState
func init(){
	//log.Println("初始化")
	value := make(map[interface{}]interface{})
	value["foo"] = "bar"
	value["kak"] = "lal"

	conf = NewDefaultRuntimeConfig()
	runtime_store = NewRuntimeStore(conf)
	runtime_store.list.Init()

	session_state = NewSessionState(nil,"session_read",value)
	for i:=0;i<1000000;i++{
		session_states = append(session_states,NewSessionState(nil,"session_read"+strconv.Itoa(i),value))
		//runtime_store.SessionUpdate(NewSessionState(nil,"session_read"+strconv.FormatInt(time.Now().UnixNano(),10),value))
	}

	runtime_store.SessionUpdate(session_state)
	runtime_store.SessionUpdate(NewSessionState(nil,"session_read_1",value))
}

func TestRuntimeStore_SessionUpdate(t *testing.T) {
	//log.Println("开始 写测试")
	fmt.Println("-------------before update session state------------")
	state, _ := runtime_store.SessionRead("session_read")
	fmt.Printf("session state session_read:  %+v \n ", state)
	session_state.values["foo"] = "newbar"
	runtime_store.SessionUpdate(session_state)
	state, _ = runtime_store.SessionRead("session_read")
	fmt.Println("-------------after update session state------------")
	fmt.Printf("session state session_read:  %+v \n ",state)

}
func TestNewRuntimeStore_SessionUpdate_StateNotExist(t *testing.T) {
	fmt.Println("-------------before update session state------------")
	state, _ := runtime_store.SessionRead("session_read_2")
	fmt.Printf("session state session_read:  %+v \n ", state)
	state.values["make"] = "new"
	runtime_store.SessionUpdate(state)
	state, _ = runtime_store.SessionRead("session_read")
	fmt.Println("-------------after update session state------------")
	fmt.Printf("session state session_read:  %+v \n ",state)
}

func TestRuntimeStore_SessionRead(t *testing.T) {
	//log.Println("开始读测试")
	fmt.Printf("runtime_store:  %+v \n",*runtime_store)
	read,_ := runtime_store.SessionRead("session_read")
	if read == nil {
		fmt.Println("cannot find sessionId")
		return
	}
	fmt.Println("start read : ")
	fmt.Printf("sessionid : %v ,  values : %v \n", read.SessionID(),read.values)
}

func TestRuntimeStore_SessionExist(t *testing.T) {
	//log.Println("测试 session 存在")
	fmt.Println("is session exist: ", runtime_store.SessionExist("session_read"))

}

func TestRuntimeStore_SessionRemove(t *testing.T) {
	log.Println("session 删除测试")
	fmt.Println("------------------------")
	fmt.Println("before remove : ")
	read,err := runtime_store.SessionRead("session_read")
	if err != nil {
		panic(err)
	}
	fmt.Println("read : ")
	fmt.Printf("sessionid : %s ,  values : %v \n", read.SessionID(),read.values)

	err = runtime_store.SessionRemove("session_read")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("------------------------")
	fmt.Println("after remove : ")
	read,err = runtime_store.SessionRead("session_read")
	if err != nil {
		panic(err)
	}
	fmt.Println("read : ")
	fmt.Printf("sessionid : %s ,  values : %v \n", read.SessionID(),read.values)
}

func TestRuntimeStore_SessionGC(t *testing.T) {

}

func TestRuntimeStore_SessionCount(t *testing.T) {
	fmt.Println(runtime_store.SessionCount())
}

func TestRuntimeStore_SessionAccess(t *testing.T) {
	state ,_  := runtime_store.SessionRead("session_read")
	fmt.Println("------------------")
	fmt.Println("before session access")
	fmt.Println(state.timeAccessed.String())
	fmt.Println("------------------")
	fmt.Println("after session access")
	time.Sleep(10*time.Second)
	runtime_store.SessionAccess("session_read")
	fmt.Println(state.timeAccessed.String())

}


/**
	性能测试 | 基准测试
 */


func BenchmarkRuntimeStore_SessionRead_1(b *testing.B) {
	for i:=0;i<b.N;i++{
		runtime_store.SessionRead("session_read")
	}
	b.ReportAllocs()
}
func BenchmarkRuntimeStore_SessionRead_Parallel(b *testing.B) {
	b.RunParallel(func (pb *testing.PB){
		for pb.Next(){
			runtime_store.SessionRead("session_read")
		}
	})
	b.ReportAllocs()

}

func BenchmarkRuntimeStore_SessionCount_1(b *testing.B) {
	for i:=0;i<b.N ;i++  {
		runtime_store.SessionCount()
	}
}

func BenchmarkRuntimeStore_SessionCount_Parallel(b *testing.B) {
	b.RunParallel(func (pb *testing.PB){
		for pb.Next(){
			runtime_store.SessionCount()
		}
	})
	b.ReportAllocs()

}

func BenchmarkRuntimeStore_SessionRemove_1(b *testing.B) {

}

func BenchmarkRuntimeStore_SessionUpdate_1(b *testing.B) {
	for i:=0;i<b.N;i++{
		runtime_store.SessionUpdate(session_states[i])
	}
	b.ReportAllocs()
}

func BenchmarkRuntimeStore_SessionUpdate_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			runtime_store.SessionUpdate(session_state)
		}
	})
	b.ReportAllocs()
	fmt.Println(len(runtime_store.sessions))
}



