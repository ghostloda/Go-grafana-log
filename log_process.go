package main

import (
	"time"
	"strings"
	"fmt"
	"os"
	"bufio"
	"io"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc chan []byte
	wc chan string
	read Reader
	write Writer
}

type ReadFromfile struct {
	path string //读取文件的路径
}

type WriteToInfluxDB struct {
	influxdbsn string // influxdb data source
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

func (w *WriteToInfluxDB) Write(wc chan string)  {
	//写入数据库
	for v := range wc{
		fmt.Println(v)
	}


}

func (l *LogProcess) Process()  {
	//解析数据
	for v := range l.rc{

		l.wc <- strings.ToUpper(string(v))

	}

}


func main()  {
	r := &ReadFromfile{
		path: "F:\\数据结构与算法\\goLog\\access.log",
	}
	w := &WriteToInfluxDB{
		influxdbsn: "username@password..",
	}

	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan string),
		read: r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(30*time.Second)
	
}

