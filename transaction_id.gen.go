package iotago

// Code generated by go generate; DO NOT EDIT. Check gen/ directory instead.

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"

	"golang.org/x/crypto/blake2b"

	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/iota.go/v4/hexutil"
)

const (
	TransactionIDLength = IdentifierLength + SlotIndexLength
)

var (
	ErrInvalidTransactionIDLength = ierrors.New("invalid transactionID length")

	EmptyTransactionID = TransactionID{}
)

// TransactionID is a 32 byte hash value together with an 4 byte slot index.
type TransactionID [TransactionIDLength]byte

// TransactionIDRepresentingData returns a new TransactionID for the given data by hashing it with blake2b and associating it with the given slot index.
func TransactionIDRepresentingData(slot SlotIndex, data []byte) TransactionID {
	return NewTransactionID(slot, blake2b.Sum256(data))
}

func NewTransactionID(slot SlotIndex, idBytes Identifier) TransactionID {
	t := TransactionID{}
	copy(t[:], idBytes[:])
	binary.LittleEndian.PutUint32(t[IdentifierLength:], uint32(slot))

	return t
}

// TransactionIDFromHexString converts the hex to a TransactionID representation.
func TransactionIDFromHexString(hex string) (TransactionID, error) {
	b, err := hexutil.DecodeHex(hex)
	if err != nil {
		return EmptyTransactionID, err
	}

	s, _, err := TransactionIDFromBytes(b)

	return s, err
}

// TransactionIDFromBytes returns a new TransactionID represented by the passed bytes.
func TransactionIDFromBytes(b []byte) (TransactionID, int, error) {
	if len(b) < TransactionIDLength {
		return EmptyTransactionID, 0, ErrInvalidTransactionIDLength
	}

	return TransactionID(b), TransactionIDLength, nil
}

// MustTransactionIDFromHexString converts the hex to a TransactionID representation.
func MustTransactionIDFromHexString(hex string) TransactionID {
	t, err := TransactionIDFromHexString(hex)
	if err != nil {
		panic(err)
	}

	return t
}

func (t TransactionID) Bytes() ([]byte, error) {
	return t[:], nil
}

func (t TransactionID) MarshalText() (text []byte, err error) {
	dst := make([]byte, hex.EncodedLen(len(EmptyTransactionID)))
	hex.Encode(dst, t[:])

	return dst, nil
}

func (t *TransactionID) UnmarshalText(text []byte) error {
	_, err := hex.Decode(t[:], text)

	return err
}

// Empty tells whether the TransactionID is empty.
func (t TransactionID) Empty() bool {
	return t == EmptyTransactionID
}

// ToHex converts the Identifier to its hex representation.
func (t TransactionID) ToHex() string {
	return hexutil.EncodeHex(t[:])
}

func (t TransactionID) String() string {
	return fmt.Sprintf("TransactionID(%s:%d)", t.Alias(), t.Slot())
}

func (t TransactionID) Slot() SlotIndex {
	return SlotIndex(binary.LittleEndian.Uint32(t[IdentifierLength:]))
}

// Index returns a slot index to conform with hive's IndexedID interface.
func (t TransactionID) Index() SlotIndex {
	return t.Slot()
}

func (t TransactionID) Identifier() Identifier {
	return Identifier(t[:IdentifierLength])
}

var (
	// TransactionIDAliases contains a dictionary of identifiers associated to their human-readable alias.
	TransactionIDAliases = make(map[TransactionID]string)

	// transactionIDAliasesMutex is the mutex that is used to synchronize access to the previous map.
	transactionIDAliasesMutex = sync.RWMutex{}
)

// RegisterAlias allows to register a human-readable alias for the Identifier which will be used as a replacement for
// the String method.
func (t TransactionID) RegisterAlias(alias string) {
	transactionIDAliasesMutex.Lock()
	defer transactionIDAliasesMutex.Unlock()

	TransactionIDAliases[t] = alias
}

// Alias returns the human-readable alias of the Identifier (or the base58 encoded bytes of no alias was set).
func (t TransactionID) Alias() (alias string) {
	transactionIDAliasesMutex.RLock()
	defer transactionIDAliasesMutex.RUnlock()

	if existingAlias, exists := TransactionIDAliases[t]; exists {
		return existingAlias
	}

	return t.ToHex()
}

// UnregisterAlias allows to unregister a previously registered alias.
func (t TransactionID) UnregisterAlias() {
	transactionIDAliasesMutex.Lock()
	defer transactionIDAliasesMutex.Unlock()

	delete(TransactionIDAliases, t)
}

// Compare compares two TransactionIDs.
func (t TransactionID) Compare(other TransactionID) int {
	return bytes.Compare(t[:], other[:])
}

type TransactionIDs []TransactionID

// ToHex converts the TransactionIDs to their hex representation.
func (ids TransactionIDs) ToHex() []string {
	hexIDs := make([]string, len(ids))
	for i, t := range ids {
		hexIDs[i] = hexutil.EncodeHex(t[:])
	}

	return hexIDs
}

// RemoveDupsAndSort removes duplicated TransactionIDs and sorts the slice by the lexical ordering.
func (ids TransactionIDs) RemoveDupsAndSort() TransactionIDs {
	sorted := append(TransactionIDs{}, ids...)
	sort.Slice(sorted, func(i, j int) bool {
		return bytes.Compare(sorted[i][:], sorted[j][:]) == -1
	})

	var result TransactionIDs
	var prev TransactionID
	for i, t := range sorted {
		if i == 0 || !bytes.Equal(prev[:], t[:]) {
			result = append(result, t)
		}
		prev = t
	}

	return result
}

// Sort sorts the TransactionIDs lexically and in-place.
func (ids TransactionIDs) Sort() {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Compare(ids[j]) < 0
	})
}

// TransactionIDsFromHexString converts the given block IDs from their hex to TransactionID representation.
func TransactionIDsFromHexString(TransactionIDsHex []string) (TransactionIDs, error) {
	result := make(TransactionIDs, len(TransactionIDsHex))

	for i, hexString := range TransactionIDsHex {
		TransactionID, err := TransactionIDFromHexString(hexString)
		if err != nil {
			return nil, err
		}
		result[i] = TransactionID
	}

	return result, nil
}
