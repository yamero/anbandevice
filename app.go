package main

import (
	"anbandevice/cardinfo"
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	beginEndTag = "7e"
	returnOk = "210100" // 应答OK
	returnEnd = "3602ff" // 数据传输结束
	eventReadCard = "190100" // 读卡消息
	eventGoOutSwitch = "190200" // 出门开关消息
	eventDoorMagnetism = "190300" // 门磁消息
	eventRemoteOpenDoor = "190400" // 远程开门消息
	eventAlarm = "190500" // 报警消息
	eventSystem = "190600" // 系统消息
)

var (
	wg sync.WaitGroup
)

// 指令结构体
type command struct {
	begin byte // 开始标志码
	snPwd []byte // 设备sn和密码
	sn string // 设备SN，必须为16位长度
	password []byte // 设备通讯密码，必须为4位长度，默认为十六进制FFFFFFFF
	infoCode string // 发送端信息，必须为4位长度，最好为字母或数字
	controlCode []byte // 控制指令，必须为3位长度
	data []byte // 数据长度
}

func newCommand(sn string, pwd []byte, info string) *command {
	return &command{
		begin:       0x7e,
		sn:          sn,
		password:    pwd,
		infoCode:    info,
	}
}

// 将指令结构体转换为要发送的字节数据
func (c *command) getByteData() []byte {
	data := []byte{c.begin}
	data = append(data, c.sn...)
	data = append(data, c.password...)
	data = append(data, c.infoCode...)
	data = append(data, c.controlCode...)
	dataLen := computeDataLength(c.data)
	data = append(data, dataLen...)
	if len(c.data) > 0 {
		data = append(data, c.data...)
	}
	data = getFullData(data)
	return data
}

// 开启数据监控
func (c *command) openDataMonitor(conn net.Conn) {
	control := []byte{0x01, 0x0b, 0x00}
	data := []byte{}
	c.controlCode = control
	c.data = data
	commandData := c.getByteData()
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
		values := url.Values{"card": {card}, "door": {door}, "status": {status}}
		ret := httpPostForm("http://127.0.0.1:8001/read_card", values)
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

// 设置开门时段
func (c *command) setOpenDoorTimes(conn net.Conn, times map[int][8]string) {
	control := []byte{0x06, 0x03, 0x00}
	data := []byte{0x01} // 总共64(0x40)组，每个设备只用到了第一组时间段
	c.controlCode = control
	for i := 1; i < 8; i++ {
		for _, t := range times[i] {
			if t == "0" {
				t = "00000000"
			} else {
				t = strings.ReplaceAll(t, ":", "")
				t = strings.ReplaceAll(t, "-", "")
			}
			h, _ := hex.DecodeString(t)
			data = append(data, h...)
		}
	}
	c.data = data
	commandData := c.getByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，开门时间段设置失败！")
		return
	}
	var recvMsg [1024]byte
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，开门时间段设置失败！")
		return
	}
	sta := fmt.Sprintf("%x", recvMsg[25:28])
	if sta == returnOk {
		fmt.Println("开门时间段设置成功")
	} else {
		fmt.Println("开门时间段设置失败！")
	}
}

// 获取所有开门时间段
func (c *command) getOpenDoorTimes(conn net.Conn)  {
	defer conn.Close()
	control := []byte{0x06, 0x02, 0x00}
	data := []byte{}
	c.controlCode = control
	c.data = data
	commandData := c.getByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，开门时间段获取失败！")
		return
	}
	var recvMsg [1024]byte
	bufReader := bufio.NewReader(conn)
	for {
		n, err := bufReader.Read(recvMsg[:])
		if err != nil {
			fmt.Println("接收数据失败！")
			continue
		}
		sta := fmt.Sprintf("%x", recvMsg[25:28])
		if sta == returnEnd {
			break
		}
		recvData := fmt.Sprintf("%x", recvMsg[:n])
		fmt.Println(recvData)
	}
}

// 获取TCP参数
func (c *command) getTcpParameter(conn net.Conn)  {
	control := []byte{0x01, 0x06, 0x00}
	data := []byte{}
	c.controlCode = control
	c.data = data
	commandData := c.getByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，TCP参数获取失败！")
		return
	}
	var recvMsg [1024]byte
	bufReader := bufio.NewReader(conn)
	n, err := bufReader.Read(recvMsg[:])
	if err != nil {
		fmt.Println("接收数据失败！")
		return
	}
	recvData := fmt.Sprintf("%x", recvMsg[:n])
	fmt.Println(recvData)
}

// 获得完整的命令数据
func getFullData(data []byte) []byte {
	// 计算检验码（除标志码和检验码，命令中所有字节都相加然后取尾子节）
	sum := 0
	for k, v := range data {
		if k == 0 { // 排除开始标志码
			continue
		}
		sum += int(v)
	}
	s1 := fmt.Sprintf("%x", sum)
	s2, _ := hex.DecodeString(s1[len(s1)-2:])
	data = append(data, s2[0], 0x7e)
	return data
}

// 计算数据长度，并返回数据长度byte切片
func computeDataLength(data []byte) []byte {
	s := fmt.Sprintf("%x", len(data))
	s = strings.Repeat("0", 8-len(s)) + s
	s1, _ := hex.DecodeString(s)
	return s1
}

// 执行post请求
func httpPostForm(url string, values url.Values) string {
	resp, err := http.PostForm(url, values)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// 执行get请求
// u = "http://www.xx.com"
func httpGet(u string) string {
	resp, err := http.Get(u)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

// 搜索设备，从设备返回的信息中，拿到设备sn和设备密码
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
	commandObj := newCommand("YC-0000000000000", []byte{0xff, 0xff, 0xff, 0xff}, "nini")
	commandObj.controlCode = []byte{0x01, 0x06, 0x00}
	commandObj.data = []byte{}
	sendData := commandObj.getByteData()
	m, err := udpConn.WriteToUDP(sendData, remoteAddr)
	if err != nil || m <= 0 {
		fmt.Println("无法搜索设备")
		os.Exit(1)
	}

	wg.Wait()
	udpConn.Close()

	deviceInfo := <- deviceInfoCh
	fmt.Printf("已搜索到设备，开始连接设备 %s\n", deviceInfo)

	conn, err := net.Dial("tcp", "192.168.1.150:8000")
	if err != nil {
		fmt.Println("无法连接设备！")
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("设备连接成功")

	commandObj = newCommand(deviceInfo, []byte{0xff, 0xff, 0xff, 0xff}, "nini")
	//commandObj.getTcpParameter(conn) // 获取TCP参数
	//commandObj.getOpenDoorTimes(conn) // 获取所有开门时间段
	commandObj.openDataMonitor(conn) // 开启数据监控

	// 一周七天，每天可以设置八个时间段，"0"表示不设置
	/*times := map[int][8]string{
		1: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周一
		2: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周二
		3: {"19:00-20:00", "0", "0", "0", "0", "0", "0", "0", }, // 周三
		4: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周四
		5: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周五
		6: {"19:00-20:00", "0", "0", "0", "0", "0", "0", "0", }, // 周六
		7: {"08:00-12:00", "0", "0", "0", "0", "0", "0", "0", }, // 周日
	}
	commandObj.setOpenDoorTimes(conn, times)*/

}

