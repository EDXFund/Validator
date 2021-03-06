package types

import (
	"encoding/binary"
	"io"
	"math/big"
	"sort"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/EDXFund/Validator/common"
	"github.com/EDXFund/Validator/common/hexutil"

	"github.com/EDXFund/Validator/crypto/sha3"
	"github.com/EDXFund/Validator/rlp"
)

var (
	EmptyRootHash = DeriveSha(Transactions{})

//	EmptyUncleHash = CalcUncleHash(nil)
)

type ShardStatus uint16

var (
	ShardMaster ShardStatus = 0xFFFF
)
var (
	ShardEnableLen = 32
)

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

type HeaderIntf interface {
	ShardId() uint16
	Time() *big.Int
	ParentHash() common.Hash
	Coinbase() common.Address
	Difficulty() *big.Int
	Number() *big.Int
	Extra() []byte
	Nonce() BlockNonce
	MixDigest() common.Hash
}

// Header represents a block header in the Ethereum blockchain.
type Header struct {
	ShardId    uint16      `json:"shardId"			gencodec:"required"`
	ParentHash common.Hash `json:"parentHash"       gencodec:"required"`


	Coinbase common.Address `json:"miner"            gencodec:"required"`
	Root     common.Hash    `json:"stateRoot,omitempty"        gencodec:"nil"`
	TxHash   common.Hash    `json:"transactionsRoot,omitempty" gencodec:"nil"`
	ReceiptHash common.Hash `json:"receiptsRoot,omitempty"     gencodec:"nil"`

	Bloom       Bloom       `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int    `json:"difficulty"       gencodec:"required"`
	Number      *big.Int    `json:"number"           gencodec:"required"`
	GasLimit    uint64      `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64      `json:"gasUsed"          gencodec:"required"`
	Time        *big.Int    `json:"timestamp"        gencodec:"required"`
	Extra       []byte      `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash `json:"mixHash"          gencodec:"required"`
	Nonce       BlockNonce  `json:"nonce"            gencodec:"required"`
}

type HeaderMarshal struct {
	ShardId    uint16      `json:"shardId"			gencodec:"required"`
	ParentHash common.Hash `json:"parentHash"       gencodec:"required"`
	Coinbase     common.Address `json:"miner"            gencodec:"required"`
	Root         common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash       common.Hash    `json:"transactionsRoot" gencodec:"required"`
	// hash of shard block included in this master block
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *hexutil.Big   `json:"difficulty"       gencodec:"required"`
	Number      *hexutil.Big   `json:"number"           gencodec:"required"`
	GasLimit    hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
	GasUsed     hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
	Time        *hexutil.Big   `json:"timestamp"        gencodec:"required"`
	Extra       hexutil.Bytes  `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"          gencodec:"required"`
	Nonce       BlockNonce     `json:"nonce"            gencodec:"required"`
	Hash        common.Hash    `json:"hash"`
}
type HeaderUnmarshal struct {
	ShardId    uint16       `json:"shardId"			gencodec:"required"`
	ParentHash *common.Hash `json:"parentHash"       gencodec:"required"`
	Coinbase     *common.Address `json:"miner"            gencodec:"required"`
	Root         *common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash       *common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash *common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       *Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int        `json:"difficulty"       gencodec:"required"`
	Number      *big.Int        `json:"number"           gencodec:"required"`
	GasLimit    *hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
	GasUsed     *hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
	Time        *big.Int        `json:"timestamp"        gencodec:"required"`
	Extra       *hexutil.Bytes  `json:"extraData"        gencodec:"required"`
	MixDigest   *common.Hash    `json:"mixHash"          gencodec:"required"`
	Nonce       *BlockNonce     `json:"nonce"            gencodec:"required"`
}

// field type overrides for gencodec
type scheaderMarshaling struct {
	ShardId    hexutil.Uint64
	Difficulty *hexutil.Big
	Number     *hexutil.Big
	GasLimit   hexutil.Uint64
	GasUsed    hexutil.Uint64
	Time       *hexutil.Big
	Extra      hexutil.Bytes
	Hash       common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h)
}

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Size() common.StorageSize {
	return common.StorageSize(unsafe.Sizeof(*h)) + common.StorageSize(len(h.Extra)+(h.Difficulty.BitLen()+h.Number.BitLen()+h.Time.BitLen())/8)
}




// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions and uncles) together.
type Body struct {
	Transactions []*Transaction

	//receipts
	Receipts ContractResults


}

// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	header *Header

	transactions Transactions
	receipts     ContractResults


	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}
type Blocks []*Block

// DeprecatedTd is an old relic for extracting the TD of a block. It is in the
// code solely to facilitate upgrading the database from the old format to the
// new, after which it should be deleted. Do not use!
func (b *Block) DeprecatedTd() *big.Int {
	return b.td
}

// [deprecated by eth/63]
// StorageBlock defines the RLP encoding of a Block stored in the
// state database. The StorageBlock encoding contains fields that
// would otherwise need to be recomputed.
type StorageBlock Block

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header   *Header
	Txs      []*Transaction
	Receipts []*ContractResult
}

// [deprecated by eth/63]
// "storage" block encoding. used for database.
type storageblock struct {
	Header   *Header
	Txs      []*Transaction
	Receipts []*ContractResult
	TD       *big.Int
}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs, uncles
// and receipts.
func NewBlock(header *Header, txs []*Transaction, receipts []*ContractResult) *Block {
	b := &Block{header: CopyHeader(header), td: new(big.Int)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		b.header.TxHash = DeriveSha(Transactions(txs))
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = EmptyRootHash
	} else {
		b.header.ReceiptHash = DeriveSha(ContractResults(receipts))
		b.header.Bloom = CreateBloom(receipts)
	}


	return b
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Time = new(big.Int); h.Time != nil {
		cpy.Time.Set(h.Time)
	}
	if cpy.Difficulty = new(big.Int); h.Difficulty != nil {
		cpy.Difficulty.Set(h.Difficulty)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	return &cpy
}

// DecodeRLP decodes the Ethereum
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.transactions, b.receipts = eb.Header, eb.Txs, eb.Receipts
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		Header:   b.header,
		Txs:      b.transactions,
		Receipts: b.receipts,
	})
}

// [deprecated by eth/63]
func (b *StorageBlock) DecodeRLP(s *rlp.Stream) error {
	var sb storageblock
	if err := s.Decode(&sb); err != nil {
		return err
	}
	b.header, b.transactions, b.receipts, b.td = sb.Header, sb.Txs, sb.Receipts, sb.TD
	return nil
}

// TODO: copies


func (b *Block) Transactions() Transactions { return b.transactions }

func (b *Block) Receiptions() ContractResults { return b.receipts }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}
func (b *Block) ContractReceipts() ContractResults { return b.receipts }

func (b *Block) ContrcatReceipt(hash common.Hash) *ContractResult {
	for _, receipt := range b.receipts {
		if receipt.TxHash == hash {
			return receipt
		}
	}
	return nil
}
func (b *Block) ShardId() uint16      { return b.header.ShardId }
func (b *Block) Number() *big.Int     { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasLimit() uint64     { return b.header.GasLimit }
func (b *Block) GasUsed() uint64      { return b.header.GasUsed }
func (b *Block) Difficulty() *big.Int { return new(big.Int).Set(b.header.Difficulty) }
func (b *Block) Time() *big.Int       { return new(big.Int).Set(b.header.Time) }

func (b *Block) NumberU64() uint64        { return b.header.Number.Uint64() }
func (b *Block) MixDigest() common.Hash   { return b.header.MixDigest }
func (b *Block) Nonce() uint64            { return binary.BigEndian.Uint64(b.header.Nonce[:]) }
func (b *Block) Bloom() Bloom             { return b.header.Bloom }
func (b *Block) Coinbase() common.Address { return b.header.Coinbase }
func (b *Block) Root() common.Hash        { return b.header.Root }
func (b *Block) ParentHash() common.Hash  { return b.header.ParentHash }
func (b *Block) TxHash() common.Hash      { return b.header.TxHash }
func (b *Block) ReceiptHash() common.Hash { return b.header.ReceiptHash }

//func (b *Block) UncleHash() common.Hash   { return b.header.UncleHash }
func (b *Block) Extra() []byte                 { return common.CopyBytes(b.header.Extra) }

func (b *Block) Header() *Header               { return CopyHeader(b.header) }

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.transactions, b.receipts} }

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header

	return &Block{
		header:       &cpy,
		transactions: b.transactions,
		receipts:     b.receipts,
	}
}

// WithBody returns a new block with the given transaction and uncle contents.
func (b *Block) WithBody(transactions []*Transaction, contractReceipts ContractResults) *Block {
	block := &Block{
		header: CopyHeader(b.header),
	}
	block.transactions = make([]*Transaction, len(transactions))
	block.receipts = make(ContractResults, len(contractReceipts))

	copy(block.transactions, transactions)
	copy(block.receipts, contractReceipts)

	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := scblockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type scblockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self scblockSorter) Len() int { return len(self.blocks) }
func (self scblockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self scblockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.header.Number.Cmp(b2.header.Number) < 0 }
