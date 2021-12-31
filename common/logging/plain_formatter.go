package logging

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type PlainFormatterWithTsWithCaller struct {
}
type PlainFormatterWithoutTsWithCaller struct {
}

type PlainFormatterWithTsWithoutCaller struct {
}
type PlainFormatterWithoutTsWithoutCaller struct {
}

func (f *PlainFormatterWithTsWithCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] [%s] %s%s\n", entry.Time.Format("2006-1-2 15:04:05"), entry.Level, entry.Caller.File, entry.Message, formatFields(&entry.Data))), nil
}
func (f *PlainFormatterWithoutTsWithCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] %s%s\n", entry.Level, entry.Caller.File, entry.Message, formatFields(&entry.Data))), nil
}

func (f *PlainFormatterWithTsWithoutCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] %s%s\n", entry.Time.Format("2006-1-2 15:04:05"), entry.Level, entry.Message, formatFields(&entry.Data))), nil
}
func (f *PlainFormatterWithoutTsWithoutCaller) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] %s%s\n", entry.Level, entry.Message, formatFields(&entry.Data))), nil
}

func formatFields(fields *log.Fields) string {
	if len(*fields) < 1 {
		return ""
	}
	entries := make([]string, 0, len(*fields))
	for k, v := range *fields {
		entries = append(entries, fmt.Sprintf("%s: %v", k, v))
	}
	return ", " + strings.Join(entries, ", ")
}
