package iotago

import (
	"context"
	"fmt"
	"time"

	"github.com/iotaledger/hive.go/lo"
	"github.com/iotaledger/hive.go/runtime/options"
)

// V3ProtocolParameters defines the parameters of the protocol.
type V3ProtocolParameters struct {
	v3ProtocolParameters `serix:"0"`
}

type v3ProtocolParameters struct {
	// Version defines the version of the protocol this protocol parameters are for.
	Version Version `serix:"0,mapKey=version"`

	// NetworkName defines the human friendly name of the network.
	NetworkName string `serix:"1,lengthPrefixType=uint8,mapKey=networkName"`
	// Bech32HRP defines the HRP prefix used for Bech32 addresses in the network.
	Bech32HRP NetworkPrefix `serix:"2,lengthPrefixType=uint8,mapKey=bech32Hrp"`

	// RentStructure defines the rent structure used by given node/network.
	RentStructure RentStructure `serix:"3,mapKey=rentStructure"`
	// WorkScoreStructure defines the work score structure used by given node/network.
	WorkScoreStructure WorkScoreStructure `serix:"4,mapKey=workScoreStructure"`
	// TokenSupply defines the current token supply on the network.
	TokenSupply BaseToken `serix:"5,mapKey=tokenSupply"`

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
	// LivenessThreshold is used by tip-selection to determine the if a block is eligible by evaluating issuingTimes
	// and commitments in its past-cone to ATT and lastCommittedSlot respectively.
	LivenessThreshold SlotIndex `serix:"17,mapKey=livenessThreshold"`
	// EpochNearingThreshold is used by the epoch orchestrator to detect the slot that should trigger a new committee
	// selection for the next and upcoming epoch.
	EpochNearingThreshold SlotIndex `serix:"18,mapKey=epochNearingThreshold"`

	VersionSignaling VersionSignaling `serix:"19,mapKey=versionSignaling"`
}

func (p v3ProtocolParameters) Equals(other v3ProtocolParameters) bool {
	return p.Version == other.Version &&
		p.NetworkName == other.NetworkName &&
		p.Bech32HRP == other.Bech32HRP &&
		p.RentStructure.Equals(other.RentStructure) &&
		p.WorkScoreStructure.Equals(other.WorkScoreStructure) &&
		p.TokenSupply == other.TokenSupply &&
		p.GenesisUnixTimestamp == other.GenesisUnixTimestamp &&
		p.SlotDurationInSeconds == other.SlotDurationInSeconds &&
		p.SlotsPerEpochExponent == other.SlotsPerEpochExponent &&
		p.ManaGenerationRate == other.ManaGenerationRate &&
		p.ManaGenerationRateExponent == other.ManaGenerationRateExponent &&
		lo.Equal(p.ManaDecayFactors, other.ManaDecayFactors) &&
		p.ManaDecayFactorsExponent == other.ManaDecayFactorsExponent &&
		p.ManaDecayFactorEpochsSum == other.ManaDecayFactorEpochsSum &&
		p.ManaDecayFactorEpochsSumExponent == other.ManaDecayFactorEpochsSumExponent &&
		p.StakingUnbondingPeriod == other.StakingUnbondingPeriod &&
		p.EvictionAge == other.EvictionAge &&
		p.LivenessThreshold == other.LivenessThreshold &&
		p.EpochNearingThreshold == other.EpochNearingThreshold &&
		p.VersionSignaling.Equals(other.VersionSignaling)
}

func NewV3ProtocolParameters(opts ...options.Option[V3ProtocolParameters]) *V3ProtocolParameters {
	return options.Apply(
		new(V3ProtocolParameters),
		append([]options.Option[V3ProtocolParameters]{
			WithNetworkOptions("testnet", PrefixTestnet),
			WithSupplyOptions(1813620509061365, 100, 1, 10),
			WithWorkScoreOptions(1, 1, 1, 10, 5, 1, 2, 2, 2, 10, 4),
			WithTimeProviderOptions(time.Now().Unix(), 10, 13),
			// TODO: add sane default values
			WithManaOptions(1,
				0,
				[]uint32{
					10,
					20,
				},
				0,
				0,
				0,
			),
			WithLivenessOptions(10, 3, 4),
			WithStakingOptions(10),
			WithVersionSignalingOptions(7, 5, 7),
		},
			opts...,
		),
		func(p *V3ProtocolParameters) {
			p.v3ProtocolParameters.Version = apiV3Version
		},
	)
}

var _ ProtocolParameters = &V3ProtocolParameters{}

func (p *V3ProtocolParameters) Version() Version {
	return p.v3ProtocolParameters.Version
}

func (p *V3ProtocolParameters) Bech32HRP() NetworkPrefix {
	return p.v3ProtocolParameters.Bech32HRP
}

func (p *V3ProtocolParameters) NetworkName() string {
	return p.v3ProtocolParameters.NetworkName
}

func (p *V3ProtocolParameters) RentStructure() *RentStructure {
	return &p.v3ProtocolParameters.RentStructure
}

func (p *V3ProtocolParameters) WorkScoreStructure() *WorkScoreStructure {
	return &p.v3ProtocolParameters.WorkScoreStructure
}

func (p *V3ProtocolParameters) TokenSupply() BaseToken {
	return p.v3ProtocolParameters.TokenSupply
}

func (p *V3ProtocolParameters) NetworkID() NetworkID {
	return NetworkIDFromString(p.v3ProtocolParameters.NetworkName)
}

func (p *V3ProtocolParameters) TimeProvider() *TimeProvider {
	return NewTimeProvider(p.v3ProtocolParameters.GenesisUnixTimestamp, int64(p.v3ProtocolParameters.SlotDurationInSeconds), p.v3ProtocolParameters.SlotsPerEpochExponent)
}

// EpochDurationInSlots defines the amount of slots in an epoch.
func (p *V3ProtocolParameters) ParamEpochDurationInSlots() SlotIndex {
	return 1 << p.v3ProtocolParameters.SlotsPerEpochExponent
}

func (p *V3ProtocolParameters) StakingUnbondingPeriod() EpochIndex {
	return p.v3ProtocolParameters.StakingUnbondingPeriod
}

func (p *V3ProtocolParameters) LivenessThreshold() SlotIndex {
	return p.v3ProtocolParameters.LivenessThreshold
}

func (p *V3ProtocolParameters) EvictionAge() SlotIndex {
	return p.v3ProtocolParameters.EvictionAge
}

func (p *V3ProtocolParameters) EpochNearingThreshold() SlotIndex {
	return p.v3ProtocolParameters.EpochNearingThreshold
}

func (p *V3ProtocolParameters) VersionSignaling() *VersionSignaling {
	return &p.v3ProtocolParameters.VersionSignaling
}

func (p *V3ProtocolParameters) Bytes() ([]byte, error) {
	return commonSerixAPI().Encode(context.TODO(), p)
}

func (p *V3ProtocolParameters) Hash() (Identifier, error) {
	bytes, err := p.Bytes()
	if err != nil {
		return Identifier{}, err
	}

	return IdentifierFromData(bytes), nil
}

func (p *V3ProtocolParameters) String() string {
	return fmt.Sprintf("ProtocolParameters: {\n\tVersion: %d\n\tNetwork Name: %s\n\tBech32 HRP Prefix: %s\n\tRent Structure: %v\n\tWorkScore Structure: %v\n\tToken Supply: %d\n\tGenesis Unix Timestamp: %d\n\tSlot Duration in Seconds: %d\n\tSlots per Epoch Exponent: %d\n\tMana Generation Rate: %d\n\tMana Generation Rate Exponent: %d\t\nMana Decay Factors: %v\n\tMana Decay Factors Exponent: %d\n\tMana Decay Factor Epochs Sum: %d\n\tMana Decay Factor Epochs Sum Exponent: %d\n\tStaking Unbonding Period: %d\n\tEviction Age: %d\n\tLiveness Threshold: %d\n}",
		p.v3ProtocolParameters.Version, p.v3ProtocolParameters.NetworkName, p.v3ProtocolParameters.Bech32HRP, p.v3ProtocolParameters.RentStructure, p.v3ProtocolParameters.WorkScoreStructure, p.v3ProtocolParameters.TokenSupply, p.v3ProtocolParameters.GenesisUnixTimestamp, p.v3ProtocolParameters.SlotDurationInSeconds, p.v3ProtocolParameters.SlotsPerEpochExponent, p.v3ProtocolParameters.ManaGenerationRate, p.v3ProtocolParameters.ManaGenerationRateExponent, p.v3ProtocolParameters.ManaDecayFactors, p.v3ProtocolParameters.ManaDecayFactorsExponent, p.v3ProtocolParameters.ManaDecayFactorEpochsSum, p.v3ProtocolParameters.ManaDecayFactorEpochsSumExponent, p.v3ProtocolParameters.StakingUnbondingPeriod, p.v3ProtocolParameters.EvictionAge, p.v3ProtocolParameters.LivenessThreshold)
}

func (p *V3ProtocolParameters) ManaDecayProvider() *ManaDecayProvider {
	return NewManaDecayProvider(p.TimeProvider(), p.v3ProtocolParameters.SlotsPerEpochExponent, p.v3ProtocolParameters.ManaGenerationRate, p.v3ProtocolParameters.ManaGenerationRateExponent, p.v3ProtocolParameters.ManaDecayFactors, p.v3ProtocolParameters.ManaDecayFactorsExponent, p.v3ProtocolParameters.ManaDecayFactorEpochsSum, p.v3ProtocolParameters.ManaDecayFactorEpochsSumExponent)
}

func (p *V3ProtocolParameters) Equals(other *V3ProtocolParameters) bool {
	return p.v3ProtocolParameters.Equals(other.v3ProtocolParameters)
}

func WithNetworkOptions(networkName string, bech32HRP NetworkPrefix) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.NetworkName = networkName
		p.v3ProtocolParameters.Bech32HRP = bech32HRP
	}
}

func WithSupplyOptions(totalSupply BaseToken, vByteCost uint32, vBFactorData VByteCostFactor, vBFactorKey VByteCostFactor) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.TokenSupply = totalSupply
		p.v3ProtocolParameters.RentStructure = RentStructure{
			VByteCost:    vByteCost,
			VBFactorData: vBFactorData,
			VBFactorKey:  vBFactorKey,
		}
	}
}

func WithWorkScoreOptions(output WorkScore, staking WorkScore, blockIssuer WorkScore, ed25519Signature WorkScore, nativeToken WorkScore, data WorkScoreFactor, input WorkScoreFactor, contextInput WorkScoreFactor, allotment WorkScoreFactor, missingParent WorkScoreFactor, minStrongParentsThreshold byte) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.WorkScoreStructure = WorkScoreStructure{
			WorkScores: WorkScores{
				Output:           output,
				Staking:          staking,
				BlockIssuer:      blockIssuer,
				Ed25519Signature: ed25519Signature,
				NativeToken:      nativeToken,
			},
			Factors: WorkScoreFactors{
				Data:          data,
				Input:         input,
				ContextInput:  contextInput,
				Allotment:     allotment,
				MissingParent: missingParent,
			},
			MinStrongParentsThreshold: minStrongParentsThreshold,
		}
	}
}

func WithTimeProviderOptions(genesisTimestamp int64, slotDuration uint8, slotsPerEpochExponent uint8) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.GenesisUnixTimestamp = genesisTimestamp
		p.v3ProtocolParameters.SlotDurationInSeconds = slotDuration
		p.v3ProtocolParameters.SlotsPerEpochExponent = slotsPerEpochExponent
	}
}

func WithManaOptions(manaGenerationRate uint8, manaGenerationRateExponent uint8, manaDecayFactors []uint32, manaDecayFactorsExponent uint8, manaDecayFactorEpochsSum uint32, manaDecayFactorEpochsSumExponent uint8) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.ManaGenerationRate = manaGenerationRate
		p.v3ProtocolParameters.ManaGenerationRateExponent = manaGenerationRateExponent
		p.v3ProtocolParameters.ManaDecayFactors = manaDecayFactors
		p.v3ProtocolParameters.ManaDecayFactorsExponent = manaDecayFactorsExponent
		p.v3ProtocolParameters.ManaDecayFactorEpochsSum = manaDecayFactorEpochsSum
		p.v3ProtocolParameters.ManaDecayFactorEpochsSumExponent = manaDecayFactorEpochsSumExponent
	}
}

func WithLivenessOptions(evictionAge SlotIndex, livenessThreshold SlotIndex, epochNearingThreshold SlotIndex) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.EvictionAge = evictionAge
		p.v3ProtocolParameters.LivenessThreshold = livenessThreshold
		p.v3ProtocolParameters.EpochNearingThreshold = epochNearingThreshold
	}
}

func WithStakingOptions(unboundPeriod EpochIndex) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.StakingUnbondingPeriod = unboundPeriod
	}
}

func WithVersionSignalingOptions(windowSize uint8, windowTargetRatio uint8, activationOffset uint8) options.Option[V3ProtocolParameters] {
	return func(p *V3ProtocolParameters) {
		p.v3ProtocolParameters.VersionSignaling = VersionSignaling{
			WindowSize:        windowSize,
			WindowTargetRatio: windowTargetRatio,
			ActivationOffset:  activationOffset,
		}
	}
}
