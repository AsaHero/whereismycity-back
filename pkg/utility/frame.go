package utility

import (
	"fmt"
	"runtime"
	"strings"
)

func GetFrameData(skip int) (scope, caller, location string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown", "unknown"
	}

	// Get the full function name
	funcName := runtime.FuncForPC(pc).Name()

	// Split the full path
	parts := strings.Split(funcName, "/")

	// Get the last part which contains the type and method
	lastPart := parts[len(parts)-1]

	// Split the last part by dots to separate package, type, and method
	components := strings.Split(lastPart, ".")

	switch len(components) {
	case 1:
		// Just a function
		scope = "package"
		caller = components[0]
	case 2:
		// Package and function
		scope = components[0]
		caller = components[1]

	case 3:
		if components[1][0] == '(' {
			// Method on a type: e.g., "app.(*App).Run"
			scope = components[0] + "." + strings.Trim(components[1], "()*")
			caller = components[2]
		} else {
			// Package with subpackage and function
			scope = components[0] + "." + components[1]
			caller = components[2]
		}
	default:
		if components[len(components)-2][0] == '(' {
			// Method on a type with more complex path
			scope = components[len(components)-3] + "." + strings.Trim(components[len(components)-2], "()*")
			caller = components[len(components)-1]
		} else {
			// Complex package path with function
			scope = strings.Join(components[:len(components)-1], ".")
			caller = components[len(components)-1]
		}
	}

	location = fmt.Sprintf("%s:%d", file, line)
	return
}
