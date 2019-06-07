package main

import (
	"testing"
	"fmt"
)

/* func TestRand(t *testing.T){
	b := Targetboard
	b.Rand(150)
	b.Print()
} */

func TestSearch(t *testing.T) {
	var b = Targetboard 
	b.Rand(100)
	b.Print()
	var config = configuration{
		origin_board: b,
		maxdepth: 1024,
		expectdepth: 999999,
		searchtime: 60,
	}
	var c context
	c.init(&config)
	path := Search(&c)
	print_path(path)
	println("Trying to shorten the path")
	done := make(chan struct{})
	for i:=90;i>0;i-=30{
		config.expectdepth = len(path) * i / 100
		var c context
		c.init(&config)
		go func(){
			path = Search(&c)
			println(len(path))
			done <- struct{}{}
		}()
	}

	for i:=0;i<3;i++{
		<-done
	}
}

func print_path(path []int8){
	for _,mv := range path {
		fmt.Print(move_description[mv])
	}
	println(len(path))
}