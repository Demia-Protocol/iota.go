package iotago

// Code generated by go generate; DO NOT EDIT. Check gen/ directory instead.

import (
	"encoding/hex"
	"sync"

	"golang.org/x/crypto/blake2b"

	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/iota.go/v4/hexutil"
)

const (
	// AccountIDLength defines the length of an AccountID.
	AccountIDLength = blake2b.Size256
)

var (
	EmptyAccountID = AccountID{}

	ErrInvalidAccountIDLength = ierrors.New("invalid AccountID length")
)

// AccountID is a 32 byte hash value.
type AccountID [AccountIDLength]byte

type AccountIDs []AccountID

// AccountIDFromData returns a new AccountID for the given data by hashing it with blake2b.
func AccountIDFromData(data []byte) AccountID {
	return blake2b.Sum256(data)
}

// AccountIDFromHexString converts the hex to an AccountID representation.
func AccountIDFromHexString(hex string) (AccountID, error) {
	bytes, err := hexutil.DecodeHex(hex)
	if err != nil {
		return EmptyAccountID, err
	}

	a, _, err := AccountIDFromBytes(bytes)

	return a, err
}

// MustAccountIDFromHexString converts the hex to an AccountID representation.
func MustAccountIDFromHexString(hex string) AccountID {
	a, err := AccountIDFromHexString(hex)
	if err != nil {
		panic(err)
	}

	return a
}

func AccountIDFromBytes(bytes []byte) (AccountID, int, error) {
	var a AccountID
	if len(bytes) < AccountIDLength {
		return a, 0, ErrInvalidAccountIDLength
	}
	copy(a[:], bytes)

	return a, len(bytes), nil
}

func (a AccountID) Bytes() ([]byte, error) {
	return a[:], nil
}

func (a AccountID) MarshalText() (text []byte, err error) {
	dst := make([]byte, hex.EncodedLen(len(EmptyAccountID)))
	hex.Encode(dst, a[:])

	return dst, nil
}

func (a *AccountID) UnmarshalText(text []byte) error {
	_, err := hex.Decode(a[:], text)

	return err
}

// Empty tells whether the AccountID is empty.
func (a AccountID) Empty() bool {
	return a == EmptyAccountID
}

// ToHex converts the AccountID to its hex representation.
func (a AccountID) ToHex() string {
	return hexutil.EncodeHex(a[:])
}

func (a AccountID) String() string {
	return a.Alias()
}

var (
	// accountIDAliases contains a dictionary of AccountIDs associated to their human-readable alias.
	accountIDAliases = make(map[AccountID]string)

	// accountIDAliasesMutex is the mutex that is used to synchronize access to the previous map.
	accountIDAliasesMutex = sync.RWMutex{}
)

// RegisterAlias allows to register a human-readable alias for the AccountID which will be used as a replacement for
// the String method.
func (a AccountID) RegisterAlias(alias string) {
	accountIDAliasesMutex.Lock()
	defer accountIDAliasesMutex.Unlock()

	accountIDAliases[a] = alias
}

// Alias returns the human-readable alias of the AccountID (or the hex encoded bytes if no alias was set).
func (a AccountID) Alias() (alias string) {
	accountIDAliasesMutex.RLock()
	defer accountIDAliasesMutex.RUnlock()

	if existingAlias, exists := accountIDAliases[a]; exists {
		return existingAlias
	}

	return a.ToHex()
}

// UnregisterAlias allows to unregister a previously registered alias.
func (a AccountID) UnregisterAlias() {
	accountIDAliasesMutex.Lock()
	defer accountIDAliasesMutex.Unlock()

	delete(accountIDAliases, a)
}

// UnregisterAccountIDAliases allows to unregister all previously registered aliases.
func UnregisterAccountIDAliases() {
	accountIDAliasesMutex.Lock()
	defer accountIDAliasesMutex.Unlock()

	accountIDAliases = make(map[AccountID]string)
}

// Matches checks whether other matches this ChainID.
func (a AccountID) Matches(other ChainID) bool {
	otherAccountID, isAccountID := other.(AccountID)
	if !isAccountID {
		return false
	}

	return a == otherAccountID
}

// Addressable tells whether this ChainID can be converted into a ChainAddress.
func (a AccountID) Addressable() bool {
	return true
}

// ToAddress converts this ChainID into an ChainAddress.
func (a AccountID) ToAddress() ChainAddress {
	var addr AccountAddress
	copy(addr[:], a[:])

	return &addr
}

// Key returns a key to use to index this ChainID.
func (a AccountID) Key() interface{} {
	return a.String()
}

// FromOutputID returns the ChainID computed from a given OutputID.
func (a AccountID) FromOutputID(in OutputID) ChainID {
	return AccountIDFromOutputID(in)
}

// AccountIDFromOutputID returns the AccountID computed from a given OutputID.
func AccountIDFromOutputID(outputID OutputID) AccountID {
	return blake2b.Sum256(outputID[:])
}
