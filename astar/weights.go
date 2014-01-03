package astar

type WeightCalculation func(p *PathPoint, fill_weight int, target Point)

func RawDist(p *PathPoint, fill_weight int, target Point) {
    p.DistTraveled = p.Parent.DistTraveled + 1

    p.Weight = fill_weight
    p.Weight += p.Point.Dist(target)
    p.Weight += p.DistTraveled
}
