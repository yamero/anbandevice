package command

import (
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
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，清空授权卡信息失败！")
		return
	}
	sta := fmt.Sprintf("%x", recvMsg[25:28])
	fmt.Printf("清空授权卡信息：%x\n", recvMsg[:n])
	if sta == returnOk || sta == "3703ff" {
		fmt.Println("授权卡信息已清空")
	} else {
		fmt.Println("授权卡信息清空失败！")
	}
}