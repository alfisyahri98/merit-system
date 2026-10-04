package helper

import (
	"regexp"
	"strings"
)

var (
	semuaBesar = regexp.MustCompile(`^[A-Z0-9-]{2,}$`)
	romawi     = regexp.MustCompile(`^(I|II|III|IV|V|VI|VII|VIII|IX|X|XI|XII)$`)
	singkatan  = map[string]bool{"SPKT": true, "SPN": true, "SSDM": true, "SDM": true, "PNS": true, "TNI": true}
)

// Rapikan: "SUBBIDTEKINFO BIDTIK" → "Subbidtekinfo Bidtik" (sama dengan tampilan frontend).
// Angka romawi dan singkatan baku tetap huruf besar.
func Rapikan(teks string) string {
	kata := strings.Fields(teks)
	for i, k := range kata {
		if !semuaBesar.MatchString(k) || strings.ToLower(k) == k || romawi.MatchString(k) || singkatan[k] {
			continue
		}
		if k == "DAN" {
			kata[i] = "dan"
			continue
		}
		bagian := strings.Split(k, "-")
		for j, b := range bagian {
			if b != "" {
				bagian[j] = b[:1] + strings.ToLower(b[1:])
			}
		}
		kata[i] = strings.Join(bagian, "-")
	}
	return strings.Join(kata, " ")
}
