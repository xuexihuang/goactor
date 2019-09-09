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

func (mc *MiniActor) OnInit() { //������ʼ��dataʹ�ã�����ע����Ϣ�ͻص�����
   //miniactor.OnRegister("C2GSLogin", HandleC2GSLogin2)//demo,ע����Ϣ�ʹ�����,
}

func (mc *MiniActor) OnDestroy() { //����dataʹ��

}

func OnNew() *MiniActor {

	miniactor := MiniActor{}
	miniactor.Skeleton = &module.Skeleton{
		GoLen:              10000,//�����������channel�ĳ���
		TimerDispatcherLen: 10000,//û�õ�
		AsynCallLen:        10000,//û�õ�
		ChanRPCServer:      chanrpc.NewServer(10000),
	}
	miniactor.Skeleton.Init()
	miniactor.closeSig = make(chan bool, 1)
	miniactor.wg.Add(1)
	miniactor.OnInit()
	go miniactor.run()
	return &miniactor

}

func (mc *MiniActor) OnRegister(m interface{}, h interface{}) {//ע�������Ϣ�Ĵ�����
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
/*  ���ǲ������ӣ���ģ���ڲ�ע�ᴦ��������������main����ע��ģ���ʵӦ����OnInitע��
func HandleC2GSLogin2(i []interface{}) {

	m := i[0].(int)
	log.Debug("%v", m)
}

func main() {
	miniactor := OnNew()//��̬����һ��actorʵ��
	miniactor.OnRegister("C2GSLogin", HandleC2GSLogin2)//ע����Ϣ�ʹ�����
	miniactor.Skeleton.ChanRPCServer.Go("C2GSLogin", 1212)//������Ϣ
	time.Sleep(time.Second * 10)
	miniactor.OnEnd()
	miniactor = nil

}
*/