package goactor

import (
	_ "reflect"
	"sync"
	"github.com/xuexihuang/goactor/chanrpc"
	"github.com/xuexihuang/goactor/log"
	"github.com/xuexihuang/goactor/module"
	"time"
)

type MiniActor struct {
	*module.Skeleton
	closeSig chan bool
	data     interface{}
	wg       sync.WaitGroup
}

func (mc *MiniActor) run() { 
	mc.Run(mc.closeSig)
	mc.wg.Done()
}

func (mc *MiniActor) OnInit() { //用来初始化data使用，并且注册消息和回调函数
   //miniactor.OnRegister("C2GSLogin", HandleC2GSLogin2)//demo,注册消息和处理函数,
}

func (mc *MiniActor) OnDestroy() { //析构data使用

}

func OnNew() *MiniActor {

	miniactor := MiniActor{}
	miniactor.Skeleton = &module.Skeleton{
		GoLen:              10000,//用来传输参数channel的长度
		TimerDispatcherLen: 10000,//没用到
		AsynCallLen:        10000,//没用到
		ChanRPCServer:      chanrpc.NewServer(10000),
	}
	miniactor.Skeleton.Init()
	miniactor.closeSig = make(chan bool, 1)
	miniactor.wg.Add(1)
	miniactor.OnInit()
	go miniactor.run()
	return &miniactor

}

func (mc *MiniActor) OnRegister(m interface{}, h interface{}) {//注册对于消息的处理函数
	mc.Skeleton.RegisterChanRPC(m, h)
}

func (m *MiniActor) OnEnd() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("%v", r)
		}
	}()

	m.closeSig <- true
	m.wg.Wait()
	m.OnDestroy()

}
/*  这是测试例子，在模块内部注册处理函数，例子中在main函数注册的，其实应该在OnInit注册
func HandleC2GSLogin2(i []interface{}) {

	m := i[0].(int)
	log.Debug("%v", m)
}

func main() {
	miniactor := OnNew()//动态创建一个actor实例
	miniactor.OnRegister("C2GSLogin", HandleC2GSLogin2)//注册消息和处理函数
	miniactor.Skeleton.ChanRPCServer.Go("C2GSLogin", 1212)//发送消息
	time.Sleep(time.Second * 10)
	miniactor.OnEnd()
	miniactor = nil

}
*/