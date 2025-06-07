package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Cyan  = "\033[36m"
	Reset = "\033[0m"
)

var status string

func main() {
	// Command Line Flags
	data := flag.String("data", "", "Data to encode/decode")
	tb64 := flag.Bool("tb64", false, "Convert to Base64")
	fb64 := flag.Bool("fb64", false, "Convert from Base64")
	th := flag.Bool("th", false, "Convert to Hex")
	fh := flag.Bool("fh", false, "Convert from Hex")
	tu := flag.Bool("tu", false, "Convert to URL encode")
	fu := flag.Bool("fu", false, "Convert from URL encode")
	h := flag.Bool("h", false, "Show help message")
	tascii := flag.Bool("tasc", false, "Convert to ASCII")
	fascii := flag.Bool("fasc", false, "Convert from ASCII")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *h || *help {
		PrintHelp()
		return
	}

	if *data == "" {
		status = "No data provided"
		BannerStart()
		fmt.Println(Red + "Error: " + Reset + "No data provided. Use" + Green + " -data" + Reset + " flag to specify the input data.")
		BannerEnd()
		os.Exit(1)
	}

	switch {
	case *tb64:
		status = "Text to Base64"
		BannerStart()
		Result(*data, EncodeBase64(*data))

	case *fb64:
		result, err := DecodeBase64(*data)
		if err != nil {
			fmt.Println(Red + "Base64 decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		Result(*data, result)

	case *th:
		status = "Text to Hex"
		BannerStart()
		Result(*data, EncodeHex(*data))

	case *fh:
		result, err := DecodeHex(*data)
		if err != nil {
			fmt.Println(Red + "Hex decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		Result(*data, result)

	case *tu:
		status = "Text to URL Encode"
		BannerStart()
		Result(*data, EncodeURL(*data))

	case *fu:
		result, err := DecodeURL(*data)
		if err != nil {
			fmt.Println(Red + "URL decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		Result(*data, result)

	case *tascii:
		status = "Text to ASCII"
		BannerStart()
		Result(*data, EncodeASCII(*data))

	case *fascii:
		status = "ASCII to Text"
		BannerStart()
		result, err := DecodeASCII(*data)
		if err != nil {
			fmt.Println(Red + "ASCII decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		Result(*data, result)
	default:
		status = "No valid operation selected"
		BannerStart()
		fmt.Println(Red + "Error: " + Reset + "No valid operation specified. Use one of the flags: " +
			Green + "-tb64, -fb64, -th, -fh, -tu, -fu" + Reset)
		BannerEnd()
		os.Exit(1)
	}

	BannerEnd()
}

func PrintHelp() {
	status = "Help Menu"
	BannerStart()
	fmt.Println(Green + "Scarface Encoder/Decoder Tool" + Reset)
	fmt.Println(Cyan + "Usage:" + Reset)
	fmt.Println("  " + Green + "scarface" + Reset + " [flags]")
	fmt.Println()
	fmt.Println(Cyan + "Flags:" + Reset)
	fmt.Println(Green + "  -data" + Reset + "           Input data to encode/decode")
	fmt.Println(Green + "  -tb64" + Reset + "           Encode input to Base64")
	fmt.Println(Green + "  -fb64" + Reset + "           Decode input from Base64")
	fmt.Println(Green + "  -th" + Reset + "             Encode input to Hex")
	fmt.Println(Green + "  -fh" + Reset + "             Decode input from Hex")
	fmt.Println(Green + "  -tu" + Reset + "             Encode input to URL")
	fmt.Println(Green + "  -fu" + Reset + "             Decode input from URL")
	fmt.Println(Green + "  -h" + Reset + "              Show this help message")
	fmt.Println(Green + "  -tasc" + Reset + "           Encode input to ASCII")
	fmt.Println(Green + "  -fasc" + Reset + "           Decode input from ASCII")
	fmt.Println(Green + "  -help" + Reset + "           Show this help message")
	fmt.Println()
	BannerEnd()
}

func Result(before string, after string) {

	fmt.Println(Cyan + "Input: \n\t" + Reset + before)
	fmt.Println(Cyan + "------------------")
	fmt.Println("Output: \n\t" + Reset + after)
}

// ASCII encode/decode
func EncodeASCII(input string) string {
	var result string
	for i, c := range input {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d", c)
	}
	return result
}

func DecodeASCII(input string) (string, error) {
	var result string
	var code int
	for _, s := range strings.Fields(input) {
		_, err := fmt.Sscanf(s, "%d", &code)
		if err != nil {
			return "", fmt.Errorf("invalid ASCII code: %s", s)
		}
		result += string(rune(code))
	}
	return result, nil
}

// Base64 encode/decode
func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func DecodeBase64(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	return string(data), err
}

// Hex encode/decode
func EncodeHex(input string) string {
	return hex.EncodeToString([]byte(input))
}

func DecodeHex(input string) (string, error) {
	data, err := hex.DecodeString(input)
	return string(data), err
}

// URL encode/decode
func EncodeURL(input string) string {
	return url.QueryEscape(input)
}

func DecodeURL(input string) (string, error) {
	result, err := url.QueryUnescape(input)
	return result, err
}

func BannerStart() {

	fmt.Println(Red + "------------------------------------------------------------  " + Reset)
	fmt.Println(Red + "Scarface # " + Green + "Version 1.2 # " + Cyan + status + " #")
	fmt.Println(Red + "------------------------------------------------------------  " + Reset)

}

func BannerEnd() {
	fmt.Println(Red + "------------------------------------------------------------  " + Reset)
}
