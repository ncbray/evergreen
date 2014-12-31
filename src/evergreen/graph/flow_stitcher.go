package graph

type stitch struct {
	src EdgeID
	dst NodeID
}

// Assists creating a new graph based on the structure of an old graph.
// Maintains a map of edges in the old graph onto edges in the new graph.
// Defers connecting edges in the new graph until both the source and destination are known.
type EdgeStitcher struct {
	original   *Graph
	translated *Graph
	stitches   []stitch
}

func (stitcher *EdgeStitcher) setSrc(original EdgeID, srcID EdgeID) {
	current := stitcher.stitches[original].src
	if current != NoEdge {
		panic(current)
	}
	stitcher.stitches[original].src = srcID
	dstID := stitcher.stitches[original].dst
	if dstID != NoNode {
		stitcher.translated.Connect(srcID, dstID)
	}
}

func (stitcher *EdgeStitcher) setDst(original EdgeID, dstID NodeID) {
	current := stitcher.stitches[original].dst
	if current != NoNode {
		panic(current)
	}
	stitcher.stitches[original].dst = dstID
	srcID := stitcher.stitches[original].src
	if srcID != NoEdge {
		stitcher.translated.Connect(srcID, dstID)
	}
}

func (stitcher *EdgeStitcher) MapIncomingEdges(original NodeID, translated NodeID) {
	iter := EntryIterator(stitcher.original, original)
	for iter.HasNext() {
		_, edge := iter.GetNext()
		stitcher.setDst(edge, translated)
	}
}

func (stitcher *EdgeStitcher) MapEdge(original EdgeID, translated EdgeID) {
	stitcher.setSrc(original, translated)
}

func MakeEdgeStitcher(original *Graph, translated *Graph) *EdgeStitcher {
	numStitches := original.NumEdges()
	stitches := make([]stitch, numStitches)
	for i := 0; i < numStitches; i++ {
		stitches[i] = stitch{src: NoEdge, dst: NoNode}
	}
	return &EdgeStitcher{
		original:   original,
		translated: translated,
		stitches:   stitches,
	}
}
