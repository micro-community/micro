package logger

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/bytedance/sonic"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"

	"encoding"
	"unicode"
)

type Format int8

const (
	TextFormat Format = iota
	JsonFormat
)

// Keys for "built-in" attributes.
const (
	// TimeKey is the key used by the built-in handlers for the time
	// when the log method is called. The associated Value is a [time.Time].
	TimeKey = "ts"
	// PidKey is the key used by the built-in handlers for the pid
	// when the log level is lower than InfoLevel.
	PidKey = "pid"
	// MessageKey is the key used by the built-in handlers for the
	// message of the log call. The associated value is a string.
	MessageKey = "msg"
	// NanoTimeFieldFormat indicates the format of timestamp decoded
	// from a float value (time in seconds and nanoseconds).
	NanoTimeFieldFormat = "2006-01-02 15:04:05.999999999"
)

func (f Format) String() string {
	switch f {
	case TextFormat:
		return "text"
	case JsonFormat:
		return "json"
	default:
	}
	return ""
}

type Handler struct {
	handleState
}

// HandlerOptions are options for a TextHandler or Handler.
// A zero HandlerOptions consists entirely of default values.
type HandlerOptions struct {

	// Ignore records with levels below Level.Level().
	// The default is InfoLevel.
	Level slog.Leveler

	// ReplaceAttr can be used to change the default keys of the built-in
	// attributes, convert types (for example, to replace a `time.Time` with the
	// integer seconds since the Unix epoch), sanitize personal information, or
	// remove attributes from the output.
	ReplaceAttr func(a slog.Attr) slog.Attr
}

// NewJSONHandler creates a Handler that writes to w,
// using the default options.
func NewHandler(opts *Options) (h *Handler) {
	opts.Handler = (HandlerOptions{Level: slog.Level(opts.Level * 4), ReplaceAttr: func(a slog.Attr) slog.Attr {
		for k, v := range opts.Fields {
			if a.Key == k {
				return slog.Any(k, v)
			}
		}
		return a
	}}).NewHandler(opts.Out, opts.Format == JsonFormat)
	return opts.Handler.(*Handler)
}

// NewJSONHandler creates a Handler with the given options that writes to w.
func (opts HandlerOptions) NewHandler(w io.Writer, json bool) *Handler {
	return &Handler{
		handleState: handleState{buf: NewBuffer(),
			sep: "",
			h: &commonHandler{
				json: json,
				opts: opts,
				w:    bufio.NewWriter(w),
			},
		},
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handleState.h.enabled(level)
}

// With returns a new Handler whose attributes consists
// of h's attributes followed by attrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.handleState.h = h.handleState.h.withAttrs(attrs)
	return &Handler{h.handleState}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	h.handleState.h = h.handleState.h.withGroup(name)
	return &Handler{h.handleState}
}

// Handle formats its argument Record as a JSON object on a single line.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	return h.handleState.h.handle(r)
}

// Adapted from time.Time.MarshalJSON to avoid allocation.
func appendJSONTime(s *handleState, t time.Time) {
	s.buf.WriteByte('"')
	*s.buf = t.AppendFormat(*s.buf, NanoTimeFieldFormat)
	s.buf.WriteByte('"')
}

func appendJSONValue(s *handleState, v slog.Value) error {
	switch v.Kind() {
	case slog.KindString:
		s.appendString(v.String())
	case slog.KindInt64:
		*s.buf = strconv.AppendInt(*s.buf, v.Int64(), 10)
	case slog.KindUint64:
		*s.buf = strconv.AppendUint(*s.buf, v.Uint64(), 10)
	case slog.KindFloat64:
		f := v.Float64()
		// json.Marshal fails on special floats, so handle them here.
		switch {
		case math.IsInf(f, 1):
			s.buf.WriteString(`"+Inf"`)
		case math.IsInf(f, -1):
			s.buf.WriteString(`"-Inf"`)
		case math.IsNaN(f):
			s.buf.WriteString(`"NaN"`)
		default:
			if err := appendJSONMarshal(s.buf, f); err != nil {
				return err
			}
		}
	case slog.KindBool:
		*s.buf = strconv.AppendBool(*s.buf, v.Bool())
	case slog.KindDuration:
		// Do what json.Marshal does.
		*s.buf = strconv.AppendInt(*s.buf, int64(v.Duration()), 10)
	case slog.KindTime:
		s.appendTime(v.Time())
	case slog.KindAny:
		a := v.Any()
		if err, ok := a.(error); ok {
			s.appendString(err.Error())
		} else {
			return appendJSONMarshal(s.buf, a)
		}
	default:
		panic(fmt.Sprintf("bad kind: %d", v.Kind()))
	}
	return nil
}

func appendJSONMarshal(buf *Buffer, v any) error {
	b, err := sonic.Marshal(v)
	if err != nil {
		return err
	}
	buf.Write(b)
	return nil
}

// appendEscapedJSONString escapes s for JSON and appends it to buf.
func appendEscapedJSONString(buf *Buffer, s string) *Buffer {
	char := func(b byte) { buf.WriteByte(b) }
	str := func(s string) { buf.WriteString(s) }

	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				str(s[start:i])
			}
			char('\\')
			switch b {
			case '\\', '"':
				char(b)
			case '\n':
				char('n')
			case '\r':
				char('r')
			case '\t':
				char('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// It also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				str(`u00`)
				char(hex[b>>4])
				char(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				str(s[start:i])
			}
			str(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				str(s[start:i])
			}
			str(`\u202`)
			char(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		str(s[start:])
	}
	return buf
}

type commonHandler struct {
	json              bool // true => output JSON; false => output text
	opts              HandlerOptions
	preformattedAttrs []byte
	groupPrefix       string   // for text: prefix of groups opened in preformatting
	groups            []string // all groups started from WithGroup
	nOpenGroups       int      // the number of groups opened in preformattedAttrs
	mu                sync.Mutex
	w                 *bufio.Writer
}

func (h *commonHandler) clone() *commonHandler {
	// We can't use assignment because we can't copy the mutex.
	return &commonHandler{
		json:              h.json,
		opts:              h.opts,
		preformattedAttrs: h.preformattedAttrs,
		groupPrefix:       h.groupPrefix,
		groups:            slices.Clip(h.groups),
		nOpenGroups:       h.nOpenGroups,
		w:                 h.w,
	}
}

// Enabled reports whether l is greater than or equal to the
// minimum level.
func (h *commonHandler) enabled(l slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return l >= minLevel
}

func (h *commonHandler) withAttrs(as []slog.Attr) *commonHandler {
	h2 := h.clone()
	// Pre-format the attributes as an optimization.
	prefix := NewBuffer()
	defer prefix.Free()
	prefix.WriteString(h.groupPrefix)
	state := handleState{
		h:      h2,
		buf:    (*Buffer)(&h2.preformattedAttrs),
		sep:    "",
		prefix: prefix,
	}
	if len(h2.preformattedAttrs) > 0 {
		state.sep = h.attrSep()
	}
	state.openGroups()
	for _, a := range as {
		state.appendAttr(a)
	}
	// Remember the new prefix for later keys.
	h2.groupPrefix = state.prefix.String()
	// Remember how many opened groups are in preformattedAttrs,
	// so we don't open them again when we handle a Record.
	h2.nOpenGroups = len(h2.groups)
	return h2
}

func (h *commonHandler) withGroup(name string) *commonHandler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *commonHandler) handle(r slog.Record) error {

	//    t1 := time.Now()
	rep := h.opts.ReplaceAttr
	state := handleState{h: h, buf: NewBuffer(), sep: ""}
	state.buf.Reset()
	defer state.buf.Free()

	if h.json {
		state.buf.WriteByte('{')
	}

	// Built-in attributes. They are not in a group.
	// time
	if !r.Time.IsZero() {
		state.appendKey(TimeKey)
		state.appendTime(r.Time)

	}

	// adds one entry about pid information to attrs
	if r.Level < slog.LevelInfo {
		state.appendKey(PidKey)
		state.appendValue(slog.IntValue(os.Getpid()))
	}

	state.appendAttrs(r)

	if rep == nil {
		state.appendKey(MessageKey)
		state.appendString(r.Message)
	} else {
		state.appendAttr(slog.String(MessageKey, r.Message))
	}
	state.appendNonBuiltIns()
	state.buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.w.Write(*state.buf)
	defer h.w.Flush()

	return err
}

// attrSep returns the separator between attributes.
func (h *commonHandler) attrSep() string {
	if h.json {
		return ","
	}
	return " "
}

func (s *handleState) appendAttrs(r slog.Record) {
	// preformatted Attrs
	if len(s.h.preformattedAttrs) > 0 {
		s.buf.WriteString(s.sep)
		s.buf.Write(s.h.preformattedAttrs)
		s.sep = s.h.attrSep()
	}
	// Attrs in Record -- unlike the built-in ones, they are in groups started
	// from WithGroup.
	s.prefix = NewBuffer()
	defer s.prefix.Free()

	s.prefix.WriteString(s.h.groupPrefix)
	s.openGroups()
	r.Attrs(func(a slog.Attr) {
		s.appendAttr(a)
	})
}

func (s *handleState) appendNonBuiltIns() {

	if s.h.json {
		// Close all open groups.
		for range s.h.groups {
			s.buf.WriteByte('}')
		}
		// Close the top-level object.
		s.buf.WriteByte('}')
	}
}

// handleState holds state for a single call to commonHandler.handle.
// The initial value of sep determines whether to emit a separator
// before the next key, after which it stays true.
type handleState struct {
	h      *commonHandler
	buf    *Buffer
	sep    string  // separator to write before next key
	prefix *Buffer // for text: key prefix
}

func (s *handleState) openGroups() {
	for _, n := range s.h.groups[s.h.nOpenGroups:] {
		s.openGroup(n)
	}
}

// Separator for group names and keys.
const keyComponentSep = '.'

// openGroup starts a new group of attributes
// with the given name.
func (s *handleState) openGroup(name string) {
	if s.h.json {
		s.appendKey(name)
		s.buf.WriteByte('{')
		s.sep = ""
	} else {
		s.prefix.WriteString(name)
		s.prefix.WriteByte(keyComponentSep)
	}
}

// closeGroup ends the group with the given name.
func (s *handleState) closeGroup(name string) {
	if s.h.json {
		s.buf.WriteByte('}')
	} else {
		(*s.prefix) = (*s.prefix)[:len(*s.prefix)-len(name)-1 /* forkeyComponentSep */]
	}
	s.sep = s.h.attrSep()
}

// appendAttr appends the Attr's key and value using app.
// If sep is true, it also prepends a separator.
// It handles replacement and checking for an empty key.
// It sets sep to true if it actually did the append (if the key was non-empty
// after replacement).
func (s *handleState) appendAttr(a slog.Attr) {
	if rep := s.h.opts.ReplaceAttr; rep != nil {
		a = rep(a)
	}
	if a.Key == "" {
		return
	}
	v := a.Value.Resolve()
	if v.Kind() == slog.KindGroup {
		s.openGroup(a.Key)
		for _, aa := range v.Group() {
			s.appendAttr(aa)
		}
		s.closeGroup(a.Key)
	} else {
		s.appendKey(a.Key)
		s.appendValue(v)
	}
}

func (s *handleState) appendError(err error) {
	s.appendString(fmt.Sprintf("!ERROR:%v", err))
}

func (s *handleState) appendKey(key string) {
	s.buf.WriteString(s.sep)
	if s.prefix != nil {
		// TODO: optimize by avoiding allocation.
		s.appendString(string(*s.prefix) + key)
	} else {
		s.appendString(key)
	}
	if s.h.json {
		s.buf.WriteByte(':')
	} else {
		s.buf.WriteByte('=')
	}
	s.sep = s.h.attrSep()
}

func (s *handleState) appendString(str string) {
	if s.h.json {
		s.buf.WriteByte('"')
		s.buf = appendEscapedJSONString(s.buf, str)
		s.buf.WriteByte('"')
	} else {
		// text
		if needsQuoting(str) {
			*s.buf = strconv.AppendQuote(*s.buf, str)
		} else {
			s.buf.WriteString(str)
		}
	}
}

func (s *handleState) appendValue(v slog.Value) {
	var err error
	if s.h.json {
		err = appendJSONValue(s, v)
	} else {
		err = appendTextValue(s, v)
	}
	if err != nil {
		s.appendError(err)
	}
}

func (s *handleState) appendTime(t time.Time) {
	if s.h.json {
		appendJSONTime(s, t)
	} else {
		writeTimeRFC3339Millis(s.buf, t)
	}
}

// This takes half the time of Time.AppendFormat.
func writeTimeRFC3339Millis(buf *Buffer, t time.Time) {
	*buf = t.AppendFormat(*buf, NanoTimeFieldFormat)
}

func appendTextValue(s *handleState, v slog.Value) error {
	switch v.Kind() {
	case slog.KindString:
		s.appendString(v.String())
	case slog.KindTime:
		s.appendTime(v.Time())
	case slog.KindAny:
		if tm, ok := v.Any().(encoding.TextMarshaler); ok {
			data, err := tm.MarshalText()
			if err != nil {
				return err
			}
			s.appendString(string(data))
			return nil
		}
		s.appendString(fmt.Sprint(v.Any()))
	default:
		*s.buf = appendValue(v, *s.buf)
	}
	return nil
}

// append appends a text representation of v to dst.
// v is formatted as with fmt.Sprint.
func appendValue(v slog.Value, dst []byte) []byte {
	switch v.Kind() {
	case slog.KindString:
		return append(dst, v.String()...)
	case slog.KindInt64:
		return strconv.AppendInt(dst, v.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(dst, v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.AppendFloat(dst, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		return strconv.AppendBool(dst, v.Bool())
	case slog.KindDuration:
		return strconv.AppendInt(dst, int64(v.Duration()), 10)
	case slog.KindTime:
		return append(dst, v.Time().String()...)
	case slog.KindAny, slog.KindGroup, slog.KindLogValuer:
		return append(dst, fmt.Sprint(v.Any())...)
	default:
		panic(fmt.Sprintf("bad kind: %s", v.Kind()))
	}
}

func needsQuoting(s string) bool {
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			if needsQuotingSet[b] {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

var hex = "0123456789abcdef"

// Copied from encoding/json/encode.go:encodeState.string.
//
// htmlSafeSet holds the value true if the ASCII character with the given
// array position can be safely represented inside a JSON string, embedded
// inside of HTML <script> tags, without any additional escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), the backslash character ("\"), HTML opening and closing
// tags ("<" and ">"), and the ampersand ("&").
var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

var needsQuotingSet = [utf8.RuneSelf]bool{
	'"': true,
	'=': true,
}

func init() {
	for i := 0; i < utf8.RuneSelf; i++ {
		r := rune(i)
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			needsQuotingSet[i] = true
		}
	}
}
