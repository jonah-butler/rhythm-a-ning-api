package request

import (
	"rhythmapi/model"
	"testing"
)

func polyRhythm() model.RhythmInputPoly {
	polyBeats := 3
	polySub := model.BASE

	return model.RhythmInputPoly{
		RhythmInputMono: model.RhythmInputMono{
			RhythmDiscriminator: model.RhythmDiscriminator{RhythmType: model.Polyrhythm},
			Beats:               4,
			Subdivisions:        model.Subdivisions{model.BASE, model.TRIPLET, model.DUPLET, model.BASE},
		},
		PolyBeats:        &polyBeats,
		PolyState:        []int64{1, 0, 0},
		PolySubdivision:  &polySub,
		PolySubdivisions: model.Subdivisions{model.BASE, model.BASE, model.QUINTUPLET},
	}
}

func TestValidateRhythm(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(r *model.RhythmInputPoly)
		wantErr bool
	}{
		{"valid poly", func(r *model.RhythmInputPoly) {}, false},
		{"valid mono ignores poly fields", func(r *model.RhythmInputPoly) {
			r.RhythmType = model.Monorhythm
			r.PolySubdivisions = nil
		}, false},
		{"too few subdivisions", func(r *model.RhythmInputPoly) {
			r.Subdivisions = r.Subdivisions[:3]
		}, true},
		{"too many subdivisions", func(r *model.RhythmInputPoly) {
			r.Subdivisions = append(r.Subdivisions, model.BASE)
		}, true},
		{"unknown subdivision", func(r *model.RhythmInputPoly) {
			r.Subdivisions[1] = "sixteenth"
		}, true},
		{"missing poly subdivisions", func(r *model.RhythmInputPoly) {
			r.PolySubdivisions = nil
		}, true},
		{"poly subdivisions don't match poly beats", func(r *model.RhythmInputPoly) {
			r.PolySubdivisions = r.PolySubdivisions[:2]
		}, true},
		{"unknown poly subdivision", func(r *model.RhythmInputPoly) {
			r.PolySubdivisions[0] = ""
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := polyRhythm()
			tt.mutate(&r)

			err := ValidateRhythm(r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateRhythm() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
