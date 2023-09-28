package iotago

import (
	"bytes"

	"github.com/iotaledger/hive.go/serializer/v2"
)

// TagFeature is a feature which allows to additionally tag an output by a user defined value.
type TagFeature struct {
	Tag []byte `serix:"0,lengthPrefixType=uint8,mapKey=tag,minLen=1,maxLen=64"`
}

func (s *TagFeature) Clone() Feature {
	return &TagFeature{Tag: append([]byte(nil), s.Tag...)}
}

func (s *TagFeature) VBytes(rentStruct *RentStructure, f VBytesFunc) VBytes {
	if f != nil {
		return f(rentStruct)
	}

	return rentStruct.VBFactorData().Multiply(VBytes(serializer.SmallTypeDenotationByteSize + serializer.OneByte + len(s.Tag)))
}

func (s *TagFeature) WorkScore(_ *WorkScoreStructure) (WorkScore, error) {
	return 0, nil
}

func (s *TagFeature) Equal(other Feature) bool {
	otherFeat, is := other.(*TagFeature)
	if !is {
		return false
	}

	return bytes.Equal(s.Tag, otherFeat.Tag)
}

func (s *TagFeature) Type() FeatureType {
	return FeatureTag
}

func (s *TagFeature) Size() int {
	// FeatureType + Tag
	return serializer.SmallTypeDenotationByteSize + serializer.OneByte + len(s.Tag)
}
