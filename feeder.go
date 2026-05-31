package feeder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/feeder", new(Feeder))
}

type Feeder struct{}

var client = &http.Client{Timeout: 2 * time.Second}

func base() string {
	b := os.Getenv("FEEDER_URL")
	if b == "" {
		return "http://feeder:9000"
	}
	return b
}

func (f *Feeder) Get(index int64) []string {
	url := fmt.Sprintf("%s/row?index=%d", base(), index)
	resp, _ := client.Get(url)
	defer resp.Body.Close()

	var row []string
	json.NewDecoder(resp.Body).Decode(&row)
	return row
}

func (f *Feeder) LookupAssignments(email string) []int {
	url := fmt.Sprintf("%s/lookup/assignments?email=%s", base(), email)
	resp, _ := client.Get(url)
	defer resp.Body.Close()

	var out map[string][]int
	json.NewDecoder(resp.Body).Decode(&out)
	return out["assignments"]
}
