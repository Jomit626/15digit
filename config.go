package main

// Configuration is used to form context to perform search routine.
type configuration struct {
	originBoard board	// ...
	maxCost int 		// The max cost of a board can have
	//expectdepth int		// The answer path is expect to be less then this number. It is used
	//					// to avoid the A* method behaving like BFS.
	//					// Setting it to 0 means no expect depth
	searchtime int		// After searchtime seconds, the search will be forced to stop
	costFunc func(*board,uint) uint	//Comstom costFunc(board,depth)
	
}
