package server

type serverBuilder struct {
	host        string
	port        int
	databaseNum int
}

type ServerOption func(*serverBuilder)

func WithServerHost(host string) ServerOption {
	return func(sb *serverBuilder) {
		sb.host = host
	}
}

func WithServerPort(port int) ServerOption {
	return func(sb *serverBuilder) {
		sb.port = port
	}
}

func WithDatabaseNum(num int) ServerOption {
	return func(sb *serverBuilder) {
		sb.databaseNum = num
	}
}
