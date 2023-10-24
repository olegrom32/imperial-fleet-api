package domain

const SpaceshipStatusDeleted = "deleted"

type Spaceship struct {
	ID       int64
	Name     string
	Class    string
	Armament []SpaceshipArmament
	Crew     int64
	Image    string
	Value    float64
	Status   string
}

func (s *Spaceship) Deleted() bool {
	return s.Status == SpaceshipStatusDeleted
}

type SpaceshipArmament struct {
	Armament Armament
	Quantity int64
}
