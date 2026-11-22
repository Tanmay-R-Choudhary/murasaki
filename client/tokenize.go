package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp/data"
)

type token struct {
	text   string
	quoted bool
}

func tokenize(cmd string) ([]token, error) {
	var tokens []token
	var buf strings.Builder
	inQuotes, hasToken, wasQuoted := false, false, false

	flush := func() {
		if !hasToken {
			return
		}
		tokens = append(tokens, token{buf.String(), wasQuoted})
		buf.Reset()
		hasToken, wasQuoted = false, false
	}

	for i := 0; i < len(cmd); i++ {
		c := cmd[i]
		switch {
		case c == '\\' && inQuotes && i+1 < len(cmd) && (cmd[i+1] == '"' || cmd[i+1] == '\\'):
			i++
			buf.WriteByte(cmd[i])
		case c == '"':
			inQuotes = !inQuotes
			hasToken, wasQuoted = true, true
		case (c == ' ' || c == '\t') && !inQuotes:
			flush()
		default:
			buf.WriteByte(c)
			hasToken = true
		}
	}

	if inQuotes {
		return nil, errors.New("unterminated quote in command")
	}
	flush()
	return tokens, nil
}

func toData(t token) data.Data {
	if !t.quoted {
		if n, err := strconv.ParseInt(t.text, 10, 64); err == nil &&
			strconv.FormatInt(n, 10) == t.text {
			return data.NewMuInt(int(n))
		}
	}
	return data.NewMuString(t.text)
}

func strToMwpRepr(cmd string) (data.Data, error) {
	tokens, err := tokenize(cmd)
	if err != nil {
		return nil, err
	}
	elems := make([]data.Data, len(tokens))
	for i, t := range tokens {
		elems[i] = toData(t)
	}
	return data.NewMuArray(elems), nil
}
