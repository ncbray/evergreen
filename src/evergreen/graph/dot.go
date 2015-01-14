package graph

import (
	"bytes"
	"fmt"
	"regexp"
)

func nodeDotID(node NodeID) string {
	return fmt.Sprintf("n%d", node)
}

type DotStyler interface {
	BlockLabel(node NodeID) (string, bool)
	NodeStyle(node NodeID) string
	EdgeStyle(src NodeID, edge EdgeID, dst NodeID) string
}

type edgePort struct {
	node NodeID
	port string
}

type dotDrawer struct {
	buf        *bytes.Buffer
	edgePorts  []edgePort
	fuseLinear bool
	uid        int
}

func (drawer *dotDrawer) getUID() int {
	temp := drawer.uid
	drawer.uid += 1
	return temp
}

func (drawer *dotDrawer) WriteString(message string) {
	drawer.buf.WriteString(message)
}

func (drawer *dotDrawer) WriteNode(nid NodeID, style string) {
	drawer.WriteString("  ")
	drawer.WriteString(nodeDotID(nid))
	drawer.WriteString("[")
	drawer.WriteString(style)
	drawer.WriteString("];\n")
}

func (drawer *dotDrawer) GetEdgePort(nid NodeID) string {
	port := drawer.edgePorts[nid]
	if port.port == "" {
		return nodeDotID(nid)
	} else {
		return nodeDotID(port.node) + ":" + port.port
	}
}

func (drawer *dotDrawer) IsSquashedEdge(src NodeID, dst NodeID) bool {
	sp := drawer.edgePorts[src]
	dp := drawer.edgePorts[dst]
	return sp.port != "" && dp.port != "" && sp.node == dp.node
}

func (drawer *dotDrawer) WriteEdge(src NodeID, dst NodeID, style string) {
	if drawer.IsSquashedEdge(src, dst) {
		return
	}
	drawer.WriteString("  ")
	drawer.WriteString(drawer.GetEdgePort(src))
	drawer.WriteString(" -> ")
	drawer.WriteString(drawer.GetEdgePort(dst))
	drawer.WriteString("[")
	drawer.WriteString(style)
	drawer.WriteString("];\n")
}

func drawNode(drawer *dotDrawer, node NodeID, styler DotStyler) {
	drawer.WriteNode(node, styler.NodeStyle(node))
}

func drawUnclusteredNodes(drawer *dotDrawer, order []NodeID, styler DotStyler) {
	nit := OrderedIterator(order)
	for nit.HasNext() {
		drawNode(drawer, nit.GetNext(), styler)
	}
}

var dotEscape = regexp.MustCompile("([\\\\\"\\[\\]<>{}|])")
var newline = regexp.MustCompile("\n")

func EscapeDotString(message string) string {
	return newline.ReplaceAllString(dotEscape.ReplaceAllString(message, "\\$1"), "\\l")
}

func dotString(message string) string {
	return fmt.Sprintf("\"%s\"", EscapeDotString(message))
}

func drawLinearNodes(drawer *dotDrawer, nodes []NodeID, styler DotStyler) {
	head := NoNode
	text := ""
	flush := func() {
		if head != NoNode {
			style := fmt.Sprintf("shape=record,label=\"{%s}\"", text)
			drawer.WriteNode(head, style)
			text = ""
			head = NoNode
		}
	}
	for _, n := range nodes {
		label, ok := styler.BlockLabel(n)
		if ok && drawer.fuseLinear {
			if head == NoNode {
				head = n
			} else {
				text += "|"
			}
			text += "<" + nodeDotID(n) + ">" + EscapeDotString(label)
			drawer.edgePorts[n] = edgePort{node: head, port: nodeDotID(n)}
		} else {
			flush()
			drawNode(drawer, n, styler)
		}
	}
	flush()
}

func drawCluster(drawer *dotDrawer, cluster Cluster, styler DotStyler) {
	switch cluster := cluster.(type) {
	case *ClusterLeaf:
		drawer.WriteString(fmt.Sprintf("subgraph cluster_%d {\n", drawer.getUID()))
		drawer.WriteString("  labeljust=l;\n")
		drawer.WriteString(fmt.Sprintf("  label=\"leaf %d\";\n", len(cluster.Nodes)))
		drawer.WriteString("  color=lightgrey;\n")
		drawLinearNodes(drawer, cluster.Nodes, styler)
		drawer.WriteString("}\n")
	case *ClusterLinear:
		drawer.WriteString(fmt.Sprintf("subgraph cluster_%d {\n", drawer.getUID()))
		drawer.WriteString("  labeljust=l;\n")
		drawer.WriteString(fmt.Sprintf("  label=\"linear %d\";\n", len(cluster.Clusters)))
		drawer.WriteString("  color=lightgrey;\n")
		for _, c := range cluster.Clusters {
			drawCluster(drawer, c, styler)
		}
		drawer.WriteString("}\n")
	case *ClusterSwitch:
		drawer.WriteString(fmt.Sprintf("subgraph cluster_%d {\n", drawer.getUID()))
		drawer.WriteString("  labeljust=l;\n")
		drawer.WriteString(fmt.Sprintf("  label=\"switch %d\";\n", len(cluster.Children)))
		drawer.WriteString("  color=lightgrey;\n")
		for _, c := range cluster.Children {
			drawCluster(drawer, c, styler)
		}
		drawer.WriteString("}\n")
	case *ClusterLoop:
		drawer.WriteString(fmt.Sprintf("subgraph cluster_%d {\n", drawer.getUID()))
		drawer.WriteString("  labeljust=l;\n")
		drawer.WriteString("  label=loop;\n")
		drawer.WriteString("  color=lightgrey;\n")
		drawCluster(drawer, cluster.Body, styler)
		drawer.WriteString("}\n")
	default:
		panic(cluster)
	}
}

func drawClusteredNodes(drawer *dotDrawer, g *Graph, styler DotStyler) {
	cluster := makeCluster(g)
	drawCluster(drawer, cluster, styler)
}

func GraphToDot(g *Graph, styler DotStyler) string {
	order, index := ReversePostorder(g)

	var idoms []NodeID
	visualize_idoms := false
	if visualize_idoms {
		idoms = FindDominators(g, order, index)
	}

	drawer := &dotDrawer{buf: &bytes.Buffer{}, edgePorts: make([]edgePort, g.NumNodes()), fuseLinear: true}

	drawer.WriteString("digraph G {\n")
	drawer.WriteString("  nslimit = 3;\n") // Make big graphs render faster.

	//drawUnclusteredNodes(drawer, order, styler)
	drawClusteredNodes(drawer, g, styler)

	// Draw edges.
	nit := OrderedIterator(order)
	for nit.HasNext() {
		node := nit.GetNext()
		eit := g.ExitIterator(node)
		for eit.HasNext() {
			edge, dst := eit.GetNext()
			style := styler.EdgeStyle(node, edge, dst)
			if index[node] >= index[dst] {
				style += ",weight=0"
			}
			drawer.WriteEdge(node, dst, style)
		}
	}
	if visualize_idoms {
		nit := OrderedIterator(order)
		for nit.HasNext() {
			src := nit.GetNext()
			dst := idoms[src]
			if src != dst {
				drawer.WriteEdge(src, dst, "style=dotted")
			}
		}
	}
	drawer.WriteString("}\n")
	return drawer.buf.String()
}
