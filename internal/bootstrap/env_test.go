package bootstrap

import "testing"

func TestParseEnvLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantKey string
		wantVal string
		wantOK  bool
		wantErr bool
	}{
		{name: "blank", line: "   ", wantOK: false},
		{name: "comment", line: "# x", wantOK: false},
		{name: "simple", line: "A=B", wantKey: "A", wantVal: "B", wantOK: true},
		{name: "quoted", line: `A="B C"`, wantKey: "A", wantVal: "B C", wantOK: true},
		{name: "export", line: "export A=1", wantKey: "A", wantVal: "1", wantOK: true},
		{name: "invalid", line: "A", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, val, ok, err := parseEnvLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected err state: %v", err)
			}
			if ok != tt.wantOK {
				t.Fatalf("expected ok=%v got %v", tt.wantOK, ok)
			}
			if key != tt.wantKey || val != tt.wantVal {
				t.Fatalf("unexpected key/value: %q/%q", key, val)
			}
		})
	}
}
