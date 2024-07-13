package chained

import "fmt"

type Blockchain struct {
	Blocks []BlockHash
	b      map[BlockHash]struct{}
}

func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []BlockHash{}}
}

func (b *Blockchain) AddBlock(h BlockHash) {
	if _, ok := b.b[h]; ok {
		panic(fmt.Sprintf("block %s already exists", h))
	}
	b.Blocks = append(b.Blocks, h)
}
