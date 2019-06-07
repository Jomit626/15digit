package main

import (
	"testing"
	//"fmt"
)
func TestDistance(t *testing.T) {
	var mvs  = [...]int8 {
		MoveUp,
		MoveUp,
		MoveUp,
		MoveLeft,
		MoveLeft,
		MoveLeft,
	}
	var b1 = Targetboard
	b1.Moves(mvs[:])
	b1.Distance()
}

func BenchmarkDistance(b *testing.B){
	var mvs  = [...]int8 {
		MoveUp,
		MoveUp,
		MoveUp,
		MoveLeft,
		MoveLeft,
		MoveLeft,
	}
	var b1 = Targetboard
	b1.Moves(mvs[:])
	for n := 0; n < b.N; n++ {
		b1.Distance()
	}
}