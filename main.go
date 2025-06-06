package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"net/url"
	"os"
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
	toBase64 := flag.Bool("to_base64", false, "Conver to Base64")
	fromBase64 := flag.Bool("from_base64", false, "Convert from Base64")
	toHex := flag.Bool("to_hex", false, "Convert to Hex")
	fromHex := flag.Bool("from_hex", false, "Convert from Hex")
	toURL := flag.Bool("to_url", false, "Convert to Url encode")
	fromURL := flag.Bool("from_url", false, "Convert from Url encode")

	flag.Parse()
	if *data == "" {
		status = "No data provided"
		BannerStart()
		fmt.Println(Red + "Error: " + Reset + "No data provided. Use" + Green + " -data" + Reset + " flag to specify the input data.")
		BannerEnd()
		os.Exit(1)
	}
	switch {
	case *toBase64:
		status = "Text to Base64"
		BannerStart()
		Result(*data, EncodeBase64(*data))

	case *fromBase64:

		result, err := DecodeBase64(*data)
		if err != nil {
			fmt.Println(Red + "Base64 decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		fmt.Println(Green + result + Reset)

	case *toHex:

		status = "Text to Hex"
		BannerStart()
		Result(*data, EncodeHex(*data))

	case *fromHex:
		result, err := DecodeHex(*data)
		if err != nil {
			fmt.Println(Red + "Hex decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		Result(*data, result)

	case *toURL:
		status = "Text to URL Encode"
		BannerStart()
		Result(*data, EncodeBase64(*data))

	case *fromURL:
		result, err := DecodeURL(*data)
		if err != nil {
			fmt.Println(Red + "URL decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		Result(*data, result)
	default:
		status = "No valid operation selected"
		BannerStart()
		fmt.Println(Red + "Error: " + Reset + "No valid operation specified. Use one of the flags: " +
			Green + "-to_base64, -from_base64, -to_hex, -from_hex, -to_url, -from_url" + Reset)
		BannerEnd()
		os.Exit(1)
	}

	BannerEnd()
}

func Result(before string, after string) {

	fmt.Println(Cyan + "Input: " + Reset + before)
	fmt.Println(Cyan + "------------------")
	fmt.Println("Output: " + Reset + after)
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
	fmt.Println(Red + "Scarface # " + Green + "Version 1.0 # " + Cyan + status + " #")
	fmt.Println(Red + "------------------------------------------------------------  " + Reset)

}

func BannerEnd() {
	fmt.Println(Red + "------------------------------------------------------------  " + Reset)
}
