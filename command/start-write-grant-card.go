package command

import (
	"fmt"
	"net"
)

// 作好写入授权卡准备（排序区）
func (c *Command) StartWriteGrantCard(conn net.Conn) {
	defer conn.Close()
	control := []byte{0x07, 0x07, 0x00}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，开启授权卡失败！")
		return
	}
	var recvMsg [1024]byte
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，开启授权卡失败！")
		return
	}
	sta := fmt.Sprintf("%x", recvMsg[25:28])
	if sta == returnOk {
		fmt.Println("开启授权卡成功")
	} else {
		fmt.Println("开启授权卡失败！")
	}
}