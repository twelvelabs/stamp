package modify

func String(subject string, action Action, arg string, _ ModifierConf) string {
	var modified string
	switch action {
	case ActionPrepend:
		modified = arg + subject
	case ActionAppend:
		modified = subject + arg
	case ActionReplace:
		modified = arg
	case ActionDelete:
	}
	return modified
}
