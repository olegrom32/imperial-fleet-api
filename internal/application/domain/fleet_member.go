package domain

// FleetMember is the member of the fleet. The entity is created for simple auth.
// Password is stored unencrypted for simplicity only, prod-ready version of the API must have
// passwords stored in encrypted form.
type FleetMember struct {
	Login    string
	Password string
}
