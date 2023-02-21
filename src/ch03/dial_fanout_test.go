// 여러 개의 서버에서 TCP를 통해 단 하나의 리소스만 받아 올 필요가 있다고 가정하자.
// 각 서버에 비동기적으로 연결을 요청하고, 동일한 콘텍스트를 각 다이얼러에 전달한다. 한 서버로부터
// 응답이 왔으면 다른 응답은 필요 없으니 나머지 다이얼러들은 콘텍스트를 취소하여 연결 시도를 중단할 수 있다.

package ch03

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second),
	)

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func(){
		// 하나의 연결만을 수락
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	dial := func(ctx context.Context, address string, response chan int, id int, wg *sync.WaitGroup){
		defer wg.Done()

		var d net.Dialer
		c, err := d.DialContext(ctx, "tcp", address)
		if err != nil{
			return
		}
		c.Close()

		select {
		case <- ctx.Done():
		case response <- id:
		}
	}
	res := make(chan int)
	var wg sync.WaitGroup

	for i:=0;i<10;i++{
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res
	cancel()
	wg.Wait()
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s", ctx.Err(),)
	}

	t.Logf("dialer %d retrieved the resource", response)
}