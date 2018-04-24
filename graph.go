package gographer

import (
	"encoding/json"
	"errors"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrStopIterate     = errors.New("stopping iteration")
	ErrFinishedIterate = errors.New("iteration finished")
)

type ErrorStack []error

func (errs ErrorStack) Error() string {
	res := ""
	for _, err := range errs {
		res += err.Error() + "\n"
	}
	return res
}

type Direction string

const (
	Unidirectional Direction = "unidirectional"
	Bidirectional  Direction = "bidirectional"
)

type Link struct {
	Id        uuid.UUID   `json:"id"`
	Direction Direction   `json:"direction"`
	Cost      float64     `json:"cost"`
	From      *Node       `json:"-"`
	To        *Node       `json:"-"`
	Value     interface{} `json:"value"`
}

func MakeLink(direction Direction, cost float64, from, to *Node, value interface{}) Link {
	return Link{From: from, To: to, Direction: direction, Cost: cost, Value: value, Id: uuid.NewV5(uuid.Or(from.Id, to.Id), from.Id.String()+to.Id.String())}
}

func (l Link) MarshalJSON() ([]byte, error) {
	// type MarshallerLink = Link
	type MarshallerLink Link
	marsahlLink := struct {
		MarshallerLink
		FromId uuid.UUID `json:"from"`
		ToId   uuid.UUID `json:"to"`
	}{MarshallerLink: MarshallerLink(l), FromId: l.From.Id, ToId: l.To.Id}
	return json.Marshal(marsahlLink)
}

type Node struct {
	Id         uuid.UUID          `json:"id"`
	Neighbours map[uuid.UUID]Link `json:"neighbours"`
	Value      interface{}        `json:"value"`
}

func NodeCopy(n Node) *Node {
	return &Node{Id: n.Id, Neighbours: map[uuid.UUID]Link{}, Value: n.Value}
}

func NodesContain(ns []*Node, node Node) bool {
	for _, n := range ns {
		if n.Id == node.Id {
			return true
		}
	}
	return false
}

func MakeNode(val interface{}) Node {
	return Node{Id: uuid.NewV4(), Value: val, Neighbours: map[uuid.UUID]Link{}}
}

func MakeNodeId(id uuid.UUID, val interface{}) Node {
	return Node{Id: id, Value: val, Neighbours: map[uuid.UUID]Link{}}
}

func (n Node) Iterate(depth int, cb func(node Node) error) error {
	stack := []Node{n}
	closedList := map[uuid.UUID]Node{n.Id: n}
	errorsStack := ErrorStack{}
	for len(stack) > 0 {
		newNode := stack[0]
		err := cb(newNode)
		if err == ErrStopIterate {
			return err
		} else if err != nil {
			errorsStack = append(errorsStack, err)
		}
		closedList[newNode.Id] = newNode
		stack = append(stack[:0], stack[1:]...)
		if depth > 0 {
			for _, linked := range newNode.Neighbours {
				var subNode *Node
				if linked.From.Id != newNode.Id {
					subNode = linked.From
				} else {
					subNode = linked.To
				}
				if _, isOK := closedList[subNode.Id]; isOK == false {
					stack = append(stack, *subNode)
				}
			}
		}
		depth--
	}
	return errorsStack
}

func (n Node) IsLinkedToDepth(depth int, nodeId uuid.UUID) bool {
	err := n.Iterate(depth, func(node Node) error {
		if node.Id == nodeId {
			return ErrStopIterate
		}
		return nil
	})
	return err == ErrStopIterate
}

func (n Node) IsLinkedTo(nodeId uuid.UUID) bool {
	for _, l := range n.Neighbours {
		if l.From.Id == nodeId || l.To.Id == nodeId {

			return true
		}
	}
	return false
}

func (n *Node) AddNeighbour(node *Node, cost float64, direction Direction, value interface{}) {
	l := MakeLink(direction, cost, n, node, value)
	if _, isOk := n.Neighbours[l.Id]; isOk == false {
		n.Neighbours[l.Id] = l
	}
	if direction == Bidirectional {
		if _, isOk := node.Neighbours[l.Id]; isOk == false {
			node.Neighbours[l.Id] = l
		}
	}
}

type Edge struct {
	Id        uuid.UUID
	Distance  float64
	Nodes     []*Node
	Direction Direction
	Value     interface{}
}

func (e *Edge) Append(ne Edge, distance float64, direction Direction, value interface{}) {
	e.Distance += ne.Distance + distance
	if len(ne.Nodes) > 0 {
		if len(e.Nodes) > 0 {
			e.Nodes[len(e.Nodes)-1].AddNeighbour(ne.Nodes[0], distance, direction, value)
		}
		e.Nodes = append(e.Nodes, ne.Nodes...)
	}
}

func (e *Edge) AppendNode(n *Node, distance float64, direction Direction, value interface{}) {
	e.Distance += distance
	if len(e.Nodes) > 0 {
		e.Nodes[len(e.Nodes)-1].AddNeighbour(n, distance, direction, value)
	}
	e.Nodes = append(e.Nodes, n)

}

func MakeEdgeId(id uuid.UUID, nodes []*Node, distance float64, direction Direction, value interface{}) Edge {
	return Edge{Id: id, Distance: distance, Nodes: nodes, Direction: direction, Value: value}
}

func MakeEdge(nodes []*Node, distance float64, direction Direction, value interface{}) Edge {
	return Edge{Id: uuid.NewV4(), Distance: distance, Nodes: nodes, Direction: direction, Value: value}
}

type Relation struct {
	Edeges []Edge
	Value  interface{}
}

type Graph struct {
	Nodes     map[uuid.UUID]*Node
	Edges     []Edge
	Relations []Relation
}

func (g *Graph) AddEdge(edge Edge) {

	if len(edge.Nodes) > 0 {
		n := edge.Nodes[0]
		g.Nodes[n.Id] = n
	}
	for i := 1; i < len(edge.Nodes); i++ {
		n := edge.Nodes[i]
		prevN := edge.Nodes[i-1]
		n.AddNeighbour(prevN, 1, Bidirectional, nil)
		prevN.AddNeighbour(n, 1, Bidirectional, nil)
		g.Nodes[n.Id] = n
	}
	g.Edges = append(g.Edges, edge)
}

func (g *Graph) AddNode(node *Node) {
	g.Nodes[node.Id] = node
}

func NewGraph() *Graph {
	return &Graph{Nodes: make(map[uuid.UUID]*Node), Edges: make([]Edge, 0)}
}
