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

func TesdtSearch(t *testing.T) {
	var b = Targetboard 
	b.Rand(2000)
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

func TestTask(t *testing.T){
	var b = Targetboard
	b.set(0,10)
	b.set(1,8)
	b.set(2,3)
	b.set(3,14)
	b.set(4,0)
	b.set(5,2)
	b.set(6,15)
	b.set(7,11)
	b.set(8,6)
	b.set(9,4)
	b.set(10,7)
	b.set(11,5)
	b.set(12,12)
	b.set(13,1)
	b.set(14,9)
	b.set(15,13)
	b.blank = 6
	b.Print()
	var config = configuration{
		origin_board: b,
		maxdepth: 1024,
		expectdepth: 150,
		searchtime: 60,
	}
	var c context
	c.init(&config)
	path := Search(&c)
	print_path(path)
}

func print_path(path []int8){
	for _,mv := range path {
		fmt.Print(move_description[mv])
	}
	println(len(path))
}