package command

import (
	"bufio"
	"fmt"
	"net"
)

// 获取所有开门时间段
func (c *Command) GetOpenDoorTimes(conn net.Conn)  {
	defer conn.Close()
	control := []byte{0x06, 0x02, 0x00}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
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
