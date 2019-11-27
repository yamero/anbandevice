package command

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

// 读取排序区已存卡数量
func (c *Command) GetGrantedCardInfo(conn net.Conn) int {
	control := []byte{0x07, 0x01, 0x00}
	data := []byte{}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，读取排序区已存卡数量失败！")
		return 0
	}
	var recvMsg [1024]byte
	var allRecv []byte
	for {
		n, err := conn.Read(recvMsg[:])
		if err != nil || n <= 0 {
			fmt.Println("无法接收数据，读取排序区已存卡数量失败！")
			return 0
		}
		allRecv = append(allRecv, recvMsg[:n]...)
		if bytes.HasSuffix(allRecv, []byte{0x7e}) {
			break
		}
	}
	num, _ := strconv.ParseUint(fmt.Sprintf("%x", allRecv[36:40]), 16, 32)
	return int(num)
}
