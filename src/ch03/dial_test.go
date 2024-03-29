package ch03

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	// 랜덤 포트에 리스너 생성
	listner, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)	
	}

	done := make(chan struct{})
	go func() {		
		defer func() { done <- struct{}{}}	()

		for {
			conn, err := listner.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("receieved: %q", buf[:n])
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", listner.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	
	conn.Close()
	<-done
	listner.Close()
	<-done
}

