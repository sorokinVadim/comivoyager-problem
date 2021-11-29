package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Size    = 3
	MaxDist = 100
)

func makeMatrix(chance int) (matrix [][]int){
	for i := 0; i < Size; i++{
		matrix = append(matrix, make([]int, 0))
		for j := 0; j < Size; j++{
			matrix[i] = append(matrix[i], -1)
		}
	}
	for i := 0 ; i < Size; i++{
		for j := 0 ; j < Size; j++{
			randNum := rand.Int() % (MaxDist * chance)
			if randNum < MaxDist {
				matrix[i][j], matrix[j][i] = randNum, randNum
			} else {
				matrix[i][j], matrix[j][i] = 0, 0
			}
		}
	}
	for i := 0; i < Size; i++ {
		matrix[i][i] = 0
	}
	return matrix
}

type Node struct {
	name string
	neighbors map[*Node]int
}

func NewNode(name string) *Node{
	return &Node{name, make(map[*Node]int, 0)}
}

type Graph []*Node

func NewGraph() (nodes Graph) {
	for i := 0; i < Size; i++ {
		nodes = append(nodes, NewNode(strconv.Itoa(i+1)))
	}
	return
}

func (nodes Graph) InitGraph(matrix [][]int) Graph{
	for node, data := range matrix{
		for neighbour, dist := range data{
			if dist != 0 {
				nodes[node].neighbors[nodes[neighbour]] = dist
			}
		}
	}
	return nodes
}


func (graph Graph) String() string{
	result := ""
	for _, node := range graph{
		result += fmt.Sprintf("%s := [", node.name)
		for neigbour, dist := range node.neighbors{
			result += fmt.Sprintf("%s: %v, ", neigbour.name, dist)
		}
		result += fmt.Sprint("]\n")
	}
	return result
}

type Result struct {
	dist int
	path []*Node
}

func (result Result) String() string{
	names := make([]string, 0)
	for _, p := range result.path{
		names = append(names, p.name)
	}
	return fmt.Sprintf("Distance: %v\nPath: %s", result.dist, strings.Join(names, ", "))
}

func contain(graph []*Node, node *Node) bool{
	for _, n := range graph{
		if n.name == node.name{
			return true
		}
	}
	return false
}

func (graph Graph) FindLowerDistHungry() (result Result) {
	chanResult := make(chan Result, Size)
	defer close(chanResult)
	for i := 0; i < Size; i++ {
		go func(node *Node) {
			r := new(Result)
			r.path = append(r.path, node)
			for {
				lowerDist := MaxDist + 1
				var lowerNode *Node
				for neighbour, dist := range r.path[len(r.path) - 1].neighbors{
					if lowerDist > dist && !contain(r.path, neighbour){
						lowerDist, lowerNode = dist, neighbour
					}
				}
				if lowerNode == nil{
					break
				}
				r.dist += lowerDist
				r.path = append(r.path, lowerNode)
			}
			chanResult <- *r
		}(graph[i])
	}
	result = <-chanResult
	for i := 0; i < Size- 1; i++{
		r := <-chanResult
		if r.dist < result.dist{
			result = r
		}
	}
	return result
}

// Factorial
func f(n int) int {
	for i := n - 1; i > 0; i--{
		n *= i
	}
	return n
}

func alg(r Result, cRes chan Result) {
	if len(r.path) == Size {
		time.Sleep(time.Second * 2)
		cRes <- r
		return
	}
	for neighbour, dist := range r.path[len(r.path) - 1].neighbors{
		if contain(r.path, neighbour){
			continue
		}
		newR := r
		newR.dist += dist
		newR.path = append(newR.path, neighbour)
		go alg(newR, cRes)
	}
}

func (graph Graph) FindLowerDistSlow() (result Result) {
	chanResult := make(chan Result)
	defer close(chanResult)
	for _, n := range graph{
		r := Result{0, []*Node{n}}
		go alg(r, chanResult)
	}
	result = <-chanResult
	var r Result
	for i := 0; i < f(Size) - 1; i++ {
		r = <-chanResult
		if r.dist < result.dist{
			result = r
		}
	}
	return
}

func main() {
	seed := time.Now().Unix()
	rand.Seed(seed)
	//fmt.Println(seed)

	matrix := makeMatrix(1) // 0.5 -> 2, 0.25 -> 4, 1/chance = arg
	graph := NewGraph().InitGraph(matrix)
	start := time.Now().UnixMicro()
	fmt.Println(graph.FindLowerDistSlow())
	fmt.Printf("Slow alghoritm: %v ms\n", time.Now().UnixMicro() - start)
	start = time.Now().UnixMicro()
	fmt.Println(graph.FindLowerDistHungry())
	fmt.Printf("Hungry alghoritm: %v ms\n", time.Now().UnixMicro() - start)
}
