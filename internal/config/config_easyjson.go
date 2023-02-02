// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig(in *jlexer.Lexer, out *Config) {
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
		case "url":
			out.URL = string(in.String())
		case "api_key":
			out.APIKey = string(in.String())
		case "api_secret":
			out.APISecret = string(in.String())
		case "labels":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Labels = make(map[string]string)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.Labels)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "rate":
			out.Rate = int(in.Int())
		case "timeout":
			d, err := time.ParseDuration(in.String())
			if err != nil {
				out.Timeout = time.Second * 30
			}
			out.Timeout = d
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
func easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig(out *jwriter.Writer, in Config) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.URL))
	}
	{
		const prefix string = ",\"api_key\":"
		out.RawString(prefix)
		out.String(string(in.APIKey))
	}
	{
		const prefix string = ",\"api_secret\":"
		out.RawString(prefix)
		out.String(string(in.APISecret))
	}
	{
		const prefix string = ",\"labels\":"
		out.RawString(prefix)
		if in.Labels == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Labels {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.String(string(v2Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"rate\":"
		out.RawString(prefix)
		out.Int(int(in.Rate))
	}
	{
		const prefix string = ",\"timeout\":"
		out.RawString(prefix)
		out.Int64(int64(in.Timeout))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Config) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Config) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Config) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Config) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig(l, v)
}