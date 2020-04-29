package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte
	StudentData  []byte
	Signature    []byte
	Company      []byte
	Verification []byte
	PrevHash     []byte
	Nonce        int
}

func CreateBlock(data []byte, signature []byte, company []byte, verification []byte, prevHash []byte) *Block {
	block := &Block{[]byte{}, data, signature, company, verification, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func InitFirstBlock() *Block {
	return CreateBlock([]byte("Initial Block"), []byte(""), []byte(""), []byte(""), []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
