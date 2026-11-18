package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		in      string
		wantOut string
		wantErr string
		code    int
	}{
		{
			name:    "basic",
			args:    []string{"1250000"},
			wantOut: "satu juta dua ratus lima puluh ribu\n",
		},
		{
			name:    "several",
			args:    []string{"0", "100", "-1"},
			wantOut: "nol\nseratus\nminus satu\n",
		},
		{
			name:    "negative looks like a flag",
			args:    []string{"-11"},
			wantOut: "minus sebelas\n",
		},
		{
			name:    "end of flags",
			args:    []string{"--", "-11"},
			wantOut: "minus sebelas\n",
		},
		{
			name:    "rupiah",
			args:    []string{"-rupiah", "1000"},
			wantOut: "seribu rupiah\n",
		},
		{
			name:    "rupiah after the number",
			args:    []string{"1000", "-rupiah"},
			wantOut: "seribu rupiah\n",
		},
		{
			name:    "format",
			args:    []string{"-format", "1250000"},
			wantOut: "Rp1.250.000,00\n",
		},
		{
			name:    "format fraction",
			args:    []string{"-format", "1,5"},
			wantOut: "Rp1,50\n",
		},
		{
			name:    "both",
			args:    []string{"-format", "-rupiah", "1250000"},
			wantOut: "Rp1.250.000,00\nsatu juta dua ratus lima puluh ribu rupiah\n",
		},
		{
			name:    "both several",
			args:    []string{"-format", "-rupiah", "1", "2"},
			wantOut: "Rp1,00\nsatu rupiah\n\nRp2,00\ndua rupiah\n",
		},
		{
			name:    "sen words",
			args:    []string{"-rupiah", "1,5"},
			wantOut: "satu rupiah lima puluh sen\n",
		},
		{
			name:    "grouped",
			args:    []string{"1.250.000"},
			wantOut: "satu juta dua ratus lima puluh ribu\n",
		},
		{
			name:    "stdin",
			in:      "11\n12\n",
			wantOut: "sebelas\ndua belas\n",
		},
		{
			name:    "stdin rupiah spacing",
			in:      "Rp 1.250\n",
			wantOut: "seribu dua ratus lima puluh\n",
		},
		{
			name:    "stdin several on one line",
			in:      "11 12\n",
			wantOut: "sebelas\ndua belas\n",
		},
		{
			name:    "stdin skips comments",
			in:      "# catatan\n\n11\n",
			wantOut: "sebelas\n",
		},
		{
			name:    "stdin format",
			args:    []string{"-format"},
			in:      "1000\n",
			wantOut: "Rp1.000,00\n",
		},
		{
			name:    "decimal",
			args:    []string{"1,05"},
			wantOut: "satu koma nol lima\n",
		},
		{
			name:    "version",
			args:    []string{"-version"},
			wantOut: "terbilang dev\n",
		},
		{
			name:    "help",
			args:    []string{"-h"},
			wantErr: "Usage:",
		},
		{
			name:    "empty stdin",
			wantErr: "Usage:",
			code:    2,
		},
		{
			name:    "invalid",
			args:    []string{"abc"},
			wantErr: "invalid number",
			code:    1,
		},
		{
			name:    "invalid number that looks negative",
			args:    []string{"-1.2.3"},
			wantErr: "invalid number",
			code:    1,
		},
		{
			name:    "bad flag",
			args:    []string{"-nope"},
			wantErr: "flag",
			code:    2,
		},
		{
			name:    "rp without an amount",
			args:    []string{"Rp"},
			wantErr: "invalid number",
			code:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			code := run(tt.args, strings.NewReader(tt.in), &stdout, &stderr)
			if code != tt.code {
				t.Fatalf("exit %d, want %d\nstderr: %s", code, tt.code, stderr.String())
			}
			if stdout.String() != tt.wantOut {
				t.Fatalf("stdout = %q, want %q", stdout.String(), tt.wantOut)
			}
			if tt.wantErr == "" && stderr.Len() != 0 {
				t.Fatalf("stderr = %q", stderr.String())
			}
			if tt.wantErr != "" && !strings.Contains(stderr.String(), tt.wantErr) {
				t.Fatalf("stderr = %q, want it to contain %q", stderr.String(), tt.wantErr)
			}
		})
	}
}

type failIO struct{}

func (failIO) Read([]byte) (int, error)  { return 0, errors.New("read failed") }
func (failIO) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestRunIOErrors(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	if code := run(nil, failIO{}, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "read failed") {
		t.Fatalf("read error: code=%d stderr=%q", code, stderr.String())
	}

	stderr.Reset()
	if code := run([]string{"1"}, strings.NewReader(""), failIO{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "write failed") {
		t.Fatalf("write error: code=%d stderr=%q", code, stderr.String())
	}
}
