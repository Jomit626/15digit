//This file define varies of cost func that can be used
package main

// DepthPulsManhattan func
// Manhattan distance from current board to targetbaords +
// the depth of the node
func DepthPulsManhattan(b *board, depth uint) uint {
	return depth + uint(b.Distance())
}

// Manhattan func is the Manhattan distance from current board to targetbaords
func Manhattan(b *board, depth uint) uint {
	return uint(b.Distance())
}

func test(expectDepth uint) func(*board,uint) uint{
	return func(b *board, depth uint) uint{
		return depth*expectDepth/100 + uint(b.Distance())
	}
}

func ExpectedDepthPulsManhattan(expectDepth uint) func(*board,uint) uint{
	return func(b *board, depth uint) uint{
		cost := uint(b.Distance())
		if depth > expectDepth{
			cost = cost + depth - expectDepth
		}
		return cost 
	}
}


