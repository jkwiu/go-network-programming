// 노드 간의 하트비트 구현으로 연결상의 문제가 생긴 경우 데이터가 전송될 때까지 기다리는 것이 아니라
// 네트워크상의 장애를 빠르게 파악하고 연결을 재시도할 수 있다.
// 또한, 하트비트는 TCP keepalive를 차단할 수 있는 방화벽에 대해서도 잘 작동한다.
// 먼저 일정한 간격으로 핑을 전송하기 위한 고루틴이 필요하고, 최근에 데이터를 받은 원격 노드로 불필요하게 또 다시 핑을 할 이유는 없으니,
// 핑 타이머를 초기화해야 한다.

package ch03

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration	
	select {
	case <- ctx.Done():
		return
	case interval = <-reset:
	default:
	}
	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval)
	defer func(){
		if !timer.Stop(){
			<-timer.C
		}
	}()

	for {
		select {
		case <- ctx.Done():
			return
		case newInterval := <-reset:
			if !timer.Stop(){
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}
		}

		_ = timer.Reset(interval)
	}
}