package main

import "fmt"

func main() {
	defer func() { //recover()必须放在defer func() 里才生效
		err := recover() //recover()返回panic信息，如果“本协程”内没有发生panic，则recover()返回nil
		if err != nil {
			fmt.Printf("panic信息: %s", err)
		}
	}()
	go panic("大乔乔") //panic必须发生在注册recover之后，且跟recover同协程内，才能被recover捕获
	fmt.Println("正常结束")
}
