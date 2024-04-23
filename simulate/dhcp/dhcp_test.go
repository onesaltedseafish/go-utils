package dhcp_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/onesaltedseafish/go-utils/simulate/dhcp"
	"github.com/stretchr/testify/assert"
)

var _ dhcp.Storage = (*testStorageImpl)(nil)

func TestIpAdd(t *testing.T) {
	testcases := []struct {
		Ori   string
		delta int
		Want  string
	}{
		{"192.168.1.1", 1, "192.168.1.2"}, // ipv4
		{"192.168.1.1", 11, "192.168.1.12"},
		{"192.168.1.1", 256, "192.168.2.1"},
		{"::1", 1, "::2"}, // ipv6
	}

	for _, testcase := range testcases {
		ip := net.ParseIP(testcase.Ori)
		assert.NotEqual(t, nil, ip)
		newIp := dhcp.IpAdd(ip, testcase.delta)
		assert.Equal(t, testcase.Want, newIp.String())
	}
}

type testStorageImpl struct {
	m    map[string]string // map[macAddr] ipAddr
	used map[string]bool   // map[ipAddr] bool
	last net.IP
}

func newTestStorageImpl(network net.IPNet) *testStorageImpl {
	return &testStorageImpl{
		m:    make(map[string]string),
		used: make(map[string]bool),
		last: network.IP,
	}
}

// GetAddressWithMAC if storage has a record of hardwareAddr
// then return the related ip address
// else return nil
func (s *testStorageImpl) GetAddressWithMAC(addr net.HardwareAddr) net.IP {
	return net.ParseIP(s.m[addr.String()])
}

// GetOneUnusedAddress finds the first unused record
func (s *testStorageImpl) GetOneUnusedAddress() net.IP {
	for k, v := range s.used {
		if !v {
			return net.ParseIP(k)
		}
	}
	return nil
}

// GetLastAddress finds the last used ip address
func (s *testStorageImpl) GetLastAddress() net.IP {
	return s.last
}

// SetAddressWithMAC sets record with ip address and MAC address
func (s *testStorageImpl) SetAddressWithMAC(ip net.IP, mac net.HardwareAddr) {
	s.m[mac.String()] = ip.String()
	s.used[ip.String()] = true
	s.last = ip
}

// ReleaseAddress release the address
func (s *testStorageImpl) ReleaseAddress(ip net.IP) error {
	if _, ok := s.used[ip.String()]; ok {
		s.used[ip.String()] = false
		return nil
	}
	return fmt.Errorf("no result with %s", ip.String())
}

// IsUsed judge the ip address is used or not
func (s *testStorageImpl) IsUsed(ip net.IP) bool {
	var r, ok bool
	r, ok = s.used[ip.String()]
	if ok {
		return r
	}
	return false
}

func TestDhcpClient(t *testing.T) {
	var ip net.IP
	var err error
	var (
		_, network, _ = net.ParseCIDR("127.0.0.1/30")
		mac1, _       = net.ParseMAC("00:16:3e:03:57:45")
		mac2, _       = net.ParseMAC("02:42:be:7f:b3:58")
		mac3, _       = net.ParseMAC("02:42:fe:21:ad:e3")
		mac4, _       = net.ParseMAC("9a:7e:a6:2f:f0:d0")
		mac5, _       = net.ParseMAC("00:16:3e:03:57:46")
	)

	dhcpClient := dhcp.New(*network, newTestStorageImpl(*network))

	// allocate from start
	ip, err = dhcpClient.AllocateAddress(mac1)
	assert.Equal(t, nil, err)
	assert.Equal(t, "127.0.0.1", ip.String())

	// release and allocate the same
	err = dhcpClient.ReleaseAddress(net.ParseIP("127.0.0.1"))
	assert.Equal(t, nil, err)
	ip, err = dhcpClient.AllocateAddress(mac1)
	assert.Equal(t, nil, err)
	assert.Equal(t, "127.0.0.1", ip.String())
	// allocate more 2 address
	_, err = dhcpClient.AllocateAddress(mac2)
	assert.Equal(t, nil, err)
	_, err = dhcpClient.AllocateAddress(mac3)
	assert.Equal(t, nil, err)
	// release one and alloate to mac4
	err = dhcpClient.ReleaseAddress(net.ParseIP("127.0.0.2"))
	assert.Equal(t, nil, err)
	ip, err = dhcpClient.AllocateAddress(mac4)
	assert.Equal(t, nil, err)
	assert.Equal(t, "127.0.0.2", ip.String())
	// can't allocate address for mac5
	_, err = dhcpClient.AllocateAddress(mac5)
	assert.NotEqual(t, nil, err)
}
