package main


import (
	"fmt"
	"context"
	"sync"
	"time"
	
)

type Effector func(context.Context) (string,error)

func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector{
	var tokens = max
	var once sync.Once

	return func(ctx context.Context)  (string,error){
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func ()  {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for{
					select{
					case <- ticker.C:
						t := tokens + refill 
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})

		if tokens > 0 {
			return "", fmt.Errorf("too many calls")
	}
		tokens --
		return e(ctx)
	}
 }

 func main() {
	ef := func(ctx context.Context) (string,error){
		return "ok", nil
	}	

	t := Throttle(ef, 5, 1, time.Second * 5)

	for i := 0; i < 2; i++ {
		res, err := t(context.Background())
		fmt.Println(res, err)
	}
		
	time.Sleep(time.Second * 6)
	fmt.Println("After sleep")
 }