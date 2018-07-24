package main

import (
    "net"
    "log"
    "time"

    kcp "github.com/xtaci/kcp-go"
)

func main() {
    do_kcp()
    do_udp()
    do_tcp()
    do_channel()
}


/// ----------------------- kcp ----------------------
// 1. s监听
// 2. s等待接收
// 3. s接受c连接请求(发送数据)后,停止监听
// 4. s接收c的数据
// 5. s定时3秒后关闭socket
// 6. s由于socket关闭了,因此从阻塞的Read中立即返回err: broken pipe
// 7. s的连接关闭
/// --------------------------------------------------
func do_kcp() {
    log.Print("#################### kcp:")

    l, _ := kcp.ListenWithOptions(":48001", nil, 0, 0)
    go kcpClient()
    c, _ := l.Accept()
    l.Close()
    go kcpHandler(c)

    time.Sleep(3 * time.Second)
    
    c.Close()

    log.Print("wait for end......")
    time.Sleep(1 * time.Second)
}

func kcpClient() {
    time.Sleep(1 * time.Second)

    c, _ := kcp.DialWithOptions("127.0.0.1:48001", nil, 0, 0)
    n, _ := c.Write([]byte{1,2,3})
    log.Print("send ok len: ", n)
}

func kcpHandler(conn net.Conn) {
    defer conn.Close()
    b := make([]byte, 1024)
    for {
        log.Print("wait for read......")
        n, err := conn.Read(b)
        if err != nil {
            log.Print("server err != nil: ", err.Error())
            return
        }
        log.Print(b[:n])
    }
    log.Print("run here")
}


/// ----------------------- udp ----------------------
// 1. s建立套接字
// 2. s定时3秒后关闭socket
// 3. s由于socket关闭了,因此从阻塞的ReadFrom中立即返回err: read udp `IP:PORT`: use of closed network connection
/// --------------------------------------------------
func do_udp() {
    log.Print("#################### udp:")
    addr, _ := net.ResolveUDPAddr("udp", "192.168.10.156:0")
    conn, _ := net.ListenUDP("udp", addr)

    go udpTimer(conn)
    
    b := make([]byte, 1024)
    for {
        log.Print("waiting read...")
        n, _, err := conn.ReadFrom(b)
        if err != nil {
            log.Print("server err != nil: ", err.Error())
            return
        }
        log.Print(b[:n])
    }
}

func udpTimer(c *net.UDPConn) {
    time.Sleep(3 * time.Second)
    c.Close()
}

/// ----------------------- tcp ----------------------
// 1. s监听
// 2. s等待接收
// 3. s接受c连接请求后,停止监听
// 4. s接收c的数据
// 5. s定时3秒后关闭socket
// 6. s由于socket关闭了,因此在阻塞的Read中立即返回err: tcp `sIP:PORT`->`cIP:PORT`: use of closed network connection
// 7. 同时由于连接断开了,c也立即从阻塞的Read中返回err: EOF
/// --------------------------------------------------
func do_tcp() {
    log.Print("#################### tcp:")
    l, _ := net.Listen("tcp", ":7002")
    log.Print(l.Addr().String())
    go tcpClient()

    c, _ := l.Accept()
    l.Close()
    go tcpHandler(c)

    time.Sleep(3 * time.Second)

    c.Close()
    
    time.Sleep(1 * time.Second)

}

func tcpHandler(conn net.Conn) {
    defer conn.Close()
    b := make([]byte, 1024)
    for {
        n, err := conn.Read(b)
        if err != nil {
            log.Print("server err != nil: ", err.Error())
            return
        }
        log.Print(b[:n])
    }
}

func tcpClient() {
    time.Sleep(1 * time.Second)
    cc, _ := net.Dial("tcp", "127.0.0.1:7002")
    defer cc.Close()
    cc.Write([]byte{1,2,3})

    b := make([]byte, 1024)
    for {
        n, err := cc.Read(b)
        if err != nil {
            log.Print("client err != nil: ", err.Error())
            return
        }
        log.Print(b[:n])
    }
}


/// ----------------------- channel ----------------------
// 1. 建立同步管道
// 2. 向管道传数据
// 3. 从管道读取数据
// 4. 关闭管道
// 5. 由于管道关闭了,因此从阻塞的'<-ch'操作中立即返回ok==false
// 6. 循环结束
/// --------------------------------------------------
func do_channel() {
    log.Print("#################### channel:")
    ch := make(chan int)
    go chTimer(ch)
    OUT:
    for {
        select {
        case x, ok := <- ch:
            if !ok {
                log.Print("ch recv not ok: ", ok)
                break OUT
            }
            log.Print("ch recv: ", x)
        default:
        }
        time.Sleep(1 * time.Second)
    }
    log.Print("ch finish")
}

func chTimer(ch chan int) {
    time.Sleep(1 * time.Second)
    ch <- 10
    close(ch)
}
