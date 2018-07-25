package goossselect

import (
	"encoding/base64"
	"fmt"
	"errors"
)

func toBase64(s string) string {
	str := base64.StdEncoding.EncodeToString([]byte(s))
	return str
}

func SetFileHeaderInfo(s string) func(*Selector) error {
	return func(selector *Selector) error {
		if s == "Use" || s == "Ignore" || s == "None" {
			selector.FileHeaderInfo = s
			return nil
		}
		return errors.New("invalid value for FileHeaderInfo, only allowd: Use, Ignore, None")
	}
}

func SetInRecordDelimiter(s string) func(*Selector) error {
	return func(selector *Selector) error {
		str := toBase64(s)
		selector.InRecordDelimiter = str
		return nil
	}
}

func SetInFieldDelimiter(s string) func(*Selector) error {
	return func(selector *Selector) error {
		str := toBase64(s)
		selector.InFieldDelimiter = str
		return nil
	}
}

func SetOutRecordDelimiter(s string) func(*Selector) error {
	return func(selector *Selector) error {
		str := toBase64(s)
		selector.OutRecordDelimiter = str
		return nil
	}
}

func SetOutFieldDelimiter(s string) func(*Selector) error {
	return func(selector *Selector) error {
		str := toBase64(s)
		selector.OutFieldDelimiter = str
		return nil
	}
}

func SetOutputRawData(b bool) func(*Selector) error {
	return func(selector *Selector) error {
		selector.OutputRawData = b
		return nil
	}
}

func SetKeepAllColumns(b bool) func(*Selector) error {
	return func(selector *Selector) error {
		selector.KeepAllColumns = b
		return nil
	}
}

func SetQuoteCharacter(s string) func(*Selector) error {
	return func(selector *Selector) error {
		str := toBase64(s)
		selector.QuoteCharacter = str
		return nil
	}
}

func SetCommentCharacter(s string) func(*Selector) error {
	return func(selector *Selector) error {
		str := toBase64(s)
		selector.CommentCharacter = str
		return nil
	}
}

func SetRange(s string, start, end int) func(*Selector) error {
	return func(selector *Selector) error {
		if s == "line" || s == "split" {
			str := fmt.Sprintf("%s-range=%d-%d", s, start, end)
			selector.Range = str
			return nil
		}
		return errors.New("invalid value for Range, only allowd: line, split")
	}
}

func SetMetaOverwriteIfExisting(b bool) func(*Meta) error {
	return func(m *Meta) error {
		m.OverwriteIfExisting = b
		return nil
	}
}

func SetMetaRecordDelimiter(s string) func(*Meta) error {
	str := toBase64(s)
	return func(m *Meta) error {
		m.RecordDelimiter = str
		return nil
	}
}

func SetMetaFieldDelimiter(s string) func(*Meta) error {
	str := toBase64(s)
	return func(m *Meta) error {
		m.FieldDelimiter = str
		return nil
	}
}

func SetMetaQuoteCharacter(s string) func(*Meta) error {
	str := toBase64(s)
	return func(m *Meta) error {
		m.QuoteCharacter = str
		return nil
	}
}