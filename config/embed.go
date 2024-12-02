package config

import (
	"bytes"
	"embed"
	"time"

	"github.com/google/go-jsonnet"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//go:embed *.yaml
//go:embed schedules/*.json
var ConfigFs embed.FS

func ReadConfigJsonnet(sonnetFilename string) string {
	bytes, _ := ConfigFs.ReadFile(sonnetFilename)

	vm := jsonnet.MakeVM()
	json, err := vm.EvaluateAnonymousSnippet(sonnetFilename, string(bytes))
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return json
}

func ReadJsonConfig(jsonFilename string) string {
	json_bytes, err := ConfigFs.ReadFile(jsonFilename)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return string(json_bytes)
}

func GetYamlFileReader(yamlFilename string) *bytes.Reader {
	yaml_bytes, err := ConfigFs.ReadFile(yamlFilename)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return bytes.NewReader(yaml_bytes)
}

func GetScheduleFromHost(host string, schedule string) string {
	// This cannot come from the Env, because that would be a circular import.
	// So, this is a big FIXME.
	cfg := ReadJsonConfig(schedule)
	hostSections := make(map[string]string, 0)
	for _, section := range gjson.Parse(cfg).Get("@keys").Array() {
		for _, site := range gjson.Get(cfg, section.String()).Array() {
			hostSections[site.Get("host").String()] = section.String()
		}
	}
	return hostSections[host]
}

func GetListOfHosts(schedule string) []string {
	zap.L().Debug("reading in hosts", zap.String("schedule", schedule))
	cfg := ReadJsonConfig(schedule)
	hosts := make([]string, 0)
	for _, section := range gjson.Parse(cfg).Get("@keys").Array() {
		for _, site := range gjson.Get(cfg, section.String()).Array() {
			hosts = append(hosts, site.Get("host").String())
		}
	}
	return hosts
}

func SectionToTimestamp(section string, start_time time.Time) time.Time {
	switch section {
	case "daily":
		return start_time.Add(24 * time.Hour)
	case "weekly":
		return start_time.Add(7 * 24 * time.Hour)
	case "bi-weekly":
		return start_time.Add(14 * 24 * time.Hour)
	case "monthly":
		return start_time.Add(30 * 24 * time.Hour)
	case "quarterly":
		return start_time.Add(3 * 30 * 24 * time.Hour)
	case "bi-annually":
		return start_time.Add(6 * 30 * 24 * time.Hour)
	default:
		// We will default to `montly` to be safe
		return start_time.Add(time.Duration(30*24) * time.Hour)
	}
}

func HostToPgTimestamp(host string, schedule string, start_time time.Time) pgtype.Timestamp {
	sched := GetScheduleFromHost(host, schedule)
	return SectionToPgTimestamp(sched, start_time)
}

func SectionToPgTimestamp(section string, start_time time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:             SectionToTimestamp(section, start_time),
		InfinityModifier: 0,
		Valid:            true,
	}
}
