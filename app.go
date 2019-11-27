package main

import (
	"anbandevice/command"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

// 搜索设备，从设备返回的信息中，拿到设备sn
func recvUdpMsg(conn *net.UDPConn, deviceInfoCh chan<- string)  {
	defer wg.Done()
	for {
		var recvMsg [1024]byte
		_, _, err := conn.ReadFromUDP(recvMsg[:])
		if err != nil {
			fmt.Println("接收广播数据失败")
			return
		}
		deviceInfoCh <- string(recvMsg[5:21])
		break
	}
}

func main()  {

	fmt.Println("开始搜索设备...")
	// IP地址要和设备在同已网段，端口随意，只要不是被占用的端口都可以
	localAddr, err := net.ResolveUDPAddr("udp", "192.168.1.152:8102")
	if err != nil {
		fmt.Println("无法搜索设备")
		os.Exit(1)
	}
	udpConn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Println("无法搜索设备")
		os.Exit(1)
	}
	wg.Add(1)
	deviceInfoCh := make(chan string, 1)
	go recvUdpMsg(udpConn, deviceInfoCh)

	// 开始广播
	remoteAddr := &net.UDPAddr{
		IP: net.ParseIP("255.255.255.255"),
		Port: 8101,
	}
	commandObj := command.NewCommand("YC-0000000000000", []byte{0xff, 0xff, 0xff, 0xff}, "nini")
	commandObj.ControlCode = []byte{0x01, 0x06, 0x00}
	commandObj.Data = []byte{}
	sendData := commandObj.GetByteData()
	m, err := udpConn.WriteToUDP(sendData, remoteAddr)
	if err != nil || m <= 0 {
		fmt.Println("无法搜索设备")
		os.Exit(1)
	}

	wg.Wait()
	udpConn.Close()

	deviceInfo := <- deviceInfoCh // 广播结束，收到设备返回的信息
	fmt.Printf("已搜索到设备，开始连接设备 %s\n", deviceInfo)

	conn, err := net.Dial("tcp", "192.168.1.150:8000")
	if err != nil {
		fmt.Println("无法连接设备！")
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("设备连接成功")

	commandObj = command.NewCommand(deviceInfo, []byte{0xff, 0xff, 0xff, 0xff}, "nini")
	//commandObj.GetTcpParameter(conn) // 获取TCP参数
	//commandObj.GetOpenDoorTimes(conn) // 获取所有开门时间段
	//commandObj.OpenDataMonitor(conn) // 开启数据监控

	/*num := commandObj.GetGrantedCardInfo(conn) // 获取排序区已存卡数量
	fmt.Printf("已保存卡数量：%d\n", num)*/
	cardList := commandObj.GetAllGrantedCard(conn) // 获取所有授权卡（排序区）
	fmt.Printf("已保存卡：%v\n", cardList)
	commandObj.ClearAllGrantedCard(conn) // 清空原来的所有授权卡（排序区每次写入必须先清空）
	time.Sleep(time.Second * 2)
	commandObj.StartWriteGrantCard(conn) // 开启缓冲区
	newCardList := []int{126, 127}
	cardList = append(cardList, newCardList...)
	sort.Ints(cardList)
	for i := 1; i <= len(cardList); i++ {
		commandObj.WriteGrantCard(conn, i, cardList[i-1]) // 写入授权卡 (缓冲区每次写入不可过多，否则会写入失败，这里每次写入1张卡)
	}
	time.Sleep(time.Second * 2)
	commandObj.EndWriteGrantCard(conn) // 结束写入授权卡，并更新授权卡信息
	time.Sleep(time.Second * 3)
	cardList1 := commandObj.GetAllGrantedCard(conn) // 获取所有授权卡（排序区）
	fmt.Printf("已保存卡：%v\n", cardList1)
	time.Sleep(time.Second * 1)

	// 一周七天，每天可以设置八个时间段，"0"表示不设置
	/*times := map[int][8]string{
		1: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周一
		2: {"08:00-19:00", "0", "0", "0", "0", "0", "0", "0", }, // 周二
		3: {"19:00-20:00", "0", "0", "0", "0", "0", "0", "0", }, // 周三
		4: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周四
		5: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周五
		6: {"19:00-20:00", "0", "0", "0", "0", "0", "0", "0", }, // 周六
		7: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周日
	}
	commandObj.SetOpenDoorTimes(conn, times)*/

}

