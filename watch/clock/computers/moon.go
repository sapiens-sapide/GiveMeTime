package computers

import (
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro"
	"github.com/soniakeys/meeus/base"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/moonposition"
	"github.com/soniakeys/meeus/nutation"
	"github.com/soniakeys/meeus/parallax"
	"github.com/soniakeys/unit"
)

func TopoMoonPosition(jde float64, pos globe.Coord) (α unit.RA, δ unit.Angle, Δ float64) {
	λ, β, Δ := moonposition.Position(jde) // (λ without nutation)
	//Δψ, Δε := nutation.Nutation(jde)
	obl := coord.NewObliquity(nutation.MeanObliquityLaskar(jde))
	α, δ = coord.EclToEq(λ, β, obl.S, obl.C)
	E := globe.Earth76
	ρsφʹ, ρcφʹ := E.ParallaxConstants(astro.Position.Lat, E.RadiusAtLatitude(astro.Position.Lat)*1000)
	α, δ = parallax.Topocentric(α, δ, Δ/base.AU, ρsφʹ, ρcφʹ, astro.Position.Lon, jde)
	return
}
