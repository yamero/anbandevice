package command

import (
	"encoding/hex"
	"fmt"
	"net"
	"sort"
)

// 写入授权卡（排序区）
func (c *Command) WriteGrantCard(conn net.Conn, cardList []int, card int) {
	control := []byte{0x07, 0x07, 0x01}
	start := Int2hex(len(cardList) + 1, 8) // 起始序号同时也是每次写入的数量
	data := []byte{}
	data = append(data, start...)
	data = append(data, start...)
	cardList = append(cardList, card)
	sort.Ints(cardList)
	end, _ := hex.DecodeString("ffffffff881230235901010000ffff30000000000000000000000000")
	for _, tCard := range cardList {
		data = append(data, Int2hex(tCard, 10)...)
		data = append(data, end...)
	}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，写入授权卡失败！", err)
		return
	}
	var recvMsg [1024]byte
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，写入授权卡失败！")
		return
	}
	sta := fmt.Sprintf("%x", recvMsg[25:28])
	if sta == returnOk {
		fmt.Println("写入授权卡成功")
	} else {
		fmt.Println("写入授权卡失败！")
	}
}