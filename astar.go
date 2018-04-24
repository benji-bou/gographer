package gographer

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

type NodePath struct {
	Node   *Node
	Parent *Node
	GScore float64
}

func (g Graph) AStar(source Node, goal Node) (NodePath, error) {
	closedSet := map[uuid.UUID]*Node{}
	openSet := map[uuid.UUID]*Node{source.Id: &source}

	startNode := NodePath{Node: &source, Parent: nil, GScore: 0}
	historic := map[uuid.UUID]NodePath{source.Id: startNode}
	for len(openSet) > 0 {
		n := &Node{}
		for _, tmp := range openSet {
			n = tmp
			break
		}
		if n.Id == goal.Id {
			return historic[goal.Id], nil
		}
		delete(openSet, n.Id)
		closedSet[n.Id] = n
		for _, l := range n.Neighbours {
			if _, isOk := closedSet[l.To.Id]; isOk == true {
				continue
			}
			if _, isOk := openSet[l.To.Id]; isOk == false {
				openSet[l.To.Id] = g.Nodes[l.To.Id]
			}
			historicNeighbour, isOk := historic[l.To.Id]
			tentative_gScore := historic[n.Id].GScore + l.Cost
			if isOk == false || tentative_gScore >= historicNeighbour.GScore {
				continue
			}
			historic[l.To.Id] = NodePath{Node: g.Nodes[l.To.Id], GScore: tentative_gScore, Parent: n}
		}
	}
	return NodePath{}, errors.New("No Path found")
}
