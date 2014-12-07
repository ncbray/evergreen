package framework

func FindLines(stream []rune) []int {
	lines := []int{0}
	for i, r := range stream {
		if r == '\n' {
			lines = append(lines, i+1)
		}
	}
	return lines
}

func GetLine(stream []rune, lines []int, line int) string {
	start := lines[line]
	end := len(stream)
	if line+1 < len(lines) {
		end = lines[line+1]
	}
	// HACK trim newline
	for end > start && (stream[end-1] == '\n' || stream[end-1] == '\r') {
		end -= 1
	}
	return string(stream[start:end])
}

func GetLocation(stream []rune, lines []int, pos int) (int, int, string) {
	// Stupid linear search
	var line int
	// If we don't find it, it must be on the last line.
	for line = 0; line < len(lines)-1; line++ {
		if pos >= lines[line] && pos < lines[line+1] {
			break
		}
	}
	col := pos - lines[line]
	return line, col, GetLine(stream, lines, line)
}

type LocationProvider interface {
	AddFile(filename string, stream []rune) int
	GetLocationInfo(pos int) (filename string, line int, col int, text string)
}

type fileInfo struct {
	Offset   int
	Filename string
	Stream   []rune
	Lines    []int
}

func (info *fileInfo) Contains(pos int) bool {
	pos -= info.Offset
	return pos >= 0 && pos < len(info.Stream)
}

func (info *fileInfo) GetLocationInfo(pos int) (string, int, int, string) {
	line, col, text := GetLocation(info.Stream, info.Lines, pos-info.Offset)
	return info.Filename, line, col, text
}

type simpleProvider struct {
	files     []*fileInfo
	maxOffset int
}

func (p *simpleProvider) AddFile(filename string, stream []rune) int {
	info := &fileInfo{
		Offset:   p.maxOffset,
		Filename: filename,
		Stream:   stream,
		Lines:    FindLines(stream),
	}
	p.maxOffset += len(stream)
	p.files = append(p.files, info)
	return info.Offset
}

func (p *simpleProvider) GetLocationInfo(pos int) (string, int, int, string) {
	// TODO binary search
	for _, info := range p.files {
		if info.Contains(pos) {
			return info.GetLocationInfo(pos)
		}
	}
	panic(pos)
}

func MakeProvider() LocationProvider {
	return &simpleProvider{}
}
