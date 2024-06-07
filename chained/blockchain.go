package chained

type Blockchain struct {
	Blocks []BlockHash
}

func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []BlockHash{}}
}

func (b *Blockchain) AddBlock(h BlockHash) {
	b.Blocks = append(b.Blocks, h)
}
