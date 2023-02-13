package dns

import (
	"strings"

	"github.com/asaskevich/govalidator"
	sliceutil "github.com/projectdiscovery/utils/slice"
	stringsutil "github.com/projectdiscovery/utils/strings"
	"golang.org/x/net/publicsuffix"
)

var ExcludeMap map[string]struct{}

// ChaosProgram json data item struct
type ChaosProgram struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Swag    bool     `json:"swag"`
	Domains []string `json:"domains"`
}

type ChaosList struct {
	Programs []ChaosProgram `json:"programs"`
}

func ValidateFQDN(value string) bool {
	// with publicsuffix package, it gets the top-level
	// domain of the input value
	tld, err := publicsuffix.EffectiveTLDPlusOne(value)
	if err != nil {
		return false
	}

	// with govalidator package, it checks if value provided is a
	// valid domain name system (DNS) name
	return tld == value && govalidator.IsDNSName(tld)
}

func ExtractHostname(item string) string {
	item = strings.ToLower(item)

	// Exclude if program name is in exclude.txt
	if _, ok := ExcludeMap[item]; ok {
		return ""
	}

	trimmedStr := stringsutil.TrimPrefixAny(item, "http://", "https://", "*.")

	if ValidateFQDN(trimmedStr) {
		return trimmedStr
	}

	return ""
}

func GetUniqueDomains(first, second []string) []string {
	_, diff := sliceutil.Diff(first, second)
	return sliceutil.Dedupe(diff)
}
