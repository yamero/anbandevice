package command

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

// 获取所有授权卡（排序区）
func (c *Command) GetAllGrantedCard(conn net.Conn) []int {
	var cardList []int
	control := []byte{0x07, 0x03, 0x00}
	data := []byte{1} // 只获取排序区
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，授权卡信息获取失败！")
		return cardList
	}
	var recvMsg [1024]byte
	bufReader := bufio.NewReader(conn)
	var allRecv []byte
	for {
		n, err := bufReader.Read(recvMsg[:])
		if err != nil {
			fmt.Println("接收数据失败！")
			break
		}
		sta := fmt.Sprintf("%x", recvMsg[25:28])
		if sta == "3703ff" {
			break
		}
		allRecv = append(allRecv, recvMsg[:n]...)
	}
	if len(allRecv) > 32 {
		num, _ := strconv.ParseUint(fmt.Sprintf("%x", allRecv[32:36]), 16, 32)
		for i := 0; i < int(num); i++ {
			start := i * 33 + 36
			end := start + 5
			cardNum, _ := strconv.ParseUint(fmt.Sprintf("%x", allRecv[start:end]), 16, 40)
			cardList = append(cardList, int(cardNum))
		}
	}
	return cardList
}