package tigeropen

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func buildSignContent(params map[string]interface{}) (string, error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(k)
		b.WriteByte('=')
		val := params[k]
		switch v := val.(type) {
		case string:
			b.WriteString(v)
		default:
			encoded, err := marshalWithSpaces(v)
			if err != nil {
				return "", fmt.Errorf("marshal %s: %w", k, err)
			}
			b.WriteString(encoded)
		}
	}
	return b.String(), nil
}

// marshalWithSpaces 将 JSON 编码，并在 ':' 和 ',' 后添加空格，以对齐 python 的 json.dumps 默认格式。
func marshalWithSpaces(v interface{}) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	inString := false
	escape := false
	for _, ch := range raw {
		c := rune(ch)
		if escape {
			out.WriteRune(c)
			escape = false
			continue
		}
		if c == '\\' {
			out.WriteRune(c)
			escape = true
			continue
		}
		if c == '"' {
			out.WriteRune(c)
			inString = !inString
			continue
		}
		if inString {
			out.WriteRune(c)
			continue
		}
		if c == ':' {
			out.WriteString(": ")
			continue
		}
		if c == ',' {
			out.WriteString(", ")
			continue
		}
		out.WriteRune(c)
	}
	return out.String(), nil
}

func signSHA1WithRSA(privateKey *rsa.PrivateKey, content []byte) (string, error) {
	hash := sha1.Sum(content)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// marshalDeterministic 生成按键排序、未转义 HTML 的 JSON，行为与 python 的 sort_keys + 紧凑分隔符一致。
func marshalDeterministic(v interface{}) (string, error) {
	var b strings.Builder
	if err := writeDeterministic(&b, v); err != nil {
		return "", err
	}
	return b.String(), nil
}

func writeDeterministic(b *strings.Builder, v interface{}) error {
	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(strconv.Quote(k))
			b.WriteByte(':')
			if err := writeDeterministic(b, val[k]); err != nil {
				return err
			}
		}
		b.WriteByte('}')
	case []interface{}:
		b.WriteByte('[')
		for i, item := range val {
			if i > 0 {
				b.WriteByte(',')
			}
			if err := writeDeterministic(b, item); err != nil {
				return err
			}
		}
		b.WriteByte(']')
	case json.RawMessage:
		b.Write(val)
	case string:
		b.WriteString(strconv.Quote(val))
	case nil:
		b.WriteString("null")
	default:
		// 处理非 []interface{} 的切片
		rv := reflect.ValueOf(v)
		if rv.IsValid() && rv.Kind() == reflect.Slice {
			tmp := make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				tmp[i] = rv.Index(i).Interface()
			}
			return writeDeterministic(b, tmp)
		}
		raw, err := json.Marshal(val)
		if err != nil {
			return err
		}
		b.Write(raw)
	}
	return nil
}
