//pkg/aviation/gribwind.go

package aviation

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/mmp/vice/pkg/math"
)

// GribWindModel implements WindModel using GRIB-like JSON data.
// The JSON is expected to look like:
// {"points":[{"lat":..,"lon":..,"levels":[{"h":..,"u":..,"v":..},...]},...]}

type gribPointLevel struct {
	Loc math.Point2LL
	U   float32
	V   float32
}

// GribWindModel stores wind vectors grouped by altitude level.
type GribWindModel struct {
	Levels  map[float32][]gribPointLevel // level height in meters -> points
	Heights []float32                    // sorted list of level heights
}

// LoadGribWindModel loads a GribWindModel from the given filename.
func LoadGribWindModel(filename string) (*GribWindModel, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var raw struct {
		Points []struct {
			Lat    float32 `json:"lat"`
			Lon    float32 `json:"lon"`
			Levels []struct {
				H float32 `json:"h"`
				U float32 `json:"u"`
				V float32 `json:"v"`
			} `json:"levels"`
		} `json:"points"`
	}

	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, err
	}

	m := &GribWindModel{Levels: make(map[float32][]gribPointLevel)}
	for _, p := range raw.Points {
		loc := math.Point2LL{p.Lon, p.Lat}
		for _, l := range p.Levels {
			m.Levels[l.H] = append(m.Levels[l.H], gribPointLevel{Loc: loc, U: l.U, V: l.V})
		}
	}

	for h := range m.Levels {
		m.Heights = append(m.Heights, h)
	}
	sort.Slice(m.Heights, func(i, j int) bool { return m.Heights[i] < m.Heights[j] })

	return m, nil
}

const metersToNM = 1.0 / 1852.0

func (m *GribWindModel) interpolateLevel(p math.Point2LL, h float32) (u, v float32) {
	pts := m.Levels[h]
	if len(pts) == 0 {
		return 0, 0
	}
	var sumW, sumU, sumV float32
	for _, pt := range pts {
		d := math.NMDistance2LL(p, pt.Loc)
		if d == 0 {
			return pt.U, pt.V
		}
		w := 1 / (d * d)
		sumW += w
		sumU += pt.U * w
		sumV += pt.V * w
	}
	if sumW == 0 {
		return 0, 0
	}
	return sumU / sumW, sumV / sumW
}

// GetWindVector returns the wind vector at the given point and altitude.
func (m *GribWindModel) GetWindVector(p math.Point2LL, alt float32) [2]float32 {
	if len(m.Heights) == 0 {
		return [2]float32{}
	}

	altm := alt * 0.3048
	// Clamp altitude to available range
	if altm <= m.Heights[0] {
		u, v := m.interpolateLevel(p, m.Heights[0])
		// The GRIB data gives the velocity of the wind in the
		// standard meteorological orientation, where positive values
		// indicate flow toward the east/north. The returned vector
		// should represent the effect of that wind on the aircraft,
		// so pass it through without negation.
		return [2]float32{u * metersToNM, v * metersToNM}
	}
	last := m.Heights[len(m.Heights)-1]
	if altm >= last {
		u, v := m.interpolateLevel(p, last)
		return [2]float32{u * metersToNM, v * metersToNM}
	}

	// Find surrounding levels
	i := sort.Search(len(m.Heights), func(i int) bool { return m.Heights[i] > altm })
	h0 := m.Heights[i-1]
	h1 := m.Heights[i]

	u0, v0 := m.interpolateLevel(p, h0)
	u1, v1 := m.interpolateLevel(p, h1)

	f := (altm - h0) / (h1 - h0)
	u := math.Lerp(f, u0, u1)
	v := math.Lerp(f, v0, v1)
	return [2]float32{u * metersToNM, v * metersToNM}
}

var _ WindModel = (*GribWindModel)(nil)
