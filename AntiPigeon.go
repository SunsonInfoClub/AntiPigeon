package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
)

type Cast struct {
	Time       time.Time
	UserName   string
	DeviceName string
	Status     int
	NickName   string
	Addr       string
}

var replacer = strings.NewReplacer("\u0000", "", ":", " ")

const (
	online     = 6291457
	offline    = 6291458
	keepOnline = 6291459
)

func main() {
	fmt.Println("飞鸽传书检测")
	fmt.Println("作者：书生中学信息社")
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 2425,
	})
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	go send()
	for {
		buf := make([]byte, 128)
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		handleCast(buf, addr)

	}
}

func send() {
	time.Sleep(2 * time.Second)
	listenr, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 2424,
	})
	defer listenr.Close()
	if err != nil {
		panic(err)
	}
	go func() {
		conn, err := net.DialUDP("udp", &net.UDPAddr{
			Port: 2424,
		}, &net.UDPAddr{
			IP:   net.IPv4bcast,
			Port: 2425,
		})
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("发送检测广播")
		_, err = conn.Write([]byte("1:" + strconv.FormatInt(time.Now().Add(8*time.Hour).Unix(), 10) + ":baish:DESKTOP-9JBFHEM:6291457:bs2"))
		if err != nil {
			panic(err)
		}
	}()

	for {
		buf := make([]byte, 128)
		_, addr, err := listenr.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		handleCast(buf, addr)
	}
}

func handleCast(buf []byte, addr *net.UDPAddr) {
	if addr.Port == 2424 {
		return
	}
	str := replacer.Replace(string(buf))
	if ok, _ := regexp.MatchString("^1 [0-9]+ [a-zA-Z0-9]+ [a-z-A-Z0-9]+ [0-9]+ [\u4e00-\u9fa50-9a-zA-Z]+[\u0000]*$", str); !ok {
		return
	}
	var (
		cast Cast
		t    int64
	)
	cast.Addr = addr.String()
	if _, err := fmt.Sscanf(str, "1 %d %s %s %d %s", &t, &(cast.UserName), &(cast.DeviceName), &(cast.Status), &(cast.NickName)); err != nil {
		log.Println("Fail to scan:", err, cast)
		return
	}
	cast.Time = time.Unix(t, 0).Add(-time.Hour * 8)
	fmt.Printf("[%s]%s(%s):%s\n", cast.Time.Format("15:04:05"), cast.DeviceName, cast.Addr, func() string {
		switch cast.Status {
		case online:
			walk.MsgBox(nil, "飞鸽传书上线提示", fmt.Sprintf("计算机名称:%s\n飞鸽传书用户名:%s\nIP地址:%s", cast.DeviceName, cast.NickName, cast.Addr), walk.MsgBoxIconInformation)
			return "上线"
		case offline:
			return "下线"
		case keepOnline:
			return "在线上"
		}
		return "未知状态"
	}())
}
