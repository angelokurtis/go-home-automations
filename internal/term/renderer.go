package term

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	yaml "github.com/goccy/go-yaml"
)

type Renderer interface {
	Render(data []byte, lang string) error
	RenderText(text, lang string) error
	RenderYAML(value any) error
}

type MarkdownRenderer struct {
	glamourRenderer *glamour.TermRenderer
}

func NewMarkdownRenderer(renderer *glamour.TermRenderer) *MarkdownRenderer {
	return &MarkdownRenderer{glamourRenderer: renderer}
}

func (m *MarkdownRenderer) Render(data []byte, lang string) error {
	output, err := m.glamourRenderer.Render(fmt.Sprintf("```%s\n%s\n```", lang, string(data)))
	if err != nil {
		return fmt.Errorf("rendering failed: %w", err)
	}

	println(output)

	return nil
}

func (m *MarkdownRenderer) RenderText(text, lang string) error {
	output, err := m.glamourRenderer.Render(fmt.Sprintf("```%s\n%s\n```", lang, text))
	if err != nil {
		return fmt.Errorf("rendering failed: %w", err)
	}

	println(output)

	return nil
}

func (m *MarkdownRenderer) RenderYAML(value any) error {
	data, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("YAML marshalling failed: %w", err)
	}

	return m.Render(data, "yaml")
}
