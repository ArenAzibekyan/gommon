package httpserv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddr(t *testing.T) {
	type args struct {
		host string
		port uint16
	}
	type want struct {
		addr string
	}
	cases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			args: args{},
			want: want{},
		},
		{
			name: "port only",
			args: args{
				port: 8765,
			},
			want: want{
				addr: ":8765",
			},
		},
		{
			name: "ip host only",
			args: args{
				host: "10.20.30.40",
			},
			want: want{
				addr: "10.20.30.40",
			},
		},
		{
			name: "ip host and port",
			args: args{
				host: "10.20.30.40",
				port: 8765,
			},
			want: want{
				addr: "10.20.30.40:8765",
			},
		},
		{
			name: "localhost and port",
			args: args{
				host: "localhost",
				port: 8765,
			},
			want: want{
				addr: "localhost:8765",
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			addr := Addr(c.args.host, c.args.port)
			require.Equal(t, c.want.addr, addr)
		})
	}
}
