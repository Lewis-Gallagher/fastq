package fastq

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Record struct {
	ReadId string
	Info   map[string]string
	Seq    []byte
	Qual   []byte
}

type Parser struct {
	reader *bufio.Reader
	line   uint
	atEOF  bool
}

// NewParser takes io.Reader so should be able to accept many input formats
// such as stdin, gzip readers, strings etc.
func NewParser(r io.Reader) *Parser {
	return &Parser{reader: bufio.NewReader(r)}
}

func (p *Parser) readLine() ([]byte, error) {

	line, err := p.reader.ReadBytes('\n')

	// ReadBytes returns io.EOF together with data if the file's last line has no trailing \n. In this case the data is still valid
	// So only consider io.EOF a valid error if no data was read
	if err != nil && (err != io.EOF || len(line) == 0) {
		return nil, err
	}

	// Increment line counter
	p.line++

	// Strip
	return bytes.TrimRight(line, "\r\n"), nil
}

// Parses by line with readLine()
func (p *Parser) Read() (Record, error) {

	if p.atEOF {
		return Record{}, io.EOF
	}

	// Parse header
	header, err := p.readLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			p.atEOF = true
		}
		return Record{}, err
	}

	headerStr := string(header)

	if len(header) == 0 || !strings.HasPrefix(headerStr, "@") {
		return Record{}, fmt.Errorf("line %d: header must start with '@'", p.line)
	}

	// Read ID
	// readId, err := GetReadId(headerStr)
	readId := strings.TrimPrefix(strings.Split(headerStr, " ")[0], "@")

	// Everything else
	info := ParseHeaderInfoPlainText(headerStr)

	// Parse sequence
	seq, err := p.readLine()
	if err != nil {
		return Record{}, p.malformed(err)
	}

	// Parse (skip) '+' separater
	plus, err := p.readLine()
	if err != nil {
		return Record{}, err
	}

	if len(plus) == 0 || plus[0] != '+' {
		return Record{}, fmt.Errorf("line %d: no '+' separator found", p.line)
	}

	// Parse quality
	qual, err := p.readLine()
	if err != nil {
		return Record{}, p.malformed(err)
	}

	if len(seq) != len(qual) {
		return Record{}, fmt.Errorf("sequence length (len=%d) and quality length (len=%d) mismatch on line %d",
			len(seq), len(qual), p.line)
	}

	return Record{ReadId: readId, Seq: seq, Qual: qual, Info: info}, nil
}

func ParseHeaderInfoPlainText(header string) map[string]string {

	headerParts := strings.Split(header, " ")

	info := make(map[string]string)

	// Skip first ele (the read ID)
	for _, token := range headerParts[1:] {
		infoKey, infoValue, ok := strings.Cut(token, "=")
		if ok {
			info[infoKey] = infoValue
		}
	}

	return info
}

func (p *Parser) malformed(err error) error {
	if err == io.EOF {
		return fmt.Errorf("line %d: %w", p.line, io.ErrUnexpectedEOF)
	}
	return fmt.Errorf("line %d: %w", p.line, err)
}
