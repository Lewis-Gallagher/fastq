package fastq

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func readAll(p *Parser) ([]Record, error) {
	var recs []Record
	for {
		rec, err := p.Read()
		if err == io.EOF {
			return recs, nil
		}
		if err != nil {
			return recs, err
		}
		recs = append(recs, rec)
	}
}

func TestParserRead(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Record // only ReadId, Seq, Qual are compared
		wantErr bool
		wantIs  error // optional: the error must wrap this
	}{
		{
			name:  "valid single record",
			input: "@r1 ch=1\nACGT\n+\nIIII\n",
			want:  []Record{{ReadId: "r1", Seq: []byte("ACGT"), Qual: []byte("IIII")}},
		},
		{
			name:  "no trailing newline",
			input: "@r1\nACGT\n+\nIIII",
			want:  []Record{{ReadId: "r1", Seq: []byte("ACGT"), Qual: []byte("IIII")}},
		},
		{
			name:  "windows line endings",
			input: "@r1\r\nACGT\r\n+\r\nIIII\r\n",
			want:  []Record{{ReadId: "r1", Seq: []byte("ACGT"), Qual: []byte("IIII")}},
		},
		{
			name:  "quality starts with @, followed by another record",
			input: "@r1\nACGT\n+\n@III\n@r2\nTTTT\n+\nIIII\n",
			want: []Record{
				{ReadId: "r1", Seq: []byte("ACGT"), Qual: []byte("@III")},
				{ReadId: "r2", Seq: []byte("TTTT"), Qual: []byte("IIII")},
			},
		},
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{name: "missing @", input: "r1\nACGT\n+\nIIII\n", wantErr: true},
		{name: "missing +", input: "@r1\nACGT\nXX\nIIII\n", wantErr: true},
		{name: "length mismatch", input: "@r1\nACGT\n+\nIII\n", wantErr: true},
		{name: "truncated after header", input: "@r1\n", wantErr: true, wantIs: io.ErrUnexpectedEOF},
		{name: "truncated after sequence", input: "@r1\nACGT\n", wantErr: true, wantIs: io.ErrUnexpectedEOF},
		{name: "truncated after plus", input: "@r1\nACGT\n+\n", wantErr: true, wantIs: io.ErrUnexpectedEOF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			got, err := readAll(p)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got none")
				}
				if tt.wantIs != nil && !errors.Is(err, tt.wantIs) {
					t.Fatalf("error %v does not wrap %v", err, tt.wantIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d records, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].ReadId != tt.want[i].ReadId ||
					string(got[i].Seq) != string(tt.want[i].Seq) ||
					string(got[i].Qual) != string(tt.want[i].Qual) {
					t.Errorf("record %d: got %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
