package stringx

// CompletionRight Completion content with flag on right
func CompletionRight(content, flag string, length int) string {
	if length <= 0 {
		return ""
	}
	if len(content) >= length {
		return string(content[0:length])
	}

	flagsLegth := length - len(content)
	flags := flag
	for {
		if len(flags) == flagsLegth {
			break
		} else if len(flags) > flagsLegth {
			flags = string(flags[0:flagsLegth])
			break
		} else {
			flags = flags + flag
		}
	}
	return content + flags
}

// CompletionLeft Completion content with flag on left
func CompletionLeft(content, flag string, length int) string {
	if length <= 0 {
		return ""
	}
	if len(content) >= length {
		return string(content[0:length])
	}
	flagsLegth := length - len(content)
	flags := flag
	for {
		if len(flags) == flagsLegth {
			break
		} else if len(flags) > flagsLegth {
			flags = string(flags[0:flagsLegth])
			break
		} else {
			flags = flags + flag
		}
	}
	return flags + content
}
