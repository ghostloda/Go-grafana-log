package main

import (
	"time"
	"strings"
	"fmt"
	"os"
	"bufio"
	"io"
	"regexp"
	"log"
	"strconv"
	"net/url"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc chan []byte
	wc chan *Message
	read Reader
	write Writer
}

type ReadFromfile struct {
	path string //读取文件的路径
}

type WriteToInfluxDB struct {
	influxdbsn string // influxdb data source
}

type Message struct {
	TimeLocal                 time.Time
	BytesSent                 int
	Path,Method,Scheme,Status string
	UpstreamTime,RequestTime  float64
}


func (r *ReadFromfile) Read(rc chan []byte)  {
	//读取数据

	f ,err := os.Open(r.path)
	if err != nil{
		panic(fmt.Sprintf("open file error:%s",err.Error()))
	}

	//从文件末尾开始逐行读取文件内容

	//将字符指针移动到末尾
	f.Seek(0,2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		//如果读取到文件末尾会返回io.EOF的错误，添加sleep等待新的日志产生
		if err == io.EOF {
			time.Sleep(500*time.Millisecond)
			continue
		}else if err != nil{
			panic(fmt.Sprintf("readBytes error:%s",err.Error()))
		}

		rc <- line[:len(line)-1]

	}

}

func (l *LogProcess) Process()  {
	//解析数据

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc{
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14{
			log.Println("匹配失败：",string(v))
			continue
		}

		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000",ret[4],loc)
		if err != nil{
			log.Println("时间解析失败：",err.Error(),ret[4])
			continue
		}

		message.TimeLocal = t

		byteSent, _ := strconv.Atoi(ret[8])

		message.BytesSent = byteSent
		//GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6]," ")
		if len(reqSli) != 3{
			log.Println("字符串分割失败：",ret[6])
			continue
		}
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil{
			log.Println("url 解析失败:",err)
			continue
		}
		message.Path = u.Path
		message.Scheme = ret[5]
		message.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12],64)
		requestTime, _ := strconv.ParseFloat(ret[13],64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime


		l.wc <- message

	}

}

func (w *WriteToInfluxDB) Write(wc chan *Message)  {
	//写入数据库

	for v := range wc{
		fmt.Println(v)
	}

}




func main()  {
	r := &ReadFromfile{
		path: "F:\\GitPro\\Go-grafana-log\\access.log",
	}
	w := &WriteToInfluxDB{
		influxdbsn: "username@password..",
	}

	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan *Message),
		read: r,
		write: w,
	}

	go lp.read.Read(lp.rc)

	go lp.Process()

	go lp.write.Write(lp.wc)


	time.Sleep(30*time.Second)
	
}

