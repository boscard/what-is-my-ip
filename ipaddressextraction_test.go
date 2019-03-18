package main

import (
	"testing"
	"net/http/httptest"
)


func TestGetClientsIPAddressNoProxyHeader(t *testing.T) {
	tables := []struct {
		sourceAddress string
		parsedAddress string
	}{
		{"1.2.3.4:1234","1.2.3.4"},
		{"5.6.7.8:1134","5.6.7.8"},
		{"[::1]:5678","::1"},
		{"[2003:11::ed]:3443","2003:11::ed"},
	}
	for _, table := range tables {
		fakeRequest := httptest.NewRequest("GET", "http://1.2.3.4", nil)
		fakeRequest.RemoteAddr = table.sourceAddress
		ipAddress := GetClientsIPAddress(fakeRequest)
		if ipAddress != table.parsedAddress {
			t.Errorf("Requested from %v but got %v.",table.parsedAddress, ipAddress)
		}
	}
}

func TestGetClientsIPAddressWithProxyHeaders(t *testing.T) {
	tables := []struct {
		addressFromSourceHost string
		addressFromProxyHeader string
		requiredAddress string
		proxyHeaderName string
	}{
		{"1.2.3.4:1234","5.6.7.8","5.6.7.8","X-Forwarded-For"},
		{"172.16.0.18:1134","1.2.3.4","1.2.3.4","X-fOrWaRdEd-FoR"},
		{"1.2.3.4:1234","5.6.7.8","1.2.3.4","Some-Random-Header"},
		{"[::1]:5678","2003:11ed::1","2003:11ed::1","X-Forwarded-For"},
		{"[2003:11::ed]:3443","2003:11::eb","2003:11::eb","X-Forwarded-For"},
		{"[2003:11::ed]:3443","2003:11::eb","2003:11::ed","Some-Random-Header"},
	}
	for _, table := range tables {
		fakeRequest := httptest.NewRequest("GET", "http://1.2.3.4", nil)
		fakeRequest.RemoteAddr = table.addressFromSourceHost
		fakeRequest.Header.Set(table.proxyHeaderName, table.addressFromProxyHeader)
		ipAddress := GetClientsIPAddress(fakeRequest)
		if ipAddress != table.requiredAddress {
			t.Errorf("Requested from %v but got %v.",table.requiredAddress, ipAddress)
		}
	}
}
