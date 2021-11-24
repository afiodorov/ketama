// Copyright 2016 Chao Wang <hit9@icloud.com>

// Package ketama implements a consistent hashing ring.
package ketama

import (
	"crypto/md5"
	"fmt"
	"sort"
)

// Node is the hashing ring node.
type Node struct {
	key    string
	data   interface{}
	weight uint
	hash   uint32
}

// NewNode creates a new Node.
func NewNode(key string, data interface{}, weight uint) *Node {
	return &Node{key: key, data: data, weight: weight}
}

// Key returns the Node key.
func (n *Node) Key() string {
	return n.key
}

// Data returns the Node data.
func (n *Node) Data() interface{} {
	return n.data
}

// Weight returns the Node weight.
func (n *Node) Weight() uint {
	return n.weight
}

// ByHash implements sort.Interface.
type ByHash []*Node

func (s ByHash) Len() int           { return len(s) }
func (s ByHash) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByHash) Less(i, j int) bool { return s[i].hash < s[j].hash }

// Ring is the ketama hashing ring.
type Ring struct {
	nodes []*Node
}

// alignHash returns hash value with aligment.
func alignHash(key string, align int) uint32 {
	b := md5.Sum([]byte(key))
	return ((uint32(b[3+align*4]&0xff) << 24) |
		(uint32(b[2+align*4]&0xff) << 16) |
		(uint32(b[1+align*4]&0xff) << 8) |
		(uint32(b[0+align*4] & 0xff)))
}

// NewRing creates a new Ring.
func NewRing(nodes []*Node) *Ring {
	// Create ring and init its nodes.
	r := &Ring{}
	length := 0
	for i := 0; i < len(nodes); i++ {
		length += int(nodes[i].weight) * 160
	}
	r.nodes = make([]*Node, length)
	// Init each ring node.
	k := 0
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		for j := 0; j < int(node.weight)*40; j++ {
			key := fmt.Sprintf("%s-%d", node.key, j)
			for n := 0; n < 4; n++ {
				r.nodes[k] = &Node{}
				r.nodes[k].key = node.key
				r.nodes[k].weight = node.weight
				r.nodes[k].data = node.data
				r.nodes[k].hash = alignHash(key, n)
				k++
			}
		}
	}
	sort.Sort(ByHash(r.nodes))
	return r
}

func (r *Ring) getByHash(hash uint32) int {
	if len(r.nodes) == 0 {
		return -1
	}

	return sort.Search(len(r.nodes), func(i int) bool { return r.nodes[i].hash > hash }) % len(r.nodes)
}

// Get node by key from ring.
// Returns nil if the ring is empty.
func (r *Ring) Get(key string) *Node {
	ind := r.getByHash(alignHash(key, 0))
	if ind == -1 {
		return nil
	}

	return r.nodes[ind]
}

// GetIgnoringFailed implements a failover when certain nodes are no longer responsive:
// those won't be used.
func (r *Ring) GetIgnoringFailed(key string, failedNodes map[string]struct{}) *Node {
	ind := r.getByHash(alignHash(key, 0))
	if ind == -1 {
		return nil
	}

	for i, j := ind, 0; j < len(r.nodes); i, j = (i+1)%len(r.nodes), j+1 {
		// keep going right until a responsive node is found
		if _, ok := failedNodes[r.nodes[i].key]; !ok {
			return r.nodes[i]
		}
	}

	// we have gone full circle => everything is failed
	return nil
}
