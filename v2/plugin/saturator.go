package plugin

// Interface responsible for dismissal of metric modifiers
type ModifierCloser interface {
	Dismiss()
}
