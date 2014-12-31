package graph

type stitch struct {
	src  NodeID
	flow int
	dst  NodeID
}

func (stitch *stitch) SetSrc(srcID NodeID, flow int, g *Graph) {
	if stitch.src != NoNode {
		panic(stitch.src)
	}
	stitch.src = srcID
	stitch.flow = flow
	if stitch.dst != NoNode {
		stitch.link(g)
	}
}

func (stitch *stitch) SetDst(dstID NodeID, g *Graph) {
	if stitch.dst != NoNode {
		panic(stitch.dst)
	}
	stitch.dst = dstID
	if stitch.src != NoNode {
		stitch.link(g)
	}
}

func (stitch *stitch) link(g *Graph) {
	g.Connect(stitch.src, stitch.flow, stitch.dst)
}

type FlowStitcher struct {
	original   *Graph
	translated *Graph
	stitches   [][]*stitch
}

func (stitcher *FlowStitcher) getStitch(srcID NodeID, flow int) *stitch {
	s := stitcher.stitches[srcID]
	if s == nil {
		num := stitcher.original.NumExits(srcID)
		s = make([]*stitch, num)
		for i := 0; i < num; i++ {
			s[i] = &stitch{src: NoNode, dst: NoNode}
		}
		stitcher.stitches[srcID] = s
	}
	return s[flow]
}

func (stitcher *FlowStitcher) SetHead(srcID NodeID, dstID NodeID) {
	iter := EntryIterator(stitcher.original, srcID)
	for iter.HasNext() {
		prev, edge := iter.GetNext()
		stitch := stitcher.getStitch(prev, stitcher.original.EdgeFlow(edge))
		stitch.SetDst(dstID, stitcher.translated)
	}
}

func (stitcher *FlowStitcher) Internal(dstSrc NodeID, dstFlow int, dstDst NodeID) {
	stitcher.translated.Connect(dstSrc, dstFlow, dstDst)
}

func (stitcher *FlowStitcher) SetEdge(srcID NodeID, srcFlow int, dstID NodeID, dstFlow int) {
	stitch := stitcher.getStitch(srcID, srcFlow)
	stitch.SetSrc(dstID, dstFlow, stitcher.translated)
}

func (stitcher *FlowStitcher) NumExits(srcID NodeID) int {
	return stitcher.original.NumExits(srcID)
}

func (stitcher *FlowStitcher) GetExit(srcID NodeID, flow int) NodeID {
	return stitcher.original.GetExit(srcID, flow)
}

func MakeFlowStitcher(original *Graph, translated *Graph) *FlowStitcher {
	return &FlowStitcher{
		original:   original,
		translated: translated,
		stitches:   make([][]*stitch, original.NumNodes()),
	}
}
