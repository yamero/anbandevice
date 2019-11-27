package command

import (
	"encoding/hex"
	"fmt"
	"net"
	"sort"
)

// 写入授权卡（排序区）
func (c *Command) WriteGrantCard(conn net.Conn, oldCardList []int, newCardList []int) {
	control := []byte{0x07, 0x07, 0x01}
	start := Int2hex(1, 8)
	data := []byte{}
	data = append(data, start...) // 新增的卡号从哪里开始写
	data = append(data, Int2hex(len(oldCardList) + len(newCardList), 8)...) // 每次新增几个卡号
	oldCardList = append(oldCardList, newCardList...)
	sort.Ints(oldCardList) // 卡号从小到大进行排序
	end, _ := hex.DecodeString("ffffffff881230235901010000ffff30000000000000000000000000")
	for _, tCard := range oldCardList {
		data = append(data, Int2hex(tCard, 10)...)
		data = append(data, end...)
	}
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	fmt.Printf("写入授权卡命令：%x\n", commandData)
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