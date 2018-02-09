package gographer

func Bellman_Ford(graph Graph, heuristic func(edge Edge) float64, source Node) (*Node, error) {

	// heur := make(map[uuid.UUID]float64, len(graph.Edges))
	// dist := make(map[uuid.UUID]Edge, len(graph.Nodes))
	// for _, e := range graph.Edges {
	// 	dist[e.Id] = MakeEdge(append([]*Node{&source}, e.Nodes...), math.MaxFloat64, Bidirectional, nil)
	// }
	// dist[source.Id] = MakeEdge([]*Node{&source}, 0, Bidirectional, nil)
	// for range graph.Nodes {
	// 	for _, e := range graph.Edges {
	// 		heur[e.Id] = heuristic(e)
	// 		newCost := heur[e.Id] + dist[e.Id].Distance
	// 		if newCost < dist[e.Id].Distance {
	// 			dist[e.Id] = MakeEdge(append([]*Node{&source}, e.Nodes...), newCost, Bidirectional, nil)
	// 		}
	// 	}
	// }
	// for _, e := range graph.Edges {
	// 	newCost := heur[e.Id] + dist[e.Id].Distance
	// 	if newCost < dist[e.Id].Distance {
	// 		return map[uuid.UUID]Edge{}, errors.New("Graph contains a negative-weight cycle")
	// 	}
	// }
	return nil, nil
}
