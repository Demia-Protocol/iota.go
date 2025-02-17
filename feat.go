package iotago

import (
	"fmt"
	"sort"

	"github.com/iotaledger/hive.go/constraints"
	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/hive.go/lo"
	"github.com/iotaledger/hive.go/serializer/v2"
)

var (
	// ErrNonUniqueFeatures gets returned when multiple Feature(s) with the same FeatureType exist within sets.
	ErrNonUniqueFeatures = ierrors.New("non unique features within outputs")
	// ErrInvalidFeatureTransition gets returned when a Feature's transition within a ChainOutput is invalid.
	ErrInvalidFeatureTransition = ierrors.New("invalid feature transition")
	// ErrInvalidMetadataKey gets returned when a MetadataFeature's key is invalid.
	ErrInvalidMetadataKey = ierrors.New("invalid metadata key")
	// ErrInvalidStateMetadataKey gets returned when a StateMetadataFeature's key is invalid.
	ErrInvalidStateMetadataKey = ierrors.New("invalid state metadata key")
	// ErrMetadataExceedsMaxSize gets returned when a StateMetadataFeature or MetadataFeature exceeds the max size.
	ErrMetadataExceedsMaxSize = ierrors.New("metadata exceeds max allowed size")
)

// Feature is an abstract building block extending the features of an Output.
type Feature interface {
	Sizer
	NonEphemeralObject
	ProcessableObject
	constraints.Cloneable[Feature]
	constraints.Equalable[Feature]
	constraints.Comparable[Feature]

	// Type returns the type of the Feature.
	Type() FeatureType
}

// FeatureType defines the type of features.
type FeatureType byte

const (
	// FeatureSender denotes a SenderFeature.
	FeatureSender FeatureType = iota
	// FeatureIssuer denotes an IssuerFeature.
	FeatureIssuer
	// FeatureMetadata denotes a MetadataFeature.
	FeatureMetadata
	// FeatureStateMetadata denotes a StateMetadataFeature.
	FeatureStateMetadata
	// FeatureTag denotes a TagFeature.
	FeatureTag
	// NativeTokenFeature denotes a NativeTokenFeature.
	FeatureNativeToken
	// FeatureBlockIssuer denotes a BlockIssuerFeature.
	FeatureBlockIssuer
	// FeatureStaking denotes a StakingFeature.
	FeatureStaking
)

func (featType FeatureType) String() string {
	if int(featType) >= len(featNames) {
		return fmt.Sprintf("unknown feature type: %d", featType)
	}

	return featNames[featType]
}

var featNames = [FeatureStaking + 1]string{
	"SenderFeature",
	"IssuerFeature",
	"MetadataFeature",
	"StateMetadataFeature",
	"TagFeature",
	"NativeTokenFeature",
	"BlockIssuerFeature",
	"StakingFeature",
}

// Features is a slice of Feature(s).
type Features[T Feature] []T

// Clone clones the Features.
func (f Features[T]) Clone() Features[T] {
	cpy := make(Features[T], len(f))
	for i, v := range f {
		//nolint:forcetypeassert // we can safely assume that this is of type T
		cpy[i] = v.Clone().(T)
	}

	return cpy
}

func (f Features[T]) StorageScore(storageScoreStruct *StorageScoreStructure, _ StorageScoreFunc) StorageScore {
	var sumCost StorageScore
	for _, feat := range f {
		sumCost += feat.StorageScore(storageScoreStruct, nil)
	}

	return sumCost
}

func (f Features[T]) WorkScore(workScoreParameters *WorkScoreParameters) (WorkScore, error) {
	var workScoreFeats WorkScore
	for _, feat := range f {
		workScoreFeat, err := feat.WorkScore(workScoreParameters)
		if err != nil {
			return 0, err
		}

		workScoreFeats, err = workScoreFeats.Add(workScoreFeat)
		if err != nil {
			return 0, err
		}
	}

	return workScoreFeats, nil
}

func (f Features[T]) Size() int {
	sum := serializer.OneByte // 1 byte length prefix
	for _, feat := range f {
		sum += feat.Size()
	}

	return sum
}

// Set converts the slice into a FeatureSet.
// Returns an error if a FeatureType occurs multiple times.
func (f Features[T]) Set() (FeatureSet, error) {
	set := make(FeatureSet)
	for _, feat := range f {
		if _, has := set[feat.Type()]; has {
			return nil, ErrNonUniqueFeatures
		}
		set[feat.Type()] = feat
	}

	return set, nil
}

// MustSet works like Set but panics if an error occurs.
// This function is therefore only safe to be called when it is given,
// that a Features slice does not contain the same FeatureType multiple times.
func (f Features[T]) MustSet() FeatureSet {
	set, err := f.Set()
	if err != nil {
		panic(err)
	}

	return set
}

// Equal checks whether this slice is equal to other.
func (f Features[T]) Equal(other Features[T]) bool {
	if len(f) != len(other) {
		return false
	}

	for idx, feat := range f {
		if !feat.Equal(other[idx]) {
			return false
		}
	}

	return true
}

// Upsert adds the given feature or updates the previous one if existing.
func (f *Features[T]) Upsert(feature T) {
	for i, ele := range *f {
		if ele.Type() == feature.Type() {
			(*f)[i] = feature

			return
		}
	}
	*f = append(*f, feature)
}

// Remove removes the feature with the given type.
func (f *Features[T]) Remove(featureType FeatureType) bool {
	for i, ele := range *f {
		if ele.Type() == featureType {
			*f = append((*f)[:i], (*f)[i+1:]...)
			return true
		}
	}

	return false
}

// Sort sorts the Features in place by type.
func (f Features[T]) Sort() {
	sort.Slice(f, func(i, j int) bool { return f[i].Type() < f[j].Type() })
}

// FeatureSet is a set of Feature(s).
type FeatureSet map[FeatureType]Feature

// Clone clones the FeatureSet.
func (f FeatureSet) Clone() FeatureSet {
	return lo.CloneMap(f)
}

// SenderFeature returns the SenderFeature in the set or nil.
func (f FeatureSet) SenderFeature() *SenderFeature {
	b, has := f[FeatureSender]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a SenderFeature
	return b.(*SenderFeature)
}

// Issuer returns the IssuerFeature in the set or nil.
func (f FeatureSet) Issuer() *IssuerFeature {
	b, has := f[FeatureIssuer]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a IssuerFeature
	return b.(*IssuerFeature)
}

// Metadata returns the MetadataFeature in the set or nil.
func (f FeatureSet) Metadata() *MetadataFeature {
	b, has := f[FeatureMetadata]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a MetadataFeature
	return b.(*MetadataFeature)
}

// StateMetadata returns the StateMetadataFeature in the set or nil.
func (f FeatureSet) StateMetadata() *StateMetadataFeature {
	b, has := f[FeatureStateMetadata]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a StateMetadataFeature
	return b.(*StateMetadataFeature)
}

// Tag returns the TagFeature in the set or nil.
func (f FeatureSet) Tag() *TagFeature {
	b, has := f[FeatureTag]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a TagFeature
	return b.(*TagFeature)
}

// HasNativeTokenFeature tells whether this set has a FeatureNativeToken.
func (f FeatureSet) HasNativeTokenFeature() bool {
	_, has := f[FeatureNativeToken]
	return has
}

// NativeToken returns the NativeTokenFeature in the set or nil.
func (f FeatureSet) NativeToken() *NativeTokenFeature {
	b, has := f[FeatureNativeToken]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a NativeTokenFeature
	return b.(*NativeTokenFeature)
}

// BlockIssuer returns the BlockIssuerFeature in the set or nil.
func (f FeatureSet) BlockIssuer() *BlockIssuerFeature {
	b, has := f[FeatureBlockIssuer]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a BlockIssuerFeature
	return b.(*BlockIssuerFeature)
}

// Staking returns the StakingFeature in the set or nil.
func (f FeatureSet) Staking() *StakingFeature {
	b, has := f[FeatureStaking]
	if !has {
		return nil
	}

	//nolint:forcetypeassert // we can safely assume that this is a StakingFeature
	return b.(*StakingFeature)
}

// EveryTuple runs f for every key which exists in both this set and other.
// Returns a bool indicating whether all element of this set existed on the other set.
func (f FeatureSet) EveryTuple(other FeatureSet, fun func(a Feature, b Feature) error) (bool, error) {
	hadAll := true
	for ty, featA := range f {
		featB, has := other[ty]
		if !has {
			hadAll = false

			continue
		}
		if err := fun(featA, featB); err != nil {
			return false, err
		}
	}

	return hadAll, nil
}

// FeatureUnchanged checks whether the specified Feature type is unchanged between in and out.
// Unchanged also means that the block's existence is unchanged between both sets.
func FeatureUnchanged(featType FeatureType, inFeatSet FeatureSet, outFeatSet FeatureSet) error {
	in, inHas := inFeatSet[featType]
	out, outHas := outFeatSet[featType]

	switch {
	case outHas && !inHas:
		return ierrors.Wrapf(ErrInvalidFeatureTransition, "%s in next state but not in previous", featType)
	case !outHas && inHas:
		return ierrors.Wrapf(ErrInvalidFeatureTransition, "%s in current state but not in next", featType)
	}

	// not in both sets
	if in == nil {
		return nil
	}

	if !in.Equal(out) {
		return ierrors.Wrapf(ErrInvalidFeatureTransition, "%s changed, in %v / out %v", featType, in, out)
	}

	return nil
}

// checkPrintableASCIIString returns an error if the given string contains non-printable ASCII characters (including space).
func checkPrintableASCIIString(s string) error {
	for i := 0; i < len(s); i++ {
		if s[i] < 33 || s[i] > 126 {
			return ierrors.Errorf(
				"string contains non-printable ASCII character %d at index %d (allowed range 33 <= character <= 126)", s[i], i,
			)
		}
	}

	return nil
}
