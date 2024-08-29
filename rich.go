package rich

// (c) 2024 by github.com/rvfet
// This package is a simple utility to print the formatted text with colors and styles in the terminal.
// Inspired by the rich library in Python.

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Style struct {
	Name    string
	Code    string
	IsColor bool
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

func init() {
	for _, style := range styles {
		styleMap[style.Name] = style
	}
}

func getStyle(name string) string {
	if style, ok := styleMap[name]; ok {
		return "\033[" + style.Code + "m"
	}
	return "\033[37m"
}

func parseTags(str string) string {
	var stack []string
	segments := strings.Split(str, "[")

	for i, segment := range segments {
		if i == 0 {
			continue
		}
		parts := strings.SplitN(segment, "]", 2)
		if len(parts) != 2 {
			continue
		}
		tags, rest := parts[0], parts[1]

		for _, tag := range strings.Fields(tags) {
			tag = strings.ToLower(strings.Trim(tag, "[]"))
			if strings.HasPrefix(tag, "/") {
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			} else if style, ok := styleMap[tag]; ok {
				stack = append(stack, style.Code)
			}
		}

		segments[i] = applyStyling(rest, stack)
	}

	return strings.Join(segments, "")
}

func applyStyling(str string, stack []string) string {
	return fmt.Sprintf("\033[%sm%s", strings.Join(stack, ";"), str)
}

func colorizeKeywords(input string) string {
	keywords := map[string]string{
		"success": getStyle("green"),
		"error":   getStyle("red"),
		"warning": getStyle("yellow"),
		"info":    getStyle("cyan"),
	}

	for keyword, colorCode := range keywords {
		re := regexp.MustCompile(`(?i)(\b` + keyword + `\b)`)
		input = re.ReplaceAllStringFunc(input, func(match string) string {
			return colorCode + match + "\033[0m"
		})
	}

	return input
}

func formatValue(v reflect.Value) string {
	switch reflect.TypeOf(v.Interface()).Kind() {
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return formatNumber(v)
	case reflect.Map:
		return formatMap(v)
	case reflect.Slice:
		return formatSlice(v)
	case reflect.Struct:
		return formatStruct(v)
	case reflect.Bool:
		return formatBool(v)
	default:
		return formatString(v.Interface())
	}
}

func formatString(str interface{}) string {
	return colorizeKeywords(parseTags(fmt.Sprintf("%v", str)))
}

func formatBool(v reflect.Value) string {
	if reflect.ValueOf(v.Interface()).Bool() {
		return parseTags("[green][bold]true[/]")
	}
	return parseTags("[red][bold]false[/]")
}

func formatNumber(v interface{}) string {
	return parseTags(fmt.Sprintf("[cyan][bold]%v[/]", v))
}

func formatMap(v reflect.Value) string {
	var result strings.Builder
	result.WriteString("{\n")
	for _, key := range v.MapKeys() {
		leftSide := formatValue(key)
		rightSide := formatValue(v.MapIndex(key))

		result.WriteString(fmt.Sprintf("  \"%s\": %s,\n", leftSide, rightSide))
	}
	result.WriteString("}")
	return result.String()
}

func formatSlice(v reflect.Value) string {
	var result strings.Builder
	result.WriteString("[ ")
	for i := 0; i < v.Len(); i++ {
		element := v.Index(i)
		result.WriteString(formatValue(element))
		if i < v.Len()-1 {
			result.WriteString(", ")
		}
	}
	result.WriteString(" ]")
	return result.String()
}

func formatStruct(v reflect.Value) string {
	var result strings.Builder
	result.WriteString("{\n")
	for i := 0; i < v.NumField(); i++ {
		leftSide := parseTags(fmt.Sprintf("[yellow]%s[/]", v.Type().Field(i).Name))
		rightSide := formatValue(v.Field(i))
		result.WriteString(fmt.Sprintf("  %s: %s,\n", leftSide, rightSide))
	}
	result.WriteString("}")
	return result.String()
}

func Print(args ...interface{}) {
	var formattedStrings []string

	for _, arg := range args {
		v := reflect.ValueOf(arg)
		var formattedArg string

		switch v.Kind() {
		case reflect.String:
			formattedArg = formatString(v.String())
		case reflect.Bool:
			formattedArg = formatBool(v)
		case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			formattedArg = formatNumber(v)
		case reflect.Map:
			formattedArg = formatMap(v)
		case reflect.Slice:
			formattedArg = formatSlice(v)
		case reflect.Struct:
			formattedArg = formatStruct(v)
		default:
			formattedArg = parseTags(fmt.Sprintf("%v", arg))
		}

		formattedStrings = append(formattedStrings, formattedArg)
	}

	fmt.Println(strings.Join(formattedStrings, " "))
}

func logWithPrefix(prefix string, args ...interface{}) {
	Print(append([]interface{}{prefix}, args...)...)
}

func Info(args ...interface{}) {
	logWithPrefix("[blue][b]INFO:[/b][/blue]", args...)
}

func Success(args ...interface{}) {
	logWithPrefix("[green][b]SUCC:[/b][/green]", args...)
}

func Error(args ...interface{}) {
	logWithPrefix("[red][b]ERRR:[/b][/red]", args...)
}

func Warning(args ...interface{}) {
	logWithPrefix("[yellow][b]WARN:[/b][/yellow]", args...)
}

func Debug(args ...interface{}) {
	logWithPrefix("[gray][b]DEBUG:[/b][/gray]", args...)
}
