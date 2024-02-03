package config

import (
	"bytes"
	"github.com/Zcentury/gologger/formatter"
	"github.com/Zcentury/gologger/levels"
	"github.com/logrusorgru/aurora/v4"
)

type CLI struct {
	UseColors bool
}

func NewCLI(useColors bool) *CLI {
	return &CLI{
		UseColors: useColors,
	}
}

func (c *CLI) Format(event *formatter.LogEvent) ([]byte, error) {
	c.colorizeLabel(event)

	buffer := &bytes.Buffer{}
	buffer.Grow(len(event.Message))

	if label, ok := event.Metadata["label"]; label != "" && ok {
		buffer.WriteString(label)
		buffer.WriteRune('|')
		delete(event.Metadata, "label")
	}

	if timestamp, ok := event.Metadata["timestamp"]; timestamp != "" && ok {
		buffer.WriteString(aurora.Bold(aurora.Green(timestamp)).String())
		buffer.WriteString("| ")
		delete(event.Metadata, "timestamp")
	}
	buffer.WriteString(event.Message)
	if len(event.Metadata) > 0 {
		buffer.WriteString(" (")
	}
	for k, v := range event.Metadata {
		buffer.WriteString(aurora.Bold(k).String())
		buffer.WriteRune('=')
		buffer.WriteString(v)
		buffer.WriteRune(',')
	}

	data := buffer.Bytes()

	if len(event.Metadata) > 0 {
		data = data[:len(data)-1]
		data = append(data, ')')
	}

	return data, nil
}

func (c *CLI) colorizeLabel(event *formatter.LogEvent) {
	label := event.Metadata["label"]
	if label == "" || !c.UseColors {
		return
	}
	switch event.Level {
	case levels.LevelInfo:
		event.Metadata["label"] = aurora.Bold(aurora.BgGreen(" " + label + " ")).String()
	case levels.LevelFatal:
		event.Metadata["label"] = aurora.Bold(aurora.BgMagenta(" " + label + " ")).String()
	case levels.LevelError:
		event.Metadata["label"] = aurora.Bold(aurora.BgRed(" " + label + " ")).String()
	case levels.LevelDebug:
		event.Metadata["label"] = aurora.Bold(aurora.BgBlue(" " + label + " ")).String()
	case levels.LevelWarning:
		event.Metadata["label"] = aurora.Bold(aurora.BgYellow(" " + label + " ")).String()
	}
}
