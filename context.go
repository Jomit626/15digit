//
package main 

import (
	_ "fmt"
	"container/list"
	"time"
)
	
/* 
 * Context is a data structure that stores the information needed in a searching routine.
 * With this, it is easier to perform multiple search at the same time.
 * The pending nodes are group by their cost to the TargetBoard, stored in different lists, which avoids
 * sorting.(Trade time with space)
 */
type context struct {
	nNode uint32	// number of nodes searched
	//config configuration	// The configuration of this context
	searchtime int		//
	costFunc func(*board,uint) uint	//Comstom costFunc(board,depth)
	pending []*list.List	// the pending nodes of different cost
	done 	[56]*list.List	// the nodes that have been searched, grouped by their Manhattan distance to the Targetboard
}

/* 
 * Node of the search tree.
 */
type node struct {
	pre *node	// its parent node,used to determine path
	depth uint
	b board		// the current board
	cost uint	// the cost of the board
	mv int8		// previous move
}

// Initialise context with a config
func (c *context) init(config *configuration){
	//c.config = *config	// copy config
	c.searchtime = config.searchtime
	c.costFunc = config.costFunc
	// Init the lists
	c.pending = make([]*list.List,config.maxCost,config.maxCost)

	for i:=0;i<len(c.pending);i++ {
		c.pending[i] = list.New()
	}
	for i:=0;i<len(c.done);i++ {
		c.done[i] = list.New()
	}
	
	// insert origin board to pending list
	c.insert(config.originBoard.formnode(nil,MoveNone))
}

// Search routine
// 	Return slice of path if successful, nil otherwise
func Search(c *context) []int8 {
	timeout := time.After(time.Second * time.Duration(c.searchtime))

	n := c.evolve()
	for n == nil{
		n = c.evolve()
		
		select {	// if out of time, quit search
		case <-timeout:
			return nil
		default:
		}

	}

	return n.path()
}

// Evaluate a node with lowest cost in the pending list
// 		Return pointer to a node if Targetboard is reached, nil otherwise.
func (c *context) evolve() *node{
	var n *node
	var cost int
	
	// find one with lowest cost
	for cost=0;cost<len(c.pending);cost++{
		l := c.pending[cost]
		if l.Len() > 0 {
			n = l.Remove(l.Front()).(*node)	// getting form front and putting to the front to ensure the search path go deep first
			break
		}	
	}
	//println(n.depth)
	// you will always find one node here before you get the targetboard
	/* if (n == nil){
		return nil
	} */
	b := n.b
	if b.data == Target {
		return n	// found it
	}
	for _,mv := range PossibleMoves[b.blank][n.mv + 4] {
		new := n.move(mv)	// insert the new node
		new.cost = c.costFunc(&new.b,new.depth)
		c.insert(new)
	}

	done := c.done[b.Distance()]
	done.PushFront(n)
	return nil
}


// Insert a node to pending list
func (c *context) insert(n *node) {
/* 	b := n.b
	data := b.data

	distance := int(b.Distance())
	done := c.done[distance]
	// Todo: make a new func to calcute the cost
	if tmp :=(n.depth - c.config.expectdepth);tmp > 0 {
		distance = distance + tmp
	} */

	data := n.b.data	// the board representation
	
	l := c.pending[n.cost]	// the list in which the node should be stored
	done := c.done[n.b.Distance()]

	// Multithread search for a same board in the done list
	hasFound := false
	middle := done.Len()/2;
	found := make(chan *node)
		// From front to back
	go func() {
		for e,n := done.Front(),0; n < middle; e,n = e.Next(),n+1{
			if(hasFound) {
				break
			} else {
				node := e.Value.(*node)
				if node.b.data == data {	// if their borads are same
					found <- node
					return
				}
			}
		}
		found <- nil
	}()
		// From back to front
	go func() {
		for e,n := done.Back(),done.Len(); n > middle; e,n = e.Prev(),n-1{
			if(hasFound) {
				break
			} else {
				node := e.Value.(*node)
				if node.b.data == data {	// if their borads are same
					found <- node
					return
				}
			}
		}
		found <- nil
	}()

	var t *node
	for n:=0;n<2;n++{	// wait for result
		if recv := <- found; recv != nil{
			t = recv
			hasFound = true
		}
	}
	close(found)

	if t == nil{	// not found, insert new node
		l.PushFront(n) // getting form front and putting to the front to ensure the search path go deep first
		c.nNode++
	} else {	// found, if new node's depth is smaller, insert it
		if t.depth > n.depth{
			t.pre = n.pre
			t.depth = n.depth
			t.cost = n.cost
			t.mv = n.mv
		}
	}
}

// Form a node with depth of 1 form a boad.
func (b *board) formnode(pre *node ,mv int8) *node {
	var n = &node{
		b: 		*b,
		pre:	pre,
		mv:		mv,
		depth:	1,
		cost:	0,
	}
	n.b.Move(mv)
	return n
}


// Form a new node form existing one and perform operation mv to it.
func (n *node) move(mv int8) *node{
	var nc = *n	//nc new node

	nc.mv = mv
	nc.pre = n
	nc.depth++

	nc.b.Move(mv)
	return &nc
}

// Get path from root to node n.
func (n *node) path() []int8 {
	path := make([]int8,0,n.depth)

	for n!=nil {
		path = append(path,n.mv)
		n = n.pre
	}

	// Reverse
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }
	return path
}

