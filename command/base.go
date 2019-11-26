package command

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	beginEndTag         = "7e"
	returnOk            = "210100" // 应答OK
	returnEnd           = "3602ff" // 数据传输结束
	eventReadCard       = "190100" // 读卡消息
	eventGoOutSwitch    = "190200" // 出门开关消息
	eventDoorMagnetism  = "190300" // 门磁消息
	eventRemoteOpenDoor = "190400" // 远程开门消息
	eventAlarm          = "190500" // 报警消息
	eventSystem         = "190600" // 系统消息
)

// 指令结构体
type Command struct {
	Begin       byte   // 开始标志码
	SnPwd       []byte // 设备sn和密码
	Sn          string // 设备SN，必须为16位长度
	Password    []byte // 设备通讯密码，必须为4位长度，默认为十六进制FFFFFFFF
	InfoCode    string // 发送端信息，必须为4位长度，最好为字母或数字
	ControlCode []byte // 控制指令，必须为3位长度
	Data        []byte // 数据长度
}

func NewCommand(sn string, pwd []byte, info string) *Command {
	return &Command{
		Begin:    0x7e,
		Sn:       sn,
		Password: pwd,
		InfoCode: info,
	}
}

// 将指令结构体转换为要发送的字节数据
func (c *Command) GetByteData() []byte {
	data := []byte{c.Begin}
	data = append(data, c.Sn...)
	data = append(data, c.Password...)
	data = append(data, c.InfoCode...)
	data = append(data, c.ControlCode...)
	dataLen := ComputeDataLength(c.Data)
	data = append(data, dataLen...)
	if len(c.Data) > 0 {
		data = append(data, c.Data...)
	}
	data = GetFullData(data)
	return data
}

// 获得完整的命令数据
func GetFullData(data []byte) []byte {
	// 计算检验码（除标志码和检验码，命令中所有字节都相加然后取尾子节）
	sum := 0
	for k, v := range data {
		if k == 0 { // 排除开始标志码
			continue
		}
		sum += int(v)
	}
	s1 := fmt.Sprintf("%x", sum)
	s2, _ := hex.DecodeString(s1[len(s1)-2:])
	data = append(data, s2[0], 0x7e)
	return data
}

// 计算数据长度，并返回数据长度byte切片
func ComputeDataLength(data []byte) []byte {
	return Int2hex(len(data), 8)
}

/**
将一个整数转换为指定位数的十六进制切片
@param num 要转换的数
@param n 位数（十六进制位，不是二进制位）
*/
func Int2hex(num int, n int) []byte {
	s := fmt.Sprintf("%x", num)
	s = strings.Repeat("0", n-len(s)) + s
	s1, _ := hex.DecodeString(s)
	return s1
}
