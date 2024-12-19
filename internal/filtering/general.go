package filter

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// FIXME: Perhaps these should all return error.
// It can either be "nil" or a message.
// That way, the failures contain information
// (as opposed to just being a bool.)
// FIXME: These should take URLs, not strings.

func exceedsLength(length int) func(*url.URL) error {
	return func(u *url.URL) error {
		if len(u.String()) > length {
			return fmt.Errorf("exceeds length [%d]: %s", length, u.String())
		}
		return nil
	}
}

func hasSlashHttp(u *url.URL) error {
	m, _ := regexp.MatchString(`/http`, u.Path)
	if m {
		return fmt.Errorf("http in middle of url: %s", u.Path)
	}
	return nil
}

func insecureGov(u *url.URL) error {
	m, _ := regexp.MatchString(`^http:`, u.String())
	if m {
		return fmt.Errorf("insecure URL: %s", u.String())
	}
	return nil
}

func isTooShort(length int) func(*url.URL) error {
	return func(u *url.URL) error {
		if len(u.String()) < length {
			return fmt.Errorf("too short [%d]: %s", length, u.String())
		}
		return nil
	}
}

// FIXME: move to a config file.
var skippable_prefixes = []string{"#", "mailto"}

func hasSkippablePrefixRelative(u *url.URL) error {
	for _, sp := range skippable_prefixes {
		if strings.HasPrefix(u.String(), sp) {
			return fmt.Errorf("skippable prefix [%s]: %s", sp, u.Path)
		}
	}
	return nil
}

var skippable_extensions = []string{"epub", "stl", "docx", "xlsx", "doc", "txt", "xls", "jpg", "jpeg", "png", "tiff", "tif", "gif", "svg", "raw", "psd", "mp3", "mov", "webp", "bmp", "acc", "ogg"}

func hasSkippableExtension(u *url.URL) error {
	for _, ext := range skippable_extensions {
		if strings.HasSuffix(u.Path, ext) {
			return fmt.Errorf("skippable extension [%s]: %s", ext, u.Path)
		}
	}
	return nil
}

func hasTooManyRepeats(repeatLength int, threshold int) func(*url.URL) error {
	return func(u *url.URL) error {
		s := u.String()
		end := len(s) - repeatLength
		chunks := make(map[string]bool)
		repeats := make(map[string]int)
		for ndx := 0; ndx < end; ndx++ {
			piece := s[ndx : ndx+repeatLength]
			if _, ok := chunks[piece]; ok {
				repeats[piece] = repeats[piece] + 1
			} else {
				chunks[piece] = true
				repeats[piece] = 0
			}
		}

		total := 0
		for _, v := range repeats {
			total += v
		}

		if total >= threshold {
			return fmt.Errorf("too many repeats [%d over %d]: %s", total, threshold, u.String())
		}
		return nil
	}
}

var all string = ".*"

func GeneralRules() []Rule {
	rules := make([]Rule, 0)

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "insecure gov",
		Reject: insecureGov,
	})

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "max isTooShort 5",
		Reject: isTooShort(5),
	})

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "exceedsLength 200",
		Reject: exceedsLength(200),
	})

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "hasSlashHttp",
		Reject: hasSlashHttp,
	})

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "insecureGov",
		Reject: insecureGov,
	})

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "hasSkippablePrefixRelative",
		Reject: hasSkippablePrefixRelative,
	})

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "hasSkippableExtension",
		Reject: hasSkippableExtension,
	})

	rules = append(rules, Rule{
		Match:  all,
		Msg:    "hasTooManyRepeats",
		Reject: hasTooManyRepeats(8, 50),
	})

	return rules
}
