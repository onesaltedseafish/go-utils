// Package dhcp simuates how dhcp protocol allocates ip address
package dhcp

//go:generate mockgen -source=dhcp.go -destination=mock/dhcp_mock.go -package=mock
import (
	"encoding/binary"
	"errors"
	"math/big"
	"net"
	"sync"
)

var (
	// ErrHasNotEnoughAddr means can't allocate an address
	ErrHasNotEnoughAddr = errors.New("don't have an address to allocate")
)

// Storage interface defines how to store the content
type Storage interface {
	// GetAddressWithMAC if storage has a record of hardwareAddr
	// then return the related ip address
	// else return nil
	GetAddressWithMAC(net.HardwareAddr) (net.IP, error)
	// GetOneUnusedAddress finds the first unused record
	GetOneUnusedAddress() (net.IP, error)
	// GetLastAddress finds the last used ip address
	// if no ip address was used, return the first you want use
	GetLastAddress() (net.IP, error)
	// SetAddressWithMAC sets record with ip address and MAC address
	SetAddressWithMAC(net.IP, net.HardwareAddr) error
	// ReleaseAddress release the address
	ReleaseAddress(net.IP) error
	// IsUsed judge the ip address is used or not
	IsUsed(net.IP) (bool, error)
}

// Client monitor how dhcp server work
// concurrency safty
type Client struct {
	network net.IPNet // 分配的子网
	storage Storage
	mu      sync.Mutex
}

// New a dhcp client
func New(network net.IPNet, storage Storage) *Client {
	return &Client{
		network: network,
		storage: storage,
	}
}

// AllocateAddress allocate an address for a MAC address
// return an error if no addresses can be allocated
func (cli *Client) AllocateAddress(hwAddr net.HardwareAddr) (net.IP, error) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	// 1. Got a pre used address if possible
	ip, err := cli.storage.GetAddressWithMAC(hwAddr)
	if err != nil {
		return nil, err
	}
	if ip != nil {
		err = cli.storage.SetAddressWithMAC(ip, hwAddr)
		return ip, err
	}
	// 2. Try get the last used address
	// and try to allocate 1 address after that addr
	lastIp, err := cli.storage.GetLastAddress()
	if err != nil {
		return nil, err
	}
	for {
		newIp := IpAdd(lastIp, 1)
		lastIp = newIp
		if cli.network.Contains(newIp) {
			if used, _ := cli.storage.IsUsed(newIp); used {
				continue
			}
			err = cli.storage.SetAddressWithMAC(newIp, hwAddr)
			return newIp, err
		} else {
			break
		}
	}
	// 3. Try to find a record which is not used anymore
	// and allocate that addr to current MAC
	ip, err = cli.storage.GetOneUnusedAddress()
	if err != nil {
		return nil, err
	}
	if ip != nil {
		err = cli.storage.SetAddressWithMAC(ip, hwAddr)
		return ip, err
	}
	// 4. Can't allocate addr return error
	return nil, ErrHasNotEnoughAddr
}

// ReleaseAddress release the ip address
func (cli *Client) ReleaseAddress(ip net.IP) error {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	return cli.storage.ReleaseAddress(ip)
}

// IpAdd add delta to an ip and get a new ip address
func IpAdd(ip net.IP, delta int) net.IP {
	if len(ip) == net.IPv4len {
		return ipv4Add(ip, delta)
	}
	return ipv6Add(ip, delta)
}

func ipv4Add(ip net.IP, delta int) net.IP {
	addr := ip.To4()
	result := make(net.IP, 4)
	binary.BigEndian.PutUint32(result, binary.BigEndian.Uint32(addr)+uint32(delta))
	return result
}

func ipv6Add(ip net.IP, delta int) net.IP {
	addr := ip.To16()
	ipInt := new(big.Int).SetBytes(addr) // big-endian
	ipInt = ipInt.Add(ipInt, big.NewInt(int64(delta)))

	result := make(net.IP, net.IPv6len)
	ipInt.FillBytes(result)

	return result
}
