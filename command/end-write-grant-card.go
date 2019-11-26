package command

import (
	"fmt"
	"net"
)

// 结束写入授权卡，更新信息（排序区）
func (c *Command) EndWriteGrantCard(conn net.Conn) {
	defer conn.Close()
	control := []byte{0x07, 0x07, 0x02}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，更新授权卡信息失败！")
		return
	}
	var recvMsg [1024]byte
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，更新授权卡信息失败！")
		return
	}
	sta := fmt.Sprintf("%x", recvMsg[25:28])
	if sta == returnOk {
		fmt.Println("更新授权卡信息成功")
	} else {
		fmt.Println("更新授权卡信息失败！")
	}
}