package config

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFQDN2Domain64(t *testing.T) {
	// Test a lookup against the decimal value
	v, err := FQDNToDomain64("my.goarmy.com")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(216172790703718656), v)

	// Test the same lookup against a hex parse
	v, err = FQDNToDomain64("my.goarmy.com")
	if err != nil {
		t.Error(err)
	}
	i, err := strconv.ParseInt("0300000200000100", 16, 64)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(i), v)

	// And, something in a different TLD
	v, err = FQDNToDomain64("product-guide.18f.gov")
	if err != nil {
		t.Error(err)
	}
	i, err = strconv.ParseInt("0100000100000600", 16, 64)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(i), v)

}
