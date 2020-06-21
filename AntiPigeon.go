package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/lxn/walk"
)

type Cast struct {
	Time       time.Time
	UserName   string
	DeviceName string
	Status     bool
	NickName   string
	Addr       string
}

var replacer = strings.NewReplacer("\u0000", "", ":", " ")

const (
	online  = 6291457
	offline = 6291458
)

func main() {
	fmt.Println("飞鸽传书检测")
	fmt.Println("作者：书生中学信息社")
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 2425,
	})
	if err != nil {
		panic(err)
	}
	//go send()
	for {
		buf := make([]byte, 128)
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		str := replacer.Replace(string(buf))
		if ok, _ := regexp.MatchString("^1 [0-9]+ [a-zA-Z0-9]+ [a-z-A-Z0-9]+ [0-9]+ [\u4e00-\u9fa50-9a-zA-Z]+[\u0000]*$", str); !ok {
			continue
		}
		var (
			cast   Cast
			status int32
			t      int64
		)
		cast.Addr = addr.String()
		if _, err = fmt.Sscanf(str, "1 %d %s %s %d %s", &t, &(cast.UserName), &(cast.DeviceName), &status, &(cast.NickName)); err != nil {
			log.Println("Fail to scan:", err, cast)
			continue
		}
		if status == online {
			cast.Status = true
		}
		cast.Time = time.Unix(t, 0).Add(-time.Hour * 8)
		fmt.Printf("[%s]%s(%s):%s\n", cast.Time.Format("15:04:05"), cast.DeviceName, cast.Addr, func() string {
			if cast.Status {
				return "上线"
			}
			return "下线"
		}())
		if cast.Status {
			walk.MsgBox(nil, "飞鸽传书上线提示", fmt.Sprintf("计算机名称:%s\n飞鸽传书用户名:%s\nIP地址:%s", cast.DeviceName, cast.NickName, cast.Addr), walk.MsgBoxIconInformation)
		}
	}
}

func send() {
	time.Sleep(2 * time.Second)
	dial, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 2425,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("send")
	_, err = dial.Write([]byte("1:1592732570:baish:DESKTOP-9JBFHEM:6291458:bs"))
	if err != nil {
		panic(err)
	}
	fmt.Println("")

}
