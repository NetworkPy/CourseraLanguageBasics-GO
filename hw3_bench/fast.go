package main

import (
	"bufio"
	"bytes"
	json "encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

//easyjson:json
type Users struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"-"`
	Country  string   `json:"-"`
	Email    string   `json:"email"`
	Job      string   `json:"-"`
	Name     string   `json:"name"`
	Phone    string   `json:"-"`
}

// var browserRegExp = regexp.MustCompile("MSIE")
// var androidRegExp = regexp.MustCompile("Android")

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	browsers := make(map[string]bool, 120)

	scanner := bufio.NewScanner(file)

	i := -1 // index
	user := Users{}
	fmt.Fprintln(out, "found users:")

	for scanner.Scan() {
		i++
		err := easyjson.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			panic(err)
		}
		androidDetected := false
		msieDetected := false
		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				androidDetected = true
				browsers[browser] = true
			} else if strings.Contains(browser, "MSIE") {
				msieDetected = true
				browsers[browser] = true
			}
		}

		if !(androidDetected && msieDetected) {
			continue
		}
		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, strings.Replace(user.Email, "@", " [at] ", 1))
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Total unique browsers", len(browsers))
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson8f9bd62fDecodeGithubComNetworkPyCourseraLanguageBasicsHw3BenchForEasy(in *jlexer.Lexer, out *Users) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson8f9bd62fEncodeGithubComNetworkPyCourseraLanguageBasicsHw3BenchForEasy(out *jwriter.Writer, in Users) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix[1:])
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Users) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8f9bd62fEncodeGithubComNetworkPyCourseraLanguageBasicsHw3BenchForEasy(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Users) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8f9bd62fEncodeGithubComNetworkPyCourseraLanguageBasicsHw3BenchForEasy(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Users) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8f9bd62fDecodeGithubComNetworkPyCourseraLanguageBasicsHw3BenchForEasy(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Users) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8f9bd62fDecodeGithubComNetworkPyCourseraLanguageBasicsHw3BenchForEasy(l, v)
}

func LineCounter(r io.Reader) (int, error) {

	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}
