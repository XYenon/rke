package cluster

import (
	"testing"

	"github.com/rancher/rke/hosts"
	v3 "github.com/rancher/rke/types"
	"github.com/stretchr/testify/assert"
)

func Test_getUniqStringList(t *testing.T) {
	type args struct {
		l []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"contain strings with only spaces",
			args{
				[]string{" ", "key1=value1", "   ", "key2=value2"},
			},
			[]string{"key1=value1", "key2=value2"},
		},
		{
			"contain strings with trailing or leading spaces",
			args{
				[]string{"  key1=value1", "key1=value1  ", "  key2=value2   "},
			},
			[]string{"key1=value1", "key2=value2"},
		},
		{
			"contain duplicated strings",
			args{
				[]string{"", "key1=value1", "key1=value1", "key2=value2"},
			},
			[]string{"key1=value1", "key2=value2"},
		},
		{
			"contain empty string",
			args{
				[]string{"", "key1=value1", "", "key2=value2"},
			},
			[]string{"key1=value1", "key2=value2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getUniqStringList(tt.args.l), "getUniqStringList(%v)", tt.args.l)
		})
	}
}

func Test_getEtcdListenAddress(t *testing.T) {
	tests := []struct {
		name     string
		host     *hosts.Host
		expected string
	}{
		{
			name: "ipv4 same address binds all ipv4",
			host: &hosts.Host{
				RKEConfigNode: v3.RKEConfigNode{
					Address:         "192.168.1.1",
					InternalAddress: "192.168.1.1",
				},
			},
			expected: etcdListenAllIPv4,
		},
		{
			name: "ipv6 same address binds all ipv6",
			host: &hosts.Host{
				RKEConfigNode: v3.RKEConfigNode{
					Address:         "2001:db8::1",
					InternalAddress: "2001:db8::1",
				},
			},
			expected: etcdListenAllIPv6,
		},
		{
			name: "different internal address uses internal address",
			host: &hosts.Host{
				RKEConfigNode: v3.RKEConfigNode{
					Address:         "203.0.113.1",
					InternalAddress: "10.0.0.5",
				},
			},
			expected: "10.0.0.5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, getEtcdListenAddress(tc.host))
		})
	}
}

func Test_getEtcdctlLocalEndpoint(t *testing.T) {
	tests := []struct {
		listenAddress string
		expected      string
	}{
		{etcdListenAllIPv4, "127.0.0.1:2379"},
		{etcdListenAllIPv6, "[::1]:2379"},
		{"2001:db8::1", "[2001:db8::1]:2379"},
	}

	for _, tc := range tests {
		t.Run(tc.listenAddress, func(t *testing.T) {
			assert.Equal(t, tc.expected, getEtcdctlLocalEndpoint(tc.listenAddress))
		})
	}
}
