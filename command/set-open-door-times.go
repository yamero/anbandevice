package command

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

// 设置开门时段
func (c *Command) SetOpenDoorTimes(conn net.Conn, times map[int][8]string) {
	control := []byte{0x06, 0x03, 0x00}
	data := []byte{0x01} // 总共64(0x40)组，每个设备只用到了第一组时间段
	c.ControlCode = control
	for i := 1; i < 8; i++ {
		for _, t := range times[i] {
			if t == "0" {
				t = "00000000"
			} else {
				t = strings.ReplaceAll(t, ":", "")
				t = strings.ReplaceAll(t, "-", "")
			}
			h, _ := hex.DecodeString(t)
			data = append(data, h...)
		}
	}
	c.Data = data
	commandData := c.GetByteData()
	_, err := conn.Write(commandData)
	if err != nil {
		fmt.Println("命令无法发送，开门时间段设置失败！")
		return
	}
	var recvMsg [1024]byte
	n, err := conn.Read(recvMsg[:])
	if err != nil || n <= 0 {
		fmt.Println("命令执行结果未知，开门时间段设置失败！")
		return
	}
	sta := fmt.Sprintf("%x", recvMsg[25:28])
	if sta == returnOk {
		fmt.Println("开门时间段设置成功")
	} else {
		fmt.Println("开门时间段设置失败！")
	}
}
