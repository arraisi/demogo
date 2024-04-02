package mongo

func GetRegisteredCollection() map[string]func() interface{} {
	registeredCollections := make(map[string]func() interface{})

	audit := AuditLog{}

	registeredCollections[audit.CollectionName()] = func() interface{} {
		return &audit
	}

	return registeredCollections
}
