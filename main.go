package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
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
	srv := flag.Bool("srv", false, "Run as web server with modern UI")
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
	tbinary := flag.Bool("tbinary", false, "Encode input to binary")
	fbinary := flag.Bool("fbinary", false, "Decode input from binary")
	thtml := flag.Bool("thtml", false, "Encode input to HTML entities")
	fhtml := flag.Bool("fhtml", false, "Decode input from HTML entities")
	help := flag.Bool("help", false, "Show help message")
	tcaesar := flag.Int("tcaesar", 0, "Encode input with Caesar cipher (provide shift)")
	fcaesar := flag.Int("fcaesar", 0, "Decode input with Caesar cipher (provide shift)")
	tmd5 := flag.Bool("tmd5", false, "Encode input to MD5 hash")

	flag.Parse()

	if *srv {
		runWebServer()
		return
	}

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

	case *tbinary:
		status = "Text to Binary"
		BannerStart()
		Result(*data, EncodeBinary(*data))

	case *fbinary:
		status = "Binary to Text"
		BannerStart()
		result, err := DecodeBinary(*data)
		if err != nil {
			fmt.Println(Red + "Binary decode error: " + err.Error() + Reset)
			BannerEnd()
			os.Exit(1)
		}
		Result(*data, result)

	case *thtml:
		status = "Text to HTML Entity"
		BannerStart()
		Result(*data, EncodeHTMLEntity(*data))

	case *fhtml:
		status = "HTML Entity to Text"
		BannerStart()
		Result(*data, DecodeHTMLEntity(*data))

	case *tcaesar != 0:
		status = "Text to Caesar Encode"
		BannerStart()
		Result(*data, EncodeCaesar(*data, *tcaesar))

	case *fcaesar != 0:
		status = "Caesar to Text"
		BannerStart()
		Result(*data, DecodeCaesar(*data, *fcaesar))

	case *tmd5:
		status = "Text to MD5"
		BannerStart()
		Result(*data, EncodeMD5(*data))

	default:
		status = "No valid operation selected"
		BannerStart()
		fmt.Println(Red + "Error: " + Reset + "No valid operation specified. To see flags: " +
			Green + "--help" + Reset)
		BannerEnd()
		os.Exit(1)
	}

	BannerEnd()
}

// --- WEB SERVER ---

func runWebServer() {
	http.HandleFunc("/", servePage)
	http.HandleFunc("/convert", handleConvert)
	fmt.Println(Green + "Scarface Web UI started at: " + Cyan + "http://127.0.0.1:8080" + Reset)
	fmt.Println(Green + "Only accessible from this computer." + Reset)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func servePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, pageHTML)
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	// Hem klasik hem multipart formu işle
	r.ParseForm()
	r.ParseMultipartForm(10 << 20) // 10 MB limit

	data := r.FormValue("data")
	algo := r.FormValue("algo")
	shiftStr := r.FormValue("shift")
	shift, _ := strconv.Atoi(shiftStr)
	var result string
	var err error

	fmt.Println("data:", data, "algo:", algo, "shift:", shift)

	switch algo {
	case "Base64 Encode":
		result = EncodeBase64(data)
	case "Base64 Decode":
		result, err = DecodeBase64(data)
	case "Hex Encode":
		result = EncodeHex(data)
	case "Hex Decode":
		result, err = DecodeHex(data)
	case "URL Encode":
		result = EncodeURL(data)
	case "URL Decode":
		result, err = DecodeURL(data)
	case "ASCII Encode":
		result = EncodeASCII(data)
	case "ASCII Decode":
		result, err = DecodeASCII(data)
	case "Binary Encode":
		result = EncodeBinary(data)
	case "Binary Decode":
		result, err = DecodeBinary(data)
	case "HTML Entity Encode":
		result = EncodeHTMLEntity(data)
	case "HTML Entity Decode":
		result = DecodeHTMLEntity(data)
	case "Caesar Encode":
		result = EncodeCaesar(data, shift)
	case "Caesar Decode":
		result = DecodeCaesar(data, shift)
	case "MD5 Encode":
		result = EncodeMD5(data)
	}
	if err != nil {
		result = "Error: " + err.Error()
	}
	fmt.Fprint(w, result)
}

// --- Modern HTML/CSS ---
const pageHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Scarface Encoder/Decoder</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body {
            background: #181c20;
            color: #e0e0e0;
            font-family: 'Segoe UI', Arial, sans-serif;
            margin: 0; padding: 0;
        }
        .container {
            max-width: 500px;
            margin: 40px auto;
            background: #23272b;
            border-radius: 16px;
            box-shadow: 0 8px 32px #000a;
            padding: 32px 24px 24px 24px;
        }
        h1 {
            text-align: center;
            font-size: 2.5rem;
            letter-spacing: 2px;
            margin-bottom: 16px;
            font-family: 'Segoe UI', Arial, sans-serif;
            background: linear-gradient(90deg, #00c6ff, #0072ff, #00c6ff);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        select, input, textarea, button {
            width: 100%;
            margin: 8px 0;
            padding: 10px;
            border-radius: 8px;
            border: none;
            font-size: 1rem;
            background: #181c20;
            color: #e0e0e0;
            box-sizing: border-box;
        }
        button {
            background: linear-gradient(90deg, #00c6ff, #0072ff);
            color: #fff;
            font-weight: bold;
            cursor: pointer;
            transition: background 0.2s;
        }
        button:hover {
            background: linear-gradient(90deg, #0072ff, #00c6ff);
        }
        .result {
            background: #111417;
            border-radius: 8px;
            padding: 12px;
            margin-top: 12px;
            min-height: 40px;
            white-space: pre-wrap;
            word-break: break-all;
            font-family: 'Fira Mono', 'Consolas', monospace;
            font-size: 1.1rem;
        }
        label {
            font-size: 1rem;
            margin-top: 8px;
            display: block;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Scarface</h1>
        <form id="convertForm" autocomplete="off">
            <label for="data">Input Data</label>
            <textarea id="data" name="data" rows="3" required></textarea>

            <label for="algo">Algorithm</label>
            <select id="algo" name="algo" required>
                <option>Base64 Encode</option>
                <option>Base64 Decode</option>
                <option>Hex Encode</option>
                <option>Hex Decode</option>
                <option>URL Encode</option>
                <option>URL Decode</option>
                <option>ASCII Encode</option>
                <option>ASCII Decode</option>
                <option>Binary Encode</option>
                <option>Binary Decode</option>
                <option>HTML Entity Encode</option>
                <option>HTML Entity Decode</option>
                <option>Caesar Encode</option>
                <option>Caesar Decode</option>
                <option>MD5 Encode</option>
            </select>

            <label for="shift">Caesar Shift</label>
            <input type="number" id="shift" name="shift" value="3" min="1" max="25">

            <button type="submit">Convert</button>
        </form>
        <div class="result" id="result"></div>
    </div>
    <script>
        const algo = document.getElementById('algo');
        const shiftInput = document.getElementById('shift');
        function updateShiftInput() {
            if (algo.value.includes('Caesar')) {
                shiftInput.disabled = false;
            } else {
                shiftInput.disabled = true;
            }
        }
        algo.addEventListener('change', updateShiftInput);
        // Sayfa ilk açıldığında da kontrol et
        updateShiftInput();
        document.getElementById('convertForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            const form = e.target;
            const fd = new FormData(form);
            const res = await fetch('/convert', { method: 'POST', body: fd });
            const text = await res.text();
            document.getElementById('result').textContent = text;
        });
    </script>
</body>
</html>
`

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
	var result string
	for _, c := range input {
		result += fmt.Sprintf("%%%02X", c)
	}
	return result
}

func DecodeURL(input string) (string, error) {
	var result string
	input = strings.ReplaceAll(input, "%", " %")
	fields := strings.Fields(input)
	for _, f := range fields {
		if len(f) == 3 && f[0] == '%' {
			var b byte
			_, err := fmt.Sscanf(f, "%%%02X", &b)
			if err != nil {
				return "", fmt.Errorf("invalid URL encoding: %s", f)
			}
			result += string(b)
		} else {
			result += f
		}
	}
	return result, nil
}

// Binary encode/decode
func EncodeBinary(input string) string {
	var result string
	for i, c := range input {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%08b", c)
	}
	return result
}

func DecodeBinary(input string) (string, error) {
	var result string
	for _, s := range strings.Fields(input) {
		var b byte
		_, err := fmt.Sscanf(s, "%08b", &b)
		if err != nil {
			return "", fmt.Errorf("invalid binary code: %s", s)
		}
		result += string(b)
	}
	return result, nil
}

// HTML Entity encode/decode
func EncodeHTMLEntity(input string) string {
	var result string
	for _, c := range input {
		result += fmt.Sprintf("&#%d;", c)
	}
	return result
}

func DecodeHTMLEntity(input string) string {
	return html.UnescapeString(input)
}

// Caesar encode/decode
func EncodeCaesar(input string, shift int) string {
	var result string
	for _, c := range input {
		if c >= 'a' && c <= 'z' {
			result += string('a' + (c-'a'+rune(shift))%26)
		} else if c >= 'A' && c <= 'Z' {
			result += string('A' + (c-'A'+rune(shift))%26)
		} else {
			result += string(c)
		}
	}
	return result
}

func DecodeCaesar(input string, shift int) string {
	return EncodeCaesar(input, 26-shift%26)
}

// MD5 encode
func EncodeMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func BannerStart() {

	fmt.Println(Red + "------------------------------------------------------------  " + Reset)
	fmt.Println(Red + "Scarface # " + Green + "Version 2.0 # " + Cyan + status + " #")
	fmt.Println(Red + "------------------------------------------------------------  " + Reset)

}

func BannerEnd() {
	fmt.Println(Red + "------------------------------------------------------------  " + Reset)
}

func PrintHelp() {
	status = "Help Menu"
	BannerStart()
	fmt.Println(Green + "Scarface Encoder/Decoder Tool" + Reset)
	fmt.Println(Cyan + "Usage:" + Reset)
	fmt.Println("  " + Green + "scarface" + Reset + " [flags]")
	fmt.Println()
	fmt.Println(Cyan + "Flags:" + Reset)
	fmt.Println(Green + "  -help" + Reset + "           Show this help message")
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
	fmt.Println(Green + "  -tbinary" + Reset + "        Encode input to binary")
	fmt.Println(Green + "  -fbinary" + Reset + "        Decode input from binary")
	fmt.Println(Green + "  -thtml" + Reset + "          Encode input to HTML entities")
	fmt.Println(Green + "  -fhtml" + Reset + "          Decode input from HTML entities")
	fmt.Println(Green + "  -shift" + Reset + "          Shift value for Caesar cipher (default 3)")
	fmt.Println(Green + "  -tcaesar [n]" + Reset + "    Encode input with Caesar cipher (shift n)")
	fmt.Println(Green + "  -fcaesar [n]" + Reset + "    Decode input with Caesar cipher (shift n)")
	fmt.Println(Green + "  -tmd5" + Reset + "           Encode input to MD5 hash")
	fmt.Println()
	BannerEnd()
}

func Result(before string, after string) {

	fmt.Println(Cyan + "Input: \n\t" + Reset + before)
	fmt.Println(Cyan + "------------------")
	fmt.Println("Output: \n\t" + Reset + after)
}
