package kibana

import "testing"

func TestClient(t *testing.T) {
	client := NewClient("dan.chen", "nG5=7$D6", "https://172.17.6.203")
	client.Login()
	client.DataViews()
}
