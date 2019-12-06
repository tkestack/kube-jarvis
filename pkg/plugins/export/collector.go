package export

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

// DiagnosticResultItem collect one diagnostic and it's results
type DiagnosticResultItem struct {
	Catalogue diagnose.Catalogue
	Type      string
	Name      string
	Results   []diagnose.Result
}

// Collector just collect diagnostic results and evaluation results
type Collector struct {
	Format      string
	Diagnostics []*DiagnosticResultItem
	Output      []io.Writer
}

// CoordinateBegin export information about coordinator Run begin
func (c *Collector) CoordinateBegin(ctx context.Context) error {
	if c.Format == "" {
		c.Format = "json"
	}
	return nil
}

// CoordinateFinish export information about coordinator Run finish
func (c *Collector) CoordinateFinish(ctx context.Context) error {
	data, err := c.Marshal()
	if err != nil {
		return err
	}

	for _, out := range c.Output {
		if _, err := out.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// DiagnosticBegin export information about a Diagnostic begin
func (c *Collector) DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error {
	param := dia.Meta()
	c.Diagnostics = append(c.Diagnostics, &DiagnosticResultItem{
		Catalogue: dia.Meta().Catalogue,
		Type:      param.Type,
		Name:      param.Name,
	})
	return nil
}

// DiagnosticResult export information about one diagnose.Result
func (c *Collector) DiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) error {
	dLen := len(c.Diagnostics)
	c.Diagnostics[dLen-1].Results = append(c.Diagnostics[dLen-1].Results, *result)
	return nil
}

// DiagnosticFinish export information about a Diagnostic finished
func (c *Collector) DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error {
	return nil
}

// Marshal marshal Collected results to byte data according to format
// format can be : "json" , "yaml"
func (c *Collector) Marshal() ([]byte, error) {
	switch c.Format {
	case "json":
		return json.Marshal(c)
	case "yaml":
		return yaml.Marshal(c)
	}

	return nil, fmt.Errorf("unknow format")
}

// Unmarshal unmarshal data to Collector
// format can be : "json" , "yaml"
func (c *Collector) Unmarshal(data []byte) error {
	switch c.Format {
	case "json":
		return json.Unmarshal(data, c)
	case "yaml":
		return yaml.Unmarshal(data, c)
	}
	return fmt.Errorf("unknow format")
}
