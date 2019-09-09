# goactor
利用goroutine和channel实现的类actor模块，简单易用
# 先来看例子
import (
	"github.com/xuexihuang/goactor/log"
	"github.com/xuexihuang/goactor"
	"time"
)
func HandleC2GSLogin2(i []interface{}) {

	m := i[0].(int)
	log.Debug("%v", m)
}

func main() {
	miniactor := OnNew()//动态创建一个actor实例
	miniactor.OnRegister("C2GSLogin", HandleC2GSLogin2)//注册消息和处理函数，实际开发中，这个注册动作放在MiniActor.go文件的OnInit函数中
	miniactor.Skeleton.ChanRPCServer.Go("C2GSLogin", 1212)//发送消息和参数
	time.Sleep(time.Second * 10)
	miniactor.OnEnd()
	miniactor = nil

}

如上就实现了动态创建一个actor，并通过goroutine通道发送消息和参数，actor模块内部异步实现计算，我们分部通过actor的几个特征来分析下封装的合理性
# 隔离性
actor内部和外部没有任何数据共享，内部的状态数据通过
type MiniActor struct {
	*module.Skeleton
	closeSig chan bool
	data     interface{}
	wg       sync.WaitGroup
}
的data来存储，生产中可以将interface{}换成需要的数据类型，这是actor内部的状态数据。因为属于小写变量所以外界也没办法访问。所有数据都是通过通道来传输。
# 异步性
通过channel来做邮箱，是异步的
# 串行计算
在内部是通过for{select case}结构实现的逻辑，是串行。从此上面封装符合actor思想的几大特征。
# 怎么使用如上模块
通过如上例子，也可以看出非常简单。MiniActor.go的OnInit函数注册相应的消息和处理函数，就完成了一个actor的设计，如果你程序中有多个不同的actor，就定义多个这样的go文件，内部注册相应的消息和处理函数。
同时同一类actor可以创建好多个，这种需求在很多业务中有需要，比如棋牌游戏的不同房价，每一个房间一个actor，可是这些actor内部是一样的处理逻辑。

