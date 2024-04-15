package rich

// (c) 2024 by github.com/rvfet
// This package is a simple utility to print the formatted text with colors and styles in the terminal.
// Inspired by the rich library in Python.

import (
	"fmt"     // * Required for printing the result
	"reflect" // * Required for detecting the type of the input
	"regexp"  // * Required for parsing the tags
	"strings" // * Required for string manipulation
)

// * Struct to define the color and style properties
type color struct {
	Name    string
	Code    string
	IsColor bool
}

// * List of colors and styles that can be applied to the text
var Colors = []color{
	{
		Name:    "reset",
		Code:    "0",
		IsColor: false,
	},
	{
		Name:    "unstyle",
		Code:    "22",
		IsColor: false,
	},
	{
		Name:    "b",
		Code:    "1",
		IsColor: false,
	},
	{
		Name:    "i",
		Code:    "3",
		IsColor: false,
	},
	{
		Name:    "u",
		Code:    "4",
		IsColor: false,
	},
	{
		Name:    "s",
		Code:    "9",
		IsColor: false,
	},
	{
		Name:    "blink",
		Code:    "5",
		IsColor: false,
	},
	{
		Name:    "x",
		Code:    "7",
		IsColor: false,
	},
	{
		Name:    "white",
		Code:    "97",
		IsColor: true,
	},
	{
		Name:    "gray",
		Code:    "37",
		IsColor: true,
	},
	{
		Name:    "red",
		Code:    "31",
		IsColor: true,
	},
	{
		Name:    "green",
		Code:    "32",
		IsColor: true,
	},
	{
		Name:    "cyan",
		Code:    "36",
		IsColor: true,
	},
	{
		Name:    "blue",
		Code:    "34",
		IsColor: true,
	},
	{
		Name:    "yellow",
		Code:    "33",
		IsColor: true,
	},
}

// & Helper function to get the color code by the color name
func getColorByName(name string) string {
	for _, c := range Colors {
		if c.Name == name {
			return "\033[" + c.Code + "m"
		}
	}
	return "\033[37m"
}

// & Helper function to parse the tags and forward them to `applyStyling`
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

		tagParts := strings.Fields(tags)
		for _, tag := range tagParts {
			tag = strings.ToLower(strings.Trim(tag, "[]"))
			if strings.HasPrefix(tag, "/") {
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			} else {
				found := false
				for _, c := range Colors {
					if c.Name == tag {
						stack = append(stack, c.Code)
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}

		segments[i] = applyStyling(rest, stack)
	}

	styledStr := strings.Join(segments, "")
	return styledStr
}

// & Helper function to convert the curated styles to the ANSI escape codes
func applyStyling(str string, stack []string) string {
	styleCode := strings.Join(stack, ";")
	return fmt.Sprintf("\033[%sm%s", styleCode, str)
}

// & Helper function to colorize some common keywords
func colorizeKeywords(input string) string {
	keywords := map[string]string{
		"success": getColorByName("green"),
		"error":   getColorByName("red"),
		"warning": getColorByName("yellow"),
		"info":    getColorByName("cyan"),
	}

	for keyword, colorCode := range keywords {
		re := regexp.MustCompile(`(?i)(\b` + keyword + `\b)`)
		input = re.ReplaceAllStringFunc(input, func(match string) string {
			return colorCode + match + "\033[0m"
		})
	}

	return input
}

// & Helper function formats the input value based on its type.
func formatValue(v reflect.Value) string {
	var result strings.Builder
	switch v.Kind() {
	case reflect.String:
		result.WriteString(formatString(v.String()))
	case reflect.Map:
		result.WriteString(formatMap(v))
	case reflect.Slice:
		result.WriteString(formatSlice(v))
	case reflect.Struct:
		result.WriteString(formatStruct(v))
	default:
		result.WriteString(parseTags(fmt.Sprintf("%v", v)))
	}
	return result.String()
}

// & Helper function formats a string value.
func formatString(str string) string {
	str = parseTags(str)
	str = colorizeKeywords(str)
	return str
}

// & Helper function formats a boolean value.
func formatBool(v reflect.Value) string {
	if v.Bool() {
		return getColorByName("green") + "true" + getColorByName("white")
	}
	return getColorByName("red") + "false" + getColorByName("white")
}

// & Helper function formats a float or integer or hexadecimal value.
func formatNumber(v reflect.Value) string {
	return getColorByName("cyan") + parseTags(fmt.Sprintf("%v", v)) + getColorByName("white")
}

// & Helper function formats a map value.
func formatMap(v reflect.Value) string {
	var result strings.Builder
	result.WriteString("{\n")
	for _, key := range v.MapKeys() {
		result.WriteString(" " + parseTags(fmt.Sprintf("[gray]%s[/gray]: ", key.String())))
		result.WriteString(formatValue(v.MapIndex(key)) + "\n")
	}
	result.WriteString("}")
	return result.String()
}

// & Helper function formats a slice value.
func formatSlice(v reflect.Value) string {
	var result strings.Builder
	result.WriteString("[ ")
	for i := 0; i < v.Len(); i++ {
		element := v.Index(i)
		switch value := element.Interface().(type) {
		case int, int8, int16, int32, int64:
			result.WriteString(getColorByName("cyan") + fmt.Sprint(value) + getColorByName("white"))
		case string:
			result.WriteString(getColorByName("gray") + "\"" + value + "\"" + getColorByName("white"))
		default:
			result.WriteString(getColorByName("gray") + "\"" + fmt.Sprint(value) + "\"" + getColorByName("white"))
		}
		if i < v.Len()-1 {
			result.WriteString(", ")
		}
	}
	result.WriteString(" ]")
	return result.String()
}

// & Helper function formats a struct value.
func formatStruct(v reflect.Value) string {
	var result strings.Builder
	result.WriteString("{\n")
	for i := 0; i < v.NumField(); i++ {
		result.WriteString(" " + parseTags(fmt.Sprintf("%s: ", v.Type().Field(i).Name)))
		result.WriteString(formatValue(v.Field(i)) + "\n")
	}
	result.WriteString("}")
	return result.String()
}

// * Main function that's been exposed to other packages to print the formatted value
func Print(args ...interface{}) {
	var formattedStrings []string

	for _, arg := range args {
		// Determine the type of the argument
		argType := reflect.TypeOf(arg)

		// Format the argument based on its type
		var formattedArg string
		switch argType.Kind() {
		case reflect.String:
			formattedArg = formatString(arg.(string))
		case reflect.Bool:
			formattedArg = formatBool(reflect.ValueOf(arg))
		case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			formattedArg = formatNumber(reflect.ValueOf(arg))
		case reflect.Map:
			formattedArg = formatMap(reflect.ValueOf(arg))
		case reflect.Slice:
			formattedArg = formatSlice(reflect.ValueOf(arg))
		case reflect.Struct:
			formattedArg = formatStruct(reflect.ValueOf(arg))
		default:
			formattedArg = parseTags(fmt.Sprintf("%v", arg))
		}

		// Append the formatted argument to the slice
		formattedStrings = append(formattedStrings, formattedArg)
	}

	// Concatenate all formatted strings into a single string
	formattedOutput := strings.Join(formattedStrings, " ")

	// Print the concatenated string
	fmt.Println(formattedOutput)
}

// * Info - Print the input with an abbreviation and blue color using Print
func Info(input interface{}) {
	Print("[blue][b]INFO:[/b][/blue] " + input.(string))
}

// * Success - Print the input with an abbreviation and green color using Print
func Success(input interface{}) {
	Print("[green][b]SUCC:[/b][/green] " + input.(string))
}

// * Error - Print the input with an abbreviation and red color using Print
func Error(input interface{}) {
	Print("[red][b]ERRR:[/b][/red] " + input.(string))
}

// * Warning - Print the input with an abbreviation and yellow color using Print
func Warning(input interface{}) {
	Print("[yellow][b]WARN:[/b][/yellow] " + input.(string))
}

// * Debug - Print the input with an abbreviation and gray color using Print
func Debug(input interface{}) {
	Print("[gray][b]DEBUG:[/b][/gray] " + input.(string))
}
