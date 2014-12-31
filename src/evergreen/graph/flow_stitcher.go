package graph

type stitch struct {
	src EdgeID
	dst NodeID
}

type FlowStitcher struct {
	original   *Graph
	translated *Graph
	stitches   []stitch
}

func (stitcher *FlowStitcher) setSrc(original EdgeID, srcID EdgeID) {
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

func (stitcher *FlowStitcher) setDst(original EdgeID, dstID NodeID) {
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

func (stitcher *FlowStitcher) SetHead(original NodeID, translated NodeID) {
	iter := EntryIterator(stitcher.original, original)
	for iter.HasNext() {
		_, edge := iter.GetNext()
		stitcher.setDst(edge, translated)
	}
}

func (stitcher *FlowStitcher) Internal(src NodeID, flow int, dst NodeID) {
	g := stitcher.translated
	e := g.IndexedExitEdge(src, flow)
	g.Connect(e, dst)
}

func (stitcher *FlowStitcher) SetEdge(srcID NodeID, srcFlow int, dstID NodeID, dstFlow int) {
	original := stitcher.original.IndexedExitEdge(srcID, srcFlow)
	translated := stitcher.translated.IndexedExitEdge(dstID, dstFlow)
	stitcher.setSrc(original, translated)
}

func (stitcher *FlowStitcher) NumExits(srcID NodeID) int {
	return stitcher.original.NumExits(srcID)
}

func (stitcher *FlowStitcher) GetExit(srcID NodeID, flow int) NodeID {
	return stitcher.original.GetExit(srcID, flow)
}

func MakeFlowStitcher(original *Graph, translated *Graph) *FlowStitcher {
	numStitches := original.NumEdges()
	stitches := make([]stitch, numStitches)
	for i := 0; i < numStitches; i++ {
		stitches[i] = stitch{src: NoEdge, dst: NoNode}
	}
	return &FlowStitcher{
		original:   original,
		translated: translated,
		stitches:   stitches,
	}
}
