package web

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

func (s *Config) WithPrefix(v string) *Config {
	s.Prefix = v
	return s
}
