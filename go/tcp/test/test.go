package main

import (
	"fmt"
)

func toHex(ten int) (hex []int, length int) {
	m := 0

	hex = make([]int, 0)
	length = 0;

	for{
		m = ten / 16
		ten = ten % 16

		if(m == 0){
			hex = append(hex, ten)
			length++
			break
		}

		hex = append(hex, m)
		length++;
	}
	return
}

func main(){
	for a:=1;a<256;a++ {
		hex,length := toHex(a)
		for i:=0; i < length; i++ {
			if(hex[i] >= 10){
				fmt.Printf("%c",'A'+hex[i]-10)
			} else{
				fmt.Print(hex[i])
			}
		}
	}

}