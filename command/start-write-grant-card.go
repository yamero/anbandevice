package command

import (
	"bytes"
	"fmt"
	"net"
)

// 作好写入授权卡准备（排序区）
func (c *Command) StartWriteGrantCard(conn net.Conn) {
	control := []byte{0x07, 0x07, 0x00}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，开启授权卡写入缓冲区失败！")
		return
	}
	var recvMsg [1024]byte
	var allRecv []byte
	for {
		n, err := conn.Read(recvMsg[:])
		if err != nil || n <= 0 {
			fmt.Println("无法接收信息，开启授权卡写入缓冲区失败！")
			return
		}
		allRecv = append(allRecv, recvMsg[:n]...)
		if bytes.HasSuffix(allRecv, []byte{0x7e}) {
			break
		}
	}
	sta := fmt.Sprintf("%x", allRecv[25:28])
	if sta == returnOk || sta == "3703ff" {
		fmt.Println("开启授权卡写入缓冲区成功！")
	} else {
		fmt.Println("开启授权卡写入缓冲区失败！")
	}
}