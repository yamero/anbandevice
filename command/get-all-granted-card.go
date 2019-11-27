package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
)

// 获取所有已授权卡列表（排序区）
func (c *Command) GetAllGrantedCard(conn net.Conn) []int {
	var cardList []int
	control := []byte{0x07, 0x03, 0x00}
	data := []byte{1} // 只获取排序区
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，已授权卡列表获取失败！")
		return cardList
	}
	bufReader := bufio.NewReader(conn)
	var recvMsg [1024]byte
	var recvData []byte
	var tmp []byte
	for {
		n, err := bufReader.Read(recvMsg[:])
		if err != nil {
			fmt.Println("无法接收数据，已授权卡列表获取失败！")
			break
		}
		tmp = append(tmp, recvMsg[:n]...)
		if bytes.HasSuffix(tmp, []byte{0x7e}) {
			recvStatus := fmt.Sprintf("%x", tmp[25:28])
			if recvStatus == returnOk || recvStatus == "3703ff" {
				break
			}
			tmpDataLen, _ := strconv.ParseUint(fmt.Sprintf("%x", tmp[28:32]), 16, 32)
			dataStart := 36
			dataEnd := 32 + int(tmpDataLen)
			recvData = append(recvData, tmp[dataStart:dataEnd]...)
			tmp = []byte{}
		}
	}
	if len(recvData) >= 33 {
		num := len(recvData) / 33
		for i := 0; i < int(num); i++ {
			start := i * 33
			end := start + 5
			cardNum, _ := strconv.ParseUint(fmt.Sprintf("%x", recvData[start:end]), 16, 40)
			cardList = append(cardList, int(cardNum))
		}
	}
	return cardList
}