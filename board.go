//This file defines what a board is and the operations we can do to it.
package main 

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Target = 18364758544493064720	//	the board of the target
	Blank = 0xf	//15 representing the blank 

	// Move operations
	// new pos = pos + movement
	MoveUp 		= -4	
	MoveDown 	= +4
	MoveLeft 	= -1
	MoveRight 	= +1
	MoveNone	=  0

)
/*
 	Use 16 bytes to record a borad
 	index table
 	0  1  2  3
	4  5  6  7
	8  9  10 11
	12 13 14 15
*/
type board struct {
	data uint64
	blank int8 
}

// Targetboard The board in which the number of position i is i
var Targetboard = board{18364758544493064720,15}

// MoveDescription For output description
var MoveDescription = map[int8]string{
	MoveUp		:"UP ",
	MoveDown	:"Down ",
	MoveLeft 	:"Left ",
	MoveRight	:"Right ",
	MoveNone	:"START ",
}

// PossibleMoves Given the position of the blank and the previous move([blank pos] [preivous mv + 4]), this table tells the possible moves to preform		
var PossibleMoves = [16][9][]int8 {
	//  Previous move
	//	{{/*MoveUp*/},nil,nil,				 {/*MoveLeft*/},			{/*MoveNone*/},							{/*MoveRight*/},nil,nil,			{/*MoveDown*/},},
		{{MoveRight},nil,nil,				 {MoveDown},				{MoveRight,MoveDown},					nil,nil,nil,						nil,},							// Pos 0
		{{MoveLeft,MoveRight},nil,nil,		 {MoveLeft,MoveDown},		{MoveLeft,MoveRight,MoveDown},			{MoveDown,MoveRight},nil,nil,		nil,},							// Pos 1
		{{MoveLeft,MoveRight},nil,nil,		 {MoveLeft,MoveDown},		{MoveLeft,MoveRight,MoveDown},			{MoveDown,MoveRight},nil,nil,		nil,},							// Pos 2
		{{MoveLeft},nil,nil,				 nil,						{MoveLeft,MoveDown},					{MoveDown},nil,nil,					nil,},							// Pos 3
		{{MoveUp,MoveRight},nil,nil,		 {MoveUp,MoveDown},			{MoveUp,MoveRight,MoveDown},			nil,nil,nil,						{MoveDown,MoveRight},},			// Pos 4
		{{MoveUp,MoveLeft,MoveRight},nil,nil,{MoveUp,MoveLeft,MoveDown},{MoveUp,MoveLeft,MoveRight,MoveDown},	{MoveUp,MoveDown,MoveRight},nil,nil,{MoveLeft,MoveDown,MoveRight},},// Pos 5
		{{MoveUp,MoveLeft,MoveRight},nil,nil,{MoveUp,MoveLeft,MoveDown},{MoveUp,MoveLeft,MoveRight,MoveDown},	{MoveUp,MoveDown,MoveRight},nil,nil,{MoveLeft,MoveDown,MoveRight},},// Pos 6
		{{MoveUp,MoveLeft},nil,nil,			 nil,						{MoveUp,MoveLeft,MoveDown},				{MoveUp,MoveDown},nil,nil,			{MoveLeft,MoveDown},},			// Pos 7
		{{MoveUp,MoveRight},nil,nil,		 {MoveUp,MoveDown},			{MoveUp,MoveRight,MoveDown},			nil,nil,nil,						{MoveDown,MoveRight},},			// Pos 8
		{{MoveUp,MoveLeft,MoveRight},nil,nil,{MoveUp,MoveLeft,MoveDown},{MoveUp,MoveLeft,MoveRight,MoveDown},	{MoveUp,MoveDown,MoveRight},nil,nil,{MoveLeft,MoveDown,MoveRight},},// Pos 9
		{{MoveUp,MoveLeft,MoveRight},nil,nil,{MoveUp,MoveLeft,MoveDown},{MoveUp,MoveLeft,MoveRight,MoveDown},	{MoveUp,MoveDown,MoveRight},nil,nil,{MoveLeft,MoveDown,MoveRight},},// Pos 10
		{{MoveUp,MoveLeft},nil,nil,			 nil,						{MoveUp,MoveLeft,MoveDown},				{MoveUp,MoveDown},nil,nil,			{MoveLeft,MoveDown},},			// Pos 11
		{nil,nil,nil,						 {MoveUp},					{MoveUp,MoveRight},						nil,nil,nil,						{MoveRight},},					// Pos 12
		{nil,nil,nil,						 {MoveUp,MoveLeft},			{MoveUp,MoveLeft,MoveRight},			{MoveUp,MoveRight},nil,nil,			{MoveLeft,MoveRight},},			// Pos 13
		{nil,nil,nil,						 {MoveUp,MoveLeft},			{MoveUp,MoveLeft,MoveRight},			{MoveUp,MoveRight},nil,nil,			{MoveLeft,MoveRight},},			// Pos 14
		{nil,nil,nil,						 nil,						{MoveUp,MoveLeft},						{MoveUp},nil,nil,					{MoveLeft},},					// Pos 15
}

// Move the blank of a board
// To undo the movement, call with `-mv`
func (b *board) Move(mv int8){
	blank :=  b.blank
	nblank := b.blank + mv
	b.blank = nblank

	num := b.get(nblank)	// move the number
	b.set(blank,num)

	b.set(nblank,Blank)		// clear out blank
}

// Move the blank of a board 
func (b *board) Moves(mvs []int8){
	for _,mv := range mvs {
		b.Move(mv)
	}
}

// Get the number at the given position
func (b *board) get(idx int8) uint64 {
	return (b.data >> (uint8(idx) << 2)) & 0xf
}

// Set the number at the given position
func (b *board) set(idx int8, num uint64){
	var offset = uint8(idx) << 2
	var mask uint64 = 0xf << offset
	new := b.data & ^mask
	new |= num << offset
	b.data = new
}

// Print the board to stdout in format
func (b *board) Print(){
	data := b.data
	fmt.Printf("%d %d\n",data,b.blank)
	for i:=0; i< 16; i++{
		if i==int(b.blank) {
			fmt.Printf("%2d ",0)
		} else {
			fmt.Printf("%2d ",data&0xf + 1)
		}
		data >>= 4
		if i&0x3 == 0x3 {
			fmt.Printf("\n")
		}
	}
}

// Caculate the distance to Targetboard
// 		Return the distance
// 	Method to caculate the distance of one number:
//	tmp = | number -  the position index of that number |
//	distance = tmp % 4 + [tmp / 4]
func (b *board) Distance() uint32 {
	const (
		signmask = uint32(0x80000000)
	)
	var s1,s2,s3,s4 uint32
	
	data := b.data

	// 4 X 4 loop unrolling
	for i:=uint32(0);i<16;i+=4{
		var n1,n2,n3,n4 uint32
		var t1,t2,t3,t4 uint32
		n1 = uint32(data&0xf)
		n2 = i
		t1 = n1 - n2
		if (t1&signmask != 0){	// if t is negative
			t1 = ^t1 + 1
		}
		data >>= 4
		
		n3 = uint32(data&0xf)
		n4 = i + 1
		t2 = n3 - n4
		if (t2&signmask != 0){
			t2 = ^t2 + 1
		}
		data >>= 4
		
		n1 = uint32(data&0xf)
		n2 = i + 2
		t3 = n1 - n2
		if (t3&signmask != 0){
			t3 = ^t3 + 1
		}
		data >>= 4
		
		n3 = uint32(data&0xf)
		n4 = i + 3
		t4 = n3 - n4
		if (t4&signmask != 0){
			t4 = ^t4 + 1
		}
		data >>= 4
		
		s1 += t1&0x3 + (t1>>2) // s += t%4 + t/4
		s2 += t2&0x3 + (t2>>2)
		s3 += t3&0x3 + (t3>>2)
		s4 += t4&0x3 + (t4>>2)
	}
	t := 15 - uint32(b.blank)	// subtract the Blank distance
	t = t&0x3 + (t>>2) 
	return (s1 + s2) + (s3 + s4) - t
}

// Make a board move n times randomly
func (b *board) Rand(n int){
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	mv := int8(MoveNone)
	for i:=0;i<n;i++{
		pmv := PossibleMoves[b.blank][mv + 4]
		mv = pmv[r.Intn(len(pmv))]
		b.Move(mv)
	}
}