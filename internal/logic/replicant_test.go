package logic

import (
	"testing"
	"time"
)

func TestNewReplicant(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	_ = newReplicant(
		"host.docker.internal:8500",
		"logic.grpc",
	)
	//t.Logf("%+v", data)
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//t.Logf("ch")
			}
		}
	}()
	select {}
	// fmt.Println("TestNewReplicant")
}
