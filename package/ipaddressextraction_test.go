package p

import (
	"testing"
	"strings"
	"net/http/httptest"
	"io/ioutil"
)

type addressesNoProxy struct {
	sourceAddress string
	requiredAddress string
}

type addressesWithProxy struct {
	addressFromSourceHost string
	addressFromProxyHeader string
	requiredAddress string
	proxyHeaderName string
}

func getAddressesNoProxy() []addressesNoProxy {
	addresses := []addressesNoProxy {
		{"1.2.3.4:1234","1.2.3.4"},
		{"5.6.7.8:1134","5.6.7.8"},
		{"[::1]:5678","::1"},
		{"[2003:11::ed]:3443","2003:11::ed"},
	}

	return addresses
}

func getAddressesWithProxy() []addressesWithProxy {
	addresses := []addressesWithProxy {
		{"1.2.3.4:1234","5.6.7.8","5.6.7.8","X-Forwarded-For"},
		{"172.16.0.18:1134","1.2.3.4","1.2.3.4","X-fOrWaRdEd-FoR"},
		{"1.2.3.4:1234","5.6.7.8","1.2.3.4","Some-Random-Header"},
		{"[::1]:5678","2003:11ed::1","2003:11ed::1","X-Forwarded-For"},
		{"[2003:11::ed]:3443","2003:11::eb","2003:11::eb","X-Forwarded-For"},
		{"[2003:11::ed]:3443","2003:11::eb","2003:11::ed","Some-Random-Header"},
	}

	return addresses
}

func TestGetClientsIPAddressNoProxyHeader(t *testing.T) {
	for _, table := range getAddressesNoProxy() {
		fakeRequest := httptest.NewRequest("GET", "http://1.2.3.4", nil)
		fakeRequest.RemoteAddr = table.sourceAddress
		ipAddress := GetClientsIPAddress(fakeRequest)
		if ipAddress != table.requiredAddress {
			t.Errorf("Requested from %v but got %v.",table.requiredAddress, ipAddress)
		}
	}
}

func TestGetClientsIPAddressWithProxyHeaders(t *testing.T) {
	for _, table := range getAddressesWithProxy() {
		fakeRequest := httptest.NewRequest("GET", "http://1.2.3.4", nil)
		fakeRequest.RemoteAddr = table.addressFromSourceHost
		fakeRequest.Header.Set(table.proxyHeaderName, table.addressFromProxyHeader)
		ipAddress := GetClientsIPAddress(fakeRequest)
		if ipAddress != table.requiredAddress {
			t.Errorf("Requested from %v but got %v.",table.requiredAddress, ipAddress)
		}
	}
}

func TestRespondWithPublicIPAddressNoProxyHeader(t *testing.T) {
	for _, table := range getAddressesNoProxy() {
		fakeRequest := httptest.NewRequest("GET", "http://1.2.3.4", nil)
		fakeRequest.RemoteAddr = table.sourceAddress
		w := httptest.NewRecorder()
		RespondWithPublicIPAddress(w, fakeRequest)
		response := w.Result()
		body,_ := ioutil.ReadAll(response.Body)
		if response.StatusCode != 200 {
			t.Errorf("Requested from %v but got not OK code %v.",table.requiredAddress, response.StatusCode)
		}
		if strings.TrimSpace(string(body)) != table.requiredAddress {
			t.Errorf("Requested from %v but got %v.",table.requiredAddress, strings.TrimSpace(string(body)))
		}
	}
}

func TestRespondWithPublicIPAddressWithProxyHeader(t *testing.T) {
	for _, table := range getAddressesWithProxy() {
		fakeRequest := httptest.NewRequest("GET", "http://1.2.3.4", nil)
		fakeRequest.RemoteAddr = table.addressFromSourceHost
		fakeRequest.Header.Set(table.proxyHeaderName, table.addressFromProxyHeader)
		w := httptest.NewRecorder()
		RespondWithPublicIPAddress(w, fakeRequest)
		response := w.Result()
		body,_ := ioutil.ReadAll(response.Body)
		if response.StatusCode != 200 {
			t.Errorf("Requested from %v but got not OK code %v.",table.requiredAddress, response.StatusCode)
		}
		if strings.TrimSpace(string(body)) != table.requiredAddress {
			t.Errorf("Requested from %v but got %v.",table.requiredAddress, strings.TrimSpace(string(body)))
		}
	}
}
