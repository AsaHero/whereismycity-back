package logger

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

type OrderedJSONFormatter struct{}

func (f *OrderedJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Create a struct with the desired order of fields
	serializableEntry := struct {
		Time    string                 `json:"time"`
		Level   string                 `json:"level"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data,omitempty"`
	}{
		Time:    entry.Time.Format(time.RFC3339),
		Level:   entry.Level.String(),
		Message: entry.Message,
		Data:    make(map[string]interface{}),
	}

	// Sort the keys of the entry.Data map and add them in order to serializableEntry.Data
	var keys []string
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Sorting keys alphabetically
	for _, k := range keys {
		serializableEntry.Data[k] = entry.Data[k]
	}

	serializedData, err := json.Marshal(serializableEntry)
	if err != nil {
		return nil, err
	}

	return append(serializedData, '\n'), nil
}
