package command

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
)

// 写入授权卡（排序区）
func (c *Command) WriteGrantCard(conn net.Conn, startSn int, card int) {
	control := []byte{0x07, 0x07, 0x01}
	data := []byte{}
	data = append(data, Int2hex(startSn, 8)...) // 起始序号，每次从 1 递增
	data = append(data, Int2hex(1, 8)...) // 每次向缓冲区写一个卡
	data = append(data, Int2hex(card, 10)...)
	end, _ := hex.DecodeString("ffffffff881230235901010000ffff30000000000000000000000000")
	data = append(data, end...)
	c.ControlCode = control
	c.Data = data
	commandData := c.GetByteData()
	newCommandData := []byte{}
	newCommandData = append(newCommandData, commandData[0:40]...)
	newCardByte := bytes.ReplaceAll(commandData[40:45], []byte{0x7e}, []byte{0x7f, 0x01})
	newCommandData = append(newCommandData, newCardByte...)
	newCommandData = append(newCommandData, commandData[45:]...)
	fmt.Printf("%x\n", newCommandData)
	_, err := conn.Write(newCommandData)
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