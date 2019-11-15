package receiver

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
)

type LogrusLog struct {
	Handler skogul.HandlerRef
}

var logrusLogLogger = logrus.New()

var logMetadataFields *[]string

func (lg *LogrusLog) configureLogger() error {
	// Set up internal logger so we don't cause recursive logging in case we log errors
	logrusLogLogger.SetOutput(io.Writer(os.Stdout))
	logrusLogLogger.WithFields(logrus.Fields{
		"category": "receiver",
		"receiver": "logrus",
	})

	logrusLogLogger.Debug("Configuring logger")
	metadataFields := []string{"category", "receiver", "level"}
	a := make([]string, 0)
	logMetadataFields = &a
	for _, field := range metadataFields {
		*logMetadataFields = append(*logMetadataFields, field)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(lg)
	return nil
}

// Parse logrus logs to a skogul.Container
func (lg *LogrusLog) Parse(bytes []byte) (*skogul.Container, error) {
	logrusLogLogger.Debug("Parsing log line")
	var c *skogul.Container
	err := json.Unmarshal(bytes, &c)
	if err != nil {
		logrusLogLogger.Error("Failed to marshal logrus log")
	}
	return c, nil
}

// Write logrus logs as skogul.Containers to a handler
func (lg *LogrusLog) Write(bytes []byte) (int, error) {
	metadata := make(map[string]interface{})

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		logrusLogLogger.Error("Failed to unmarshal logrus log for sending it to log receiver")
		return 0, err
	}

	// Extract metadata fields from data
	for _, field := range *logMetadataFields {
		metadata[field] = data[field]
	}

	timestampData := data["time"].(string)
	timestamp, err := time.Parse(time.RFC3339, timestampData)
	if err != nil {
		// @ToDo: Do we want to error out if this fails?
		logrusLogLogger.Error("Failed to parse timestamp '%s' '%s'\n", "time", timestampData)
	}

	m := skogul.Metric{
		Time:     &timestamp,
		Metadata: metadata,
		Data:     data,
	}
	c := skogul.Container{
		Metrics: []*skogul.Metric{&m},
	}

	logEntry, err := json.Marshal(c)
	if err != nil {
		logrusLogLogger.Error("Failed to marshal logrus log into container")
	}

	lg.Handler.H.Handle(logEntry)
	return len(bytes), nil
}

// Start initializes the logger and sets up required facilities
func (lg *LogrusLog) Start() error {
	logrusLogLogger.Debug("Starting logger")
	lg.configureLogger()
	// for {
	// 	time.Sleep(time.Hour * 1)
	// }
	return nil
}
