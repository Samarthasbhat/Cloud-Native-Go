package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"errors"
)

type Effector func(context.Context) (string,error)


func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string,error){
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}
			log.Printf("Attempt %d failed: %v. Retrying in %s...", r+1, err, delay)

			select{
				case <- time.After(delay):
				case <- ctx.Done():
					return "", ctx.Err()
			}
		}
	}
}

var count int
 
func EmulateUnreliableService(ctx context.Context) (string, error){
	count ++ 

	if count <= 3 {
		return "intentional fail", errors.New("Error")
	}else{
		return "success", nil
	}
}

func main(){
	r := Retry(EmulateUnreliableService, 2, 2*time.Second)

	res, err := r(context.Background())

	fmt.Println(res, err)
}