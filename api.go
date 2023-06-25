package iotago

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/crypto/blake2b"

	"github.com/iotaledger/hive.go/serializer/v2/serix"
)

var (
	// ErrMissingProtocolParams is returned when ProtocolParameters are missing for operations which require them.
	ErrMissingProtocolParams = errors.New("missing protocol parameters")

	// internal API instance used to encode/decode objects where protocol parameters don't matter.
	_internalAPI   API
	_internalAPIMu = sync.RWMutex{}
)

func init() {
	_internalAPI = V3API(&ProtocolParameters{})
}

// API handles en/decoding of IOTA protocol objects.
type API interface {
	// Encode encodes the given object to bytes.
	Encode(obj any, opts ...serix.Option) ([]byte, error)
	// Decode decodes the given bytes into object.
	Decode(b []byte, obj any, opts ...serix.Option) (int, error)
	// JSONEncode encodes the given object to its json representation.
	JSONEncode(obj any, opts ...serix.Option) ([]byte, error)
	// JSONDecode decodes the json data into object.
	JSONDecode(jsonData []byte, obj any, opts ...serix.Option) error
	// Underlying returns the underlying serix.API instance.
	Underlying() *serix.API
	// TimeProvider returns the underlying time provider used.
	TimeProvider() *TimeProvider
	// ManaDecayProvider returns the underlying mana decay provider used.
	ManaDecayProvider() *ManaDecayProvider
}

// LatestAPI creates a new API instance conforming to the latest IOTA protocol version.
func LatestAPI(protoParams *ProtocolParameters) API {
	return V3API(protoParams)
}

// calls the internally instantiated API to encode the given object.
//
//nolint:unparam
func internalEncode(obj any, opts ...serix.Option) ([]byte, error) {
	_internalAPIMu.RLock()
	defer _internalAPIMu.RUnlock()
	return _internalAPI.Encode(obj, opts...)
}

// calls the internally instantiated API to decode the given object.
func internalDecode(b []byte, obj any, opts ...serix.Option) (int, error) {
	_internalAPIMu.RLock()
	defer _internalAPIMu.RUnlock()
	return _internalAPI.Decode(b, obj, opts...)
}

// SwapInternalAPI swaps the internally used API of this lib with new.
func SwapInternalAPI(newAPI API) {
	_internalAPIMu.Lock()
	defer _internalAPIMu.Unlock()
	_internalAPI = newAPI
}

// NetworkID defines the ID of the network on which entities operate on.
type NetworkID = uint64

// NetworkIDFromString returns the network ID string's numerical representation.
func NetworkIDFromString(networkIDStr string) NetworkID {
	networkIDBlakeHash := blake2b.Sum256([]byte(networkIDStr))
	return binary.LittleEndian.Uint64(networkIDBlakeHash[:])
}

type protocolAPIContext string

// ProtocolAPIContextKey defines the key to use for a context containing a *ProtocolParameters.
const ProtocolAPIContextKey protocolAPIContext = "protocolParameters"

// ProtocolParameters defines the parameters of the protocol.
type ProtocolParameters struct {
	// The version of the protocol running.
	Version byte `serix:"0,mapKey=version"`
	// The human friendly name of the network.
	NetworkName string `serix:"1,lengthPrefixType=uint8,mapKey=networkName"`
	// The HRP prefix used for Bech32 addresses in the network.
	Bech32HRP NetworkPrefix `serix:"2,lengthPrefixType=uint8,mapKey=bech32Hrp"`
	// The minimum pow score of the network.
	MinPoWScore uint32 `serix:"3,mapKey=minPowScore"`
	// The rent structure used by given node/network.
	RentStructure RentStructure `serix:"4,mapKey=rentStructure"`
	// TokenSupply defines the current token supply on the network.
	TokenSupply uint64 `serix:"5,mapKey=tokenSupply"`
	// GenesisUnixTimestamp defines the genesis timestamp at which the slots start to count.
	GenesisUnixTimestamp int64 `serix:"6,mapKey=genesisUnixTimestamp"`
	// SlotDurationInSeconds defines the duration of each slot in seconds.
	SlotDurationInSeconds uint8 `serix:"7,mapKey=slotDurationInSeconds"`
	// SlotsPerEpochExponent is the number of slots in an epoch expressed as an exponent of 2.
	// (2**SlotsPerEpochExponent) == slots in an epoch.
	SlotsPerEpochExponent uint8 `serix:"8,mapKey=slotsPerEpochExponent"`
	// ManaGenerationRate is the amount of potential Mana generated by 1 IOTA in 1 slot.
	ManaGenerationRate uint8 `serix:"9,mapKey=manaGenerationRate"`
	// ManaGenerationRateExponent is the scaling of ManaGenerationRate expressed as an exponent of 2.
	ManaGenerationRateExponent uint8 `serix:"10,mapKey=manaGenerationRateExponent"`
	// ManaDecayFactors is a lookup table of epoch index diff to mana decay factor (slice index 0 = 1 epoch).
	ManaDecayFactors []uint32 `serix:"11,lengthPrefixType=uint16,mapKey=manaDecayFactors"`
	// ManaDecayFactorsExponent is the scaling of ManaDecayFactors expressed as an exponent of 2.
	ManaDecayFactorsExponent uint8 `serix:"12,mapKey=manaDecayFactorsExponent"`
	// ManaDecayFactorEpochsSum is an integer approximation of the sum of decay over epochs.
	ManaDecayFactorEpochsSum uint32 `serix:"13,mapKey=manaDecayFactorEpochsSum"`
	// ManaDecayFactorEpochsSumExponent is the scaling of ManaDecayFactorEpochsSum expressed as an exponent of 2.
	ManaDecayFactorEpochsSumExponent uint8 `serix:"14,mapKey=manaDecayFactorEpochsSumExponent"`
	// StakingUnbondingPeriod defines the unbonding period in epochs before an account can stop staking.
	StakingUnbondingPeriod EpochIndex `serix:"15,mapKey=stakingUnbondingPeriod"`
	// EvictionAge defines the age in slots when you can evict blocks by committing them into a slot commitments and
	// when slots stop being a consumable accounts' state relative to the latest committed slot.
	EvictionAge SlotIndex `serix:"16,mapKey=evictionAge"`
	// LivenessThreshold is used by tipselection to determine the if a block is eligibile by evaluating issuingTimes
	// and commitments in its pastcone to ATT and lastCommittedSlot respectively.
	LivenessThreshold SlotIndex `serix:"17,mapKey=liveNessThreshold"`
}

func (p ProtocolParameters) AsSerixContext() context.Context {
	return context.WithValue(context.Background(), ProtocolAPIContextKey, &p)
}

func (p ProtocolParameters) NetworkID() NetworkID {
	return NetworkIDFromString(p.NetworkName)
}

func (p ProtocolParameters) TimeProvider() *TimeProvider {
	return NewTimeProvider(p.GenesisUnixTimestamp, int64(p.SlotDurationInSeconds), p.SlotsPerEpochExponent)
}

// EpochDurationInSlots defines the amount of slots in an epoch.
func (p ProtocolParameters) EpochDurationInSlots() SlotIndex {
	return 1 << p.SlotsPerEpochExponent
}

func (p ProtocolParameters) String() string {
	return fmt.Sprintf("ProtocolParameters: {\n\tVersion: %d\n\tNetwork Name: %s\n\tBech32 HRP Prefix: %s\n\tMinimum PoW Score: %d\n\tRent Structure: %v\n\tToken Supply: %d\n\tGenesis Unix Timestamp: %d\n\tSlot Duration in Seconds: %d\n\tSlots per Epoch Exponent: %d\n\tMana Generation Rate: %d\n\tMana Generation Rate Exponent: %d\t\nMana Decay Factors: %v\n\tMana Decay Factors Exponent: %d\n\tMana Decay Factor Epochs Sum: %d\n\tMana Decay Factor Epochs Sum Exponent: %d\n\tStaking Unbonding Period: %d\n\tEviction Age: %d\n\tLiveness Threshold: %d\n}",
		p.Version, p.NetworkName, p.Bech32HRP, p.MinPoWScore, p.RentStructure, p.TokenSupply, p.GenesisUnixTimestamp, p.SlotDurationInSeconds, p.SlotsPerEpochExponent, p.ManaGenerationRate, p.ManaGenerationRateExponent, p.ManaDecayFactors, p.ManaDecayFactorsExponent, p.ManaDecayFactorEpochsSum, p.ManaDecayFactorEpochsSumExponent, p.StakingUnbondingPeriod, p.EvictionAge, p.LivenessThreshold)
}

func (p ProtocolParameters) ManaDecayProvider() *ManaDecayProvider {
	return NewManaDecayProvider(p.TimeProvider(), p.SlotsPerEpochExponent, p.ManaGenerationRate, p.ManaGenerationRateExponent, p.ManaDecayFactors, p.ManaDecayFactorsExponent, p.ManaDecayFactorEpochsSum, p.ManaDecayFactorEpochsSumExponent)
}

// Sizer is an object knowing its own byte size.
type Sizer interface {
	// Size returns the size of the object in terms of bytes.
	Size() int
}
