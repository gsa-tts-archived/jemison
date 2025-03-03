package searching

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/GSA-TTS/jemison/config"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"go.uber.org/zap"
)

// ////////// Searching.
func DoTheSearch(affiliate string, searchQuery string, JDB *postgres.JemisonDB) []SearchResult {
	domain64Start, domain64End := getD64(affiliate + ".gov")

	sri := SearchRequestInput{
		Host:          affiliate + ".gov",
		Path:          "",
		Terms:         searchQuery,
		Domain64Start: domain64Start,
		Domain64End:   domain64End,
	}

	rows, duration, err := runQuery(sri, JDB)

	zap.L().Info("Queried Answer:",
		zap.Any("rows: ", rows),
		zap.Any("duration", duration),
		zap.Any("err", err))

	return rows
}

func getD64(affiliate string) (string, string) {
	var subdomain, domain, tld string

	subdomain, domain, tld = parseAffiliate(affiliate)

	var d64Start, d64End int64

	// top level domain
	d64Start, _ = strconv.ParseInt(fmt.Sprintf("%02x00000000000000", tld), 16, 64)
	d64End, _ = strconv.ParseInt(fmt.Sprintf("%02xFFFFFFFFFFFF00", tld), 16, 64)

	// domain
	if domain != "" {
		start := config.RDomainToDomain64(fmt.Sprintf("%s.%s", tld, domain))
		d64Start, _ = strconv.ParseInt(fmt.Sprintf("%s00000000", start), 16, 64)
		d64End, _ = strconv.ParseInt(fmt.Sprintf("%sFFFFFF00", start), 16, 64)
	} else {
		sD64Start := fmt.Sprintf("%d", d64Start)
		sD64End := fmt.Sprintf("%d", d64End)

		return sD64Start, sD64End
	}

	// subdomain
	if subdomain != "" {
		fqdn := fmt.Sprintf("%s.%s.%s", subdomain, domain, tld)
		start, _ := config.FQDNToDomain64(fqdn)
		d64Start = start
		d64End = start + 1
	}

	sD64Start := fmt.Sprintf("%d", d64Start)
	sD64End := fmt.Sprintf("%d", d64End)

	return sD64Start, sD64End
}

func parseAffiliate(affiliate string) (string, string, string) {
	tld := ""
	domain := ""
	subdomain := ""
	delimiter := "."

	results := strings.Split(affiliate, delimiter)

	// if it has length of 3 it has a subdomain, length of 2 only a domain, length of 1 only a tld
	//nolint:mnd
	if len(results) == 3 {
		subdomain = results[0]
		domain = results[1]
		tld = results[2]
	} else if len(results) == 2 {
		domain = results[0]
		tld = results[1]
	} else {
		tld = results[0]
	}

	return subdomain, domain, tld
}

////////////////////

// ////////// Returning Results.
func ParseTheResults(results []SearchResult, reqQP RequiredQueryParameters, optQP OptionalQueryParameters) []byte {
	// create array of results {JSONResults}
	jSONResults := createJSONResults(results)

	// create webJSON
	webResults := createWebResults(jSONResults, optQP)

	// create wholeJSON
	wholeJSON := createWholeJSON(webResults, reqQP)

	return wholeJSON
}

func createJSONResults(results []SearchResult) []QueryWebResultsData {
	//nolint:prealloc
	var JSONResults []QueryWebResultsData

	for _, r := range results {
		// convert searchresult into a json object that matches current resultAPI
		qwrd := createQueryWebResultsData(r)

		// append to JSONResults
		JSONResults = append(JSONResults, qwrd)
	}

	return JSONResults
}

func createWebResults(jSONResults []QueryWebResultsData, optionalQueryParams OptionalQueryParameters) QueryWebData {
	total := 5
	nextOffset := optionalQueryParams.Offset
	spellingCorrections := "null"
	strc := QueryWebData{total, nextOffset, spellingCorrections, jSONResults}

	return strc
}

func createWholeJSON(webResults QueryWebData, requiredQueryParams RequiredQueryParameters) []byte {
	query := requiredQueryParams.SearchQuery

	var tbestBets, gBestBets, healthTopics, jobOpenings, federalRegisDocs, relatedTerms []string

	strc := QueryData{query, webResults, tbestBets, gBestBets, healthTopics, jobOpenings, federalRegisDocs, relatedTerms}
	data := structToJSON(strc)

	return data
}

func createQueryWebResultsData(strc interface{}) QueryWebResultsData {
	// convert struct to JSON
	data := structToJSON(strc)

	// from JSON convert to new struct
	var searchResultJSON QueryWebResultsData

	err := json.Unmarshal([]byte(data), &searchResultJSON)
	if err != nil {
		zap.L().Fatal("marshalledResults went fatally wrong",
			zap.Error(err),
		)
	}

	return searchResultJSON
}

func structToJSON(strc interface{}) []byte {
	data, err := json.Marshal(strc)
	if err != nil {
		zap.L().Fatal("Something went fatally wrong",
			zap.Error(err),
		)
	}

	return data
}

func JSONToStruct(marshalledResults []byte) QueryData {
	var qd QueryData
	err := json.Unmarshal([]byte(marshalledResults), &qd)
	if err != nil {
		zap.L().Fatal("Something went fatally wrong",
			zap.Error(err),
		)
	}
	return qd
}

////////////////////
