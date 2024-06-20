package utils

// func WriteNodesToChannel(uniqueWorks *sql.Rows, outChannel chan types.Node) {
// 	counter := 0

// 	for uniqueWorks.Next() {
// 		nodeID := fmt.Sprintf("n%d", counter)
// 		uniqueWorks.Scan(&nodeID)
// 		outChannel <- types.Node{ID: nodeID, Data: types.Data{Key: "oa_id", Value: nodeID}}
// 		counter++
// 	}
// 	close(outChannel)
// }

// func WriteEdgesToChannel(rows *sql.Rows, nodeMap map[string]types.Node, outChannel chan types.Edge) {
// 	counter := 0
// 	for rows.Next() {
// 		var s, t string

// 		edgeID := fmt.Sprintf("e%d", counter)

// 		rows.Scan(&s, &t)

// 		/*
// 			Given an oa_id, lookup and store the Node ID attribute for source
// 			and target
// 		*/

// 		source := nodeMap[s].ID
// 		target := nodeMap[t].ID

// 		outChannel <- types.Edge{
// 			ID:     edgeID,
// 			Source: source,
// 			Target: target,
// 		}

// 		counter++
// 	}
// 	close(outChannel)
// }

// func BufferNodes(inChannel chan types.Node) map[string]types.Node {
// 	nodeMap := map[string]types.Node{}

// 	bar := progressbar.Default(-1, "Buffering nodes...")
// 	for node := range inChannel {
// 		// Store oa_id as the key and the Node object as the value
// 		nodeMap[node.Data.Value] = node

// 		bar.Add(1)
// 	}
// 	bar.Exit()

// 	return nodeMap
// }

// func BufferEdges(inChannel chan types.Edge) []types.Edge {
// 	var edges []types.Edge

// 	bar := progressbar.Default(-1, "Buffering edges...")
// 	for edge := range inChannel {
// 		edges = append(edges, edge)
// 		bar.Add(1)
// 	}
// 	bar.Exit()

// 	return edges
// }

// func MapToNodeSlice(nodeMap map[string]types.Node) []types.Node {
// 	nodes := []types.Node{}

// 	bar := progressbar.Default(int64(len(nodeMap)), "Creating a slice of Node objects...")
// 	for _, value := range nodeMap {
// 		nodes = append(nodes, value)
// 		bar.Add(1)
// 	}
// 	bar.Exit()

// 	return nodes
// }

// func CreateGraphML(nodes []types.Node, edges []types.Edge) types.GraphML {
// 	graph := types.Graph{
// 		ID:          "G",
// 		Edgedefault: "directed",
// 		Nodes:       nodes,
// 		Edges:       edges,
// 	}

// 	graphML := types.GraphML{
// 		Xmlns:  "http://graphml.graphdrawing.org/xmlns",
// 		Graphs: []types.Graph{graph},
// 	}

// 	return graphML
// }
