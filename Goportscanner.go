package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type portscanner struct {
	ip   string
	lock *semaphore.Weighted
}

func Ulimit() int64 {
	// out, err := exec.Command("ulimit", "-n").Output()
	cmd := exec.Command("sh", "-c", "ulimit -n")
	out, err := cmd.Output()

	if err != nil {
		panic(err)
	}
	s := strings.TrimSpace(string(out))
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		panic(err)
	}
	return i
}
func ScanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)

		} else {
			fmt.Println(port, "closed")
		}
		return
	}
	conn.Close()
	fmt.Println(port, "open")
}

func (ps *portscanner) Start(f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.ip, port, timeout)
		}(port)
	}
}
func ScanSinglePort() {

}
func main() {
	ps := &portscanner{
		ip:   "127.0.0.1",
		lock: semaphore.NewWeighted(Ulimit()),
	}
	var value int
	fmt.Println("Choose 1 for individual port scan or 2 ton scan all ports ")
	fmt.Scan(&value)

	switch value {
	case 1:
		var port int
		fmt.Println("Enter Port number")
		fmt.Scan(&port)
		ScanPort(ps.ip, port, 500*time.Millisecond)
	case 2:
		ps.Start(1, 443, 500*time.Millisecond)
	default:
		fmt.Println("Invalid Input")
	}

}
