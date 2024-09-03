package rich

// (c) 2024 by github.com/rvfet
// This package is a simple utility to print the formatted text with colors and styles in the terminal.
// Inspired by the rich library in Python.

import (
	"fmt"
	"reflect"
	"strings"
)

type Style struct {
	Name    string
	Code    string
	IsColor bool
}

var KeywordMap = map[string]string{
	"SUCCESS": "green",
	"ERROR":   "red",
	"WARNING": "yellow",
	"INFO":    "cyan",
	"DEBUG":   "gray",
}

var (
	styles = []Style{
		{Name: "reset", Code: "0", IsColor: false},
		{Name: "unstyle", Code: "22", IsColor: false},
		{Name: "b", Code: "1", IsColor: false},
		{Name: "i", Code: "3", IsColor: false},
		{Name: "u", Code: "4", IsColor: false},
		{Name: "s", Code: "9", IsColor: false},
		{Name: "blink", Code: "5", IsColor: false},
		{Name: "x", Code: "7", IsColor: false},
		{Name: "white", Code: "97", IsColor: true},
		{Name: "gray", Code: "37", IsColor: true},
		{Name: "red", Code: "31", IsColor: true},
		{Name: "green", Code: "32", IsColor: true},
		{Name: "cyan", Code: "36", IsColor: true},
		{Name: "blue", Code: "34", IsColor: true},
		{Name: "yellow", Code: "33", IsColor: true},
	}

	styleMap = make(map[string]Style)
)

var formatterMap map[reflect.Kind]func(reflect.Value) string

func init() {
	for _, style := range styles {
		styleMap[style.Name] = style
	}

	formatterMap = map[reflect.Kind]func(reflect.Value) string{
		reflect.String:  formatString,
		reflect.Bool:    formatBool,
		reflect.Float32: formatNumber,
		reflect.Float64: formatNumber,
		reflect.Int:     formatNumber,
		reflect.Int8:    formatNumber,
		reflect.Int16:   formatNumber,
		reflect.Int32:   formatNumber,
		reflect.Int64:   formatNumber,
		reflect.Uint:    formatNumber,
		reflect.Uint8:   formatNumber,
		reflect.Uint16:  formatNumber,
		reflect.Uint32:  formatNumber,
		reflect.Uint64:  formatNumber,
		reflect.Map:     formatMap,
		reflect.Slice:   formatSlice,
		reflect.Struct:  formatStruct,
	}
}

func parseTags(str string) string {
	var stack []string
	segments := strings.Split(str, "[")

	for index, segment := range segments {
		if index == 0 {
			continue
		}
		parts := strings.SplitN(segment, "]", 2)
		if len(parts) != 2 {
			continue
		}
		tags, rest := parts[0], parts[1]

		for _, tag := range strings.Fields(tags) {
			tag = strings.ToLower(strings.Trim(tag, "[]"))
			if tag == "/" {
				stack = nil
			} else if strings.HasPrefix(tag, "/") {
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			} else if style, ok := styleMap[tag]; ok {
				stack = append(stack, style.Code)
			}
		}

		segments[index] = applyStyling(rest, stack)
	}

	return strings.Join(segments, "")
}

func applyStyling(str string, stack []string) string {
	return fmt.Sprintf("\033[%sm%s", strings.Join(stack, ";"), str)
}

func formatValue(value reflect.Value) string {
	if formatter, ok := formatterMap[value.Kind()]; ok {
		return formatter(value)
	}

	return formatString(value)
}

func formatString(str reflect.Value) string {
	return parseTags(fmt.Sprintf("%v", str))
}

func formatBool(value reflect.Value) string {
	if reflect.ValueOf(value.Interface()).Bool() {
		return parseTags("[green][bold]true[/]")
	}

	return parseTags("[red][bold]false[/]")
}

func formatNumber(value reflect.Value) string {
	return parseTags(fmt.Sprintf("[cyan][bold]%v[/]", value))
}

func formatMap(value reflect.Value) string {
	var result strings.Builder
	result.WriteString("{\n")
	for _, key := range value.MapKeys() {
		leftSide := formatValue(key)
		rightSide := formatValue(value.MapIndex(key))

		result.WriteString(fmt.Sprintf("  \"%s\": %s,\n", leftSide, rightSide))
	}
	result.WriteString("}")

	return result.String()
}

func formatSlice(value reflect.Value) string {
	var result strings.Builder
	result.WriteString("[ ")
	for index := 0; index < value.Len(); index++ { //nolint:intrange
		element := value.Index(index)
		result.WriteString(formatValue(element))
		if index < value.Len()-1 {
			result.WriteString(", ")
		}
	}
	result.WriteString(" ]")

	return result.String()
}

func formatStruct(value reflect.Value) string {
	var result strings.Builder
	result.WriteString("{\n")
	for index := 0; index < value.NumField(); index++ { //nolint:intrange
		leftSide := parseTags(fmt.Sprintf("[yellow]%s[/]", value.Type().Field(index).Name))
		rightSide := formatValue(value.Field(index))
		result.WriteString(fmt.Sprintf("  %s: %s,\n", leftSide, rightSide))
	}
	result.WriteString("}")

	return result.String()
}

func Print(args ...any) {
	formattedStrings := make([]string, 0, len(args))

	for _, arg := range args {
		formattedStrings = append(formattedStrings, formatValue(reflect.ValueOf(arg)))
	}

	fmt.Println(strings.Join(formattedStrings, " "))
}

func logWithPrefix(prefix string, args ...any) {
	var pad string
	if len(pad) < 10 {
		pad += strings.Repeat(" ", 8-len(prefix))
	}

	if color, ok := KeywordMap[prefix]; ok {
		prefix = fmt.Sprintf("[%s][b]%s[/]", color, prefix)
	}
	prefix += pad

	Print(append([]any{prefix}, args...)...)
}

func Info(args ...any) {
	logWithPrefix("INFO", args...)
}

func Success(args ...any) {
	logWithPrefix("SUCCESS", args...)
}

func Error(args ...any) {
	logWithPrefix("ERROR", args...)
}

func Warning(args ...any) {
	logWithPrefix("WARNING", args...)
}

func Debug(args ...any) {
	logWithPrefix("DEBUG", args...)
}
