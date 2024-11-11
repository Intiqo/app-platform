package secrets

// Manager defines methods for getting secrets
type Manager interface {
	// GetSecret returns the secret value for the given name
	GetSecret(name string) (result string, err error)
}
