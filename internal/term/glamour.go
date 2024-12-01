package term

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func NewGlamourTermRenderer() (*glamour.TermRenderer, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create renderer: %w", err)
	}

	return renderer, nil
}
