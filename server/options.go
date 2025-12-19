package server

type serverBuilder struct {
	bind      string
	port      int
	databases int

	hz                  int
	activeExpireSamples int
	requirePass         string
}

type ServerOption func(*serverBuilder)

func WithBind(bind string) ServerOption {
	return func(sb *serverBuilder) {
		sb.bind = bind
	}
}

func WithPort(port int) ServerOption {
	return func(sb *serverBuilder) {
		sb.port = port
	}
}

func WithDatabases(num int) ServerOption {
	return func(sb *serverBuilder) {
		sb.databases = num
	}
}

func WithHZ(hz int) ServerOption {
	return func(sb *serverBuilder) {
		sb.hz = hz
	}
}

func WithActiveExpireSamples(samples int) ServerOption {
	return func(sb *serverBuilder) {
		sb.activeExpireSamples = samples
	}
}

func WithRequirePass(password string) ServerOption {
	return func(sb *serverBuilder) {
		sb.requirePass = password
	}
}
