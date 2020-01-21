package main

import (
	"fmt"

	"github.com/vvatanabe/weighted-roundrobin"
)

func main() {

	rr := weighted.New([]*weighted.Node{
		{
			Value:  "apple",
			Weight: 2,
		},
		{
			Value:  "banana",
			Weight: 4,
		},
		{
			Value:  "grape",
			Weight: 4,
		},
		{
			Value:  "orange",
			Weight: 18,
		},
	})

	result := make(map[string]uint64)
	for i := 0; i < 100; i++ {
		n := rr.GetNode()
		name := n.Value.(string)
		result[name] = result[name] + 1
	}

	fmt.Println(result)
}
