/* 
 * 
 */
package main 

import (
	_ "fmt"
	"container/list"
)
	
/* 
 * Context is a data structure that stores the information needed in a searching routine.
 * With this, it is easier to perform multiple search at same time.
 * The pending nodes are the nodes that have distance of 1 to a node in the done list.
 * They are group by their cost to the TargetBoard, stored in different lists, which avoids
 * sorting.(Trade time with space)
 * 
 */
type context struct {
	nnode uint32	// number of nodes searched
	config configuration	// The configuration of this context
	pending [56 + 1024]*list.List	// the pending nodes of different distances
									//	56 (Longest distance) + 1024 (Max depth)
	done 	[56]*list.List			// the nodes that have been searched
}

/* 
 * Node of the search tree.
 */
type node struct {
	pre *node	// used to determine path
	depth int
	b board			// the current board
	mv int8			// previous move
}

// Initialise context with config
func (c *context) init(config *configuration){
	c.config = *config	// copy config
	// Init the lists
	// TODO: set maxdepth 
	for i:=0;i<len(c.pending);i++ {
		c.pending[i] = list.New()
	}
	for i:=0;i<len(c.done);i++ {
		c.done[i] = list.New()
	}
	
	// insert origin board to pending list
	c.insert(config.origin_board.formnode(nil,MoveNone))
}

// Search routine
// 	Return slice of path if successful, nil otherwise
func Search(c *context) []int8 {
	// TODO: Time out
	n := c.evolve()
	for n == nil{
		n = c.evolve()
	}

	return n.path()
}

// Evaluate a node with lowest cost in the pending list
// 		Return pointer to a node if Targetboard is reached, nil otherwise.
func (c *context) evolve() *node{
	var n *node = nil
	var d int

	// find one with shortest distance
	for d=0;d<len(c.pending);d++{
		l := c.pending[d]
		if l.Len() > 0 {
			n = l.Remove(l.Front()).(*node)	// getting form front and putting to the front to ensure the path go deep first
			break
		}	
	}
	b := n.b
	t := int(b.Distance())
	if t == 0 {
		return n	// found it
	}
	for _,mv := range possible_moves[b.blank][n.mv + 4] {
		new := n.move(mv)	// insert the new node
		c.insert(new)
	}

	done := c.done[t]
	done.PushFront(n)
	return nil
}


// Insert a node to pending list
func (c *context) insert(n *node) {
	b := n.b
	data := b.data

	t := int(b.Distance())
	done := c.done[t]
	
	if tmp :=(n.depth - c.config.expectdepth);tmp > 0 {
		t = t + tmp
	}
	l := c.pending[t]
	
	for e := done.Front(); e != nil; e = e.Next() {	// TODO: MAKE IT A CON
		node := e.Value.(*node)
		if node.b.data == data {
			if node.depth <= n.depth { 	// this means that there is a shorter path to the board n stores,
				return					// so no need for n;
			} else {	// Updata depth and ancestor
				node.pre = n.pre
				node.depth = n.depth
				node.mv = n.mv
				return
			}
		}
	}

	l.PushFront(n)
	c.nnode++
}

// Form a node with depth of 1 form a boad.
func (b *board) formnode(pre *node ,mv int8) *node {
	var n node
	n.b = *b
	n.pre = pre
	n.mv = mv
	n.b.Move(mv)
	n.depth = 1
	return &n
}


// Form a new node form existing one and perform operation mv to it.
func (n *node) move(mv int8) *node{
	var nc node = *n

	nc.mv = mv
	nc.pre = n
	nc.depth++

	nc.b.Move(mv)
	return &nc
}

// Get path from root to node n.
func (n *node) path() []int8 {
	path := make([]int8,0,128)	//TODO: Make it a constant
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
