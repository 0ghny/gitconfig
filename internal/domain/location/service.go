package location

// Service is the port that defines all location management operations.
// Consumers (e.g. CLI adapters) depend on this interface, never on the
// concrete implementation, keeping the hexagonal boundary intact.
type Service interface {
	GetLocations() ([]Location, error)
	FindLocationByKey(key string) (*Location, error)
	SaveLocation(key string, location string) error
	DeleteLocation(key string) error
}
