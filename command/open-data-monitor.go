package command

import (
	"anbandevice/cardinfo"
	myHttp "anbandevice/http"
	"bufio"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// 开启数据监控
func (c *Command) OpenDataMonitor(conn net.Conn) {
	control := []byte{0x01, 0x0b, 0x00}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，数据监控开启失败！")
		return
	}
	var recvMsg [1024]byte
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，数据监控开启失败！")
		return
	}
	sta := fmt.Sprintf("%x", recvMsg[25:28])
	if sta == returnOk {
		fmt.Println("数据监控成功开启")
		bufReader := bufio.NewReader(conn)
		for {
			fmt.Println("开始监控数据...")
			m, err := bufReader.Read(recvMsg[:])
			if err != nil || m <= 0 {
				fmt.Println("接收数据失败！")
				return
			}
			go handleEvent(recvMsg)
		}
	} else {
		fmt.Println("数据监控开启失败")
	}
}

// 处理实时监控消息（目前只处理读卡消息）
func handleEvent(recvMsg [1024]byte) {
	evt := fmt.Sprintf("%x", recvMsg[25:28])
	switch evt {
	case eventReadCard: // 读卡消息
		fmt.Println("读卡消息")
		cardData, _ := strconv.ParseUint(fmt.Sprintf("%x", recvMsg[32:37]), 16, 40)
		card := fmt.Sprintf("%d", cardData)
		doorData, _ := strconv.ParseUint(fmt.Sprintf("%x", recvMsg[43:44]), 16, 8)
		door := fmt.Sprintf("%d", doorData)
		statusData, _ := strconv.ParseUint(fmt.Sprintf("%x", recvMsg[44:45]), 16, 8)
		status := fmt.Sprintf("%d", statusData)
		fmt.Printf("卡号：%s %s 状态：%s\n", card, cardinfo.DoorInfo[door], cardinfo.CardStatus[status])
		values := url.Values{"k": {"d3be23e218c4f759be2eca23543dc243"}, "card": {card}, "door": {door}, "status": {status}}
		ret := myHttp.HttpPostForm("http://127.0.0.1:8001/read_card", values)
		fmt.Println(ret)
	case eventGoOutSwitch: // 出门开关消息
		fmt.Println("出门开关消息")
	case eventDoorMagnetism: // 门磁消息
		fmt.Println("门磁消息")
	case eventRemoteOpenDoor:
		fmt.Println("远程开门消息")
	case eventAlarm:
		fmt.Println("报警消息")
	case eventSystem:
		fmt.Println("系统消息")
	default:
		fmt.Println("其他消息")
	}
}