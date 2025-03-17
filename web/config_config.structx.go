package web

func (s *Config) Version() string { return s.version }

func (s *Config) Commit() string { return s.commit }

func (s *Config) WithVersion(v string) *Config {
	s.version = v
	return s
}

func (s *Config) WithCommit(v string) *Config {
	s.commit = v
	return s
}

func (s *Config) WithEndpoint(v string) *Config {
	s.Endpoint = v
	return s
}

func (s *Config) WithHost(v string) *Config {
	s.Host = v
	return s
}

func (s *Config) WithPort(v int) *Config {
	s.Port = v
	return s
}

func (s *Config) WithOrigins(v []string) *Config {
	s.Origins = v
	return s
}

func (s *Config) WithDebug(v bool) *Config {
	s.Debug = v
	return s
}

func (s *Config) WithPprof(v bool) *Config {
	s.Pprof = v
	return s
}

func (s *Config) WithRequestID(v bool) *Config {
	s.RequestID = v
	return s
}

func (s *Config) WithMetrics(v bool) *Config {
	s.Metrics = v
	return s
}

func (s *Config) WithSystem(v string) *Config {
	s.System = v
	return s
}

func (s *Config) WithPrefix(v string) *Config {
	s.Prefix = v
	return s
}

func (s *Config) WithHeaders(v []string) *Config {
	s.Headers = v
	return s
}

func (s *Config) WithAuth(v bool) *Config {
	s.Auth = v
	return s
}
