package model

import (
	"net"
	"reflect"
	"testing"
)

var (
	testIP    = [4]uint8{1, 2, 3, 4}
	testPort  = uint16(1234)
	testIPStr = "1.2.3.4"
)

func TestUDPAddrToAddressPort(t *testing.T) {
	udpAddr, _ := net.ResolveUDPAddr("udp", testIPStr+":1234")
	got := UDPAddrToAddressPort(*udpAddr)

	t.Run("check_port", func(t *testing.T) {
		want := testPort
		if got.port != want {
			t.Errorf("UDPAddrToAddressPort().port = %v, want %v", got.port, want)
			return
		}
	})

	t.Run("check_ip", func(t *testing.T) {
		want := testIP
		for i := range want {
			if got.ip[i] != want[i] {
				t.Errorf("UDPAddrToAddressPort().ip = %v, want %v", got.ip, want)
				return
			}
		}
	})
}

func TestAddressPort_String(t *testing.T) {
	const want = "1.2.3.4:1234"
	a := &AddressPort{
		port: testPort,
		ip:   testIP,
	}
	if got := a.String(); got != want {
		t.Errorf("AddressPort.String() = %v, want %v", got, want)
	}
}

func TestAddressPort_getIP(t *testing.T) {
	a := &AddressPort{
		port: testPort,
		ip:   testIP,
	}
	want := net.ParseIP(testIPStr)
	if got := a.getIP(); !reflect.DeepEqual(got, want) {
		t.Errorf("AddressPort.getIP() = %v, want %v", got, want)
	}
}

func TestAddressPort_getPort(t *testing.T) {
	a := &AddressPort{
		port: testPort,
		ip:   testIP,
	}
	want := int(testPort)
	if got := a.getPort(); got != want {
		t.Errorf("AddressPort.getPort() = %v, want %v", got, want)
	}
}
