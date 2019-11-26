package command

import (
	"bufio"
	"fmt"
	"net"
)

// 获取TCP参数
func (c *Command) GetTcpParameter(conn net.Conn)  {
	control := []byte{0x01, 0x06, 0x00}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
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