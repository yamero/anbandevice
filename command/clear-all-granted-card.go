package command

import (
	"bytes"
	"fmt"
	"net"
)

// 清空所有授权卡（排序区）
func (c *Command) ClearAllGrantedCard(conn net.Conn) {
	control := []byte{0x07, 0x02, 0x00}
	data := []byte{1} // 只清空排序区
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，清空授权卡信息失败！")
		return
	}
	var recvMsg [1024]byte
	var allRecv []byte
	for {
		n, err := conn.Read(recvMsg[:])
		if err != nil || n <= 0 {
			fmt.Println("无法接收信息，清空授权卡信息失败！")
			return
		}
		allRecv = append(allRecv, recvMsg[:n]...)
		if bytes.HasSuffix(allRecv, []byte{0x7e}) {
			break
		}
	}
	sta := fmt.Sprintf("%x", allRecv[25:28])
	if sta == returnOk || sta == "3703ff" {
		fmt.Println("清空授权卡信息成功！")
	} else {
		fmt.Println("清空授权卡信息失败！")
	}
}