// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

import (
	json "encoding/json"
	_v6 "github.com/brianvoe/gofakeit/v6"
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

func easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig(in *jlexer.Lexer, out *LogInfo) {
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
		case "display":
			out.Display = string(in.String())
		case "category":
			out.Category = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "example":
			out.Example = string(in.String())
		case "params":
			if in.IsNull() {
				in.Skip()
				out.Params = nil
			} else {
				in.Delim('[')
				if out.Params == nil {
					if !in.IsDelim(']') {
						out.Params = make([]_v6.Param, 0, 0)
					} else {
						out.Params = []_v6.Param{}
					}
				} else {
					out.Params = (out.Params)[:0]
				}
				for !in.IsDelim(']') {
					var v1 _v6.Param
					easyjson6615c02eDecodeGithubComBrianvoeGofakeitV6(in, &v1)
					out.Params = append(out.Params, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig(out *jwriter.Writer, in LogInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"display\":"
		out.RawString(prefix[1:])
		out.String(string(in.Display))
	}
	{
		const prefix string = ",\"category\":"
		out.RawString(prefix)
		out.String(string(in.Category))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"example\":"
		out.RawString(prefix)
		out.String(string(in.Example))
	}
	{
		const prefix string = ",\"params\":"
		out.RawString(prefix)
		if in.Params == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Params {
				if v2 > 0 {
					out.RawByte(',')
				}
				easyjson6615c02eEncodeGithubComBrianvoeGofakeitV6(out, v3)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LogInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LogInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LogInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LogInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig(l, v)
}
func easyjson6615c02eDecodeGithubComBrianvoeGofakeitV6(in *jlexer.Lexer, out *_v6.Param) {
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
		case "field":
			out.Field = string(in.String())
		case "display":
			out.Display = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "optional":
			out.Optional = bool(in.Bool())
		case "default":
			out.Default = string(in.String())
		case "options":
			if in.IsNull() {
				in.Skip()
				out.Options = nil
			} else {
				in.Delim('[')
				if out.Options == nil {
					if !in.IsDelim(']') {
						out.Options = make([]string, 0, 4)
					} else {
						out.Options = []string{}
					}
				} else {
					out.Options = (out.Options)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Options = append(out.Options, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "description":
			out.Description = string(in.String())
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
func easyjson6615c02eEncodeGithubComBrianvoeGofakeitV6(out *jwriter.Writer, in _v6.Param) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"field\":"
		out.RawString(prefix[1:])
		out.String(string(in.Field))
	}
	{
		const prefix string = ",\"display\":"
		out.RawString(prefix)
		out.String(string(in.Display))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"optional\":"
		out.RawString(prefix)
		out.Bool(bool(in.Optional))
	}
	{
		const prefix string = ",\"default\":"
		out.RawString(prefix)
		out.String(string(in.Default))
	}
	{
		const prefix string = ",\"options\":"
		out.RawString(prefix)
		if in.Options == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Options {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	out.RawByte('}')
}
func easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig1(in *jlexer.Lexer, out *LogConfig) {
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
		case "format":
			out.Format = string(in.String())
		case "structure":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Structure = make(map[string]string)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v7 string
					v7 = string(in.String())
					(out.Structure)[key] = v7
					in.WantComma()
				}
				in.Delim('}')
			}
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
func easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig1(out *jwriter.Writer, in LogConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"format\":"
		out.RawString(prefix[1:])
		out.String(string(in.Format))
	}
	{
		const prefix string = ",\"structure\":"
		out.RawString(prefix)
		if in.Structure == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v8First := true
			for v8Name, v8Value := range in.Structure {
				if v8First {
					v8First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v8Name))
				out.RawByte(':')
				out.String(string(v8Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LogConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LogConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LogConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LogConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig1(l, v)
}
func easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig2(in *jlexer.Lexer, out *DetailedLogConfig) {
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
		case "format":
			out.Format = string(in.String())
		case "structure":
			if in.IsNull() {
				in.Skip()
				out.Structure = nil
			} else {
				in.Delim('[')
				if out.Structure == nil {
					if !in.IsDelim(']') {
						out.Structure = make([]LogInfo, 0, 0)
					} else {
						out.Structure = []LogInfo{}
					}
				} else {
					out.Structure = (out.Structure)[:0]
				}
				for !in.IsDelim(']') {
					var v9 LogInfo
					(v9).UnmarshalEasyJSON(in)
					out.Structure = append(out.Structure, v9)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig2(out *jwriter.Writer, in DetailedLogConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"format\":"
		out.RawString(prefix[1:])
		out.String(string(in.Format))
	}
	{
		const prefix string = ",\"structure\":"
		out.RawString(prefix)
		if in.Structure == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v10, v11 := range in.Structure {
				if v10 > 0 {
					out.RawByte(',')
				}
				(v11).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v DetailedLogConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DetailedLogConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *DetailedLogConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DetailedLogConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig2(l, v)
}
func easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig3(in *jlexer.Lexer, out *Config) {
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
					var v12 string
					v12 = string(in.String())
					(out.Labels)[key] = v12
					in.WantComma()
				}
				in.Delim('}')
			}
		case "rate":
			out.Rate = int(in.Int())
		case "timeout":
			d, err := time.ParseDuration(in.String())
			if err != nil {
				d = time.Second * 30
			}
			out.Timeout = d
		case "log_config":
			(out.LogConfig).UnmarshalEasyJSON(in)
		case "enable_metrics":
			out.EnableMetrics = bool(in.Bool())
		case "enable_traces":
			out.EnableTraces = bool(in.Bool())
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
func easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig3(out *jwriter.Writer, in Config) {
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
			v13First := true
			for v13Name, v13Value := range in.Labels {
				if v13First {
					v13First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v13Name))
				out.RawByte(':')
				out.String(string(v13Value))
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
		out.String(time.Duration(int64(in.Timeout)).String())
	}
	{
		const prefix string = ",\"log_config\":"
		out.RawString(prefix)
		(in.LogConfig).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"enable_metrics\":"
		out.RawString(prefix)
		out.Bool(bool(in.EnableMetrics))
	}
	{
		const prefix string = ",\"enable_traces\":"
		out.RawString(prefix)
		out.Bool(bool(in.EnableTraces))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Config) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Config) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComGigapipehqLoggenInternalConfig3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Config) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Config) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComGigapipehqLoggenInternalConfig3(l, v)
}
