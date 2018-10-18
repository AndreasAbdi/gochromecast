package receiver

//GetSessionByNamespace attempts to return the first session with a specified namespace.
func (s *Status) GetSessionByNamespace(namespace string) *ApplicationSession {

	for _, app := range s.Applications {
		for _, ns := range app.Namespaces {
			if ns.Name == namespace {
				return app
			}
		}
	}

	return nil
}
