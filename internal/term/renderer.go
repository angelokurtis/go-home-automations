package term

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/glamour"
	yaml "github.com/goccy/go-yaml"
)

type Renderer interface {
	RenderBytes(input []byte, language string) error
	RenderString(input, language string) error
	RenderYAML(input any) error
	RenderJSON(input []byte) error
}

type ColorRenderer struct {
	renderer *glamour.TermRenderer
}

func NewColorRenderer(renderer *glamour.TermRenderer) *ColorRenderer {
	return &ColorRenderer{renderer: renderer}
}

func (c *ColorRenderer) RenderBytes(input []byte, language string) error {
	output, err := c.renderer.Render(fmt.Sprintf("```%s\n%s\n```", language, string(input)))
	if err != nil {
		return fmt.Errorf("unable to render: %w", err)
	}

	println(output)

	return nil
}

func (c *ColorRenderer) RenderString(input, language string) error {
	output, err := c.renderer.Render(fmt.Sprintf("```%s\n%s\n```", language, input))
	if err != nil {
		return fmt.Errorf("unable to render: %w", err)
	}

	println(output)

	return nil
}

func (c *ColorRenderer) RenderYAML(input any) error {
	bytes, err := yaml.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return c.RenderBytes(bytes, "yaml")
}

func (c *ColorRenderer) RenderJSON(input []byte) error {
	var obj interface{}
	if err := json.Unmarshal(input, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	bytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return c.RenderBytes(bytes, "yaml")
}
