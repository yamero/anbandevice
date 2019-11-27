package command

import (
	"fmt"
	"net"
	"strconv"
)

// 读取排序区已存卡数量
func (c *Command) GetGrantedCardInfo(conn net.Conn) {
	control := []byte{0x07, 0x01, 0x00}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，读取授权卡信息失败！")
		return
	}
	var recvMsg [1024]byte
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，读取授权卡信息失败！")
		return
	}
	num, _ := strconv.ParseUint(fmt.Sprintf("%x", recvMsg[36:40]), 16, 32)
	fmt.Println(num)
}
