package storage

import (
	"errors"
	"os"
	"strings"
	"testing"
)

type mockReaderWriter struct {
	data []byte
	err  error
	path string
	mode uint32
}

func (m *mockReaderWriter) ReadFile(path string) ([]byte, error) {
	return m.data, m.err
}

func (m *mockReaderWriter) WriteFile(name string, data []byte, perm os.FileMode) error {
	m.path = name
	m.mode = uint32(perm)
	m.data = data
	return m.err
}

func TestLoadFromDisk(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		err     error
		wantLen int
		wantNum int
		wantErr bool
	}{
		{
			name:    "file does not exist — create empty storage",
			data:    nil,
			err:     os.ErrNotExist,
			wantLen: 0,
			wantNum: 0,
			wantErr: false,
		},
		{
			name: "valid JSON — data loaded",
			data: []byte(`{
                "last_link_num": 2,
                "data": {"1": {"google.com": "available"}}
            }`),
			err:     nil,
			wantLen: 1,
			wantNum: 2,
			wantErr: false,
		},
		{
			name:    "corrupted JSON — error",
			data:    []byte(`{invalid json}`),
			err:     nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// Arrange
			reader := &mockReaderWriter{data: tc.data, err: tc.err}
			st, _ := NewStorage("test.json", reader)

			// Act
			err := st.LoadFromDisk()

			// Assert
			if tc.wantErr && err == nil {
				t.Fatalf("expected an error, but got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(st.Data) != tc.wantLen {
				t.Errorf("len(Data) = %d, want %d", len(st.Data), tc.wantLen)
			}

			if st.LastLinkNum != tc.wantNum {
				t.Errorf("LastLinkNum = %d, want %d", st.LastLinkNum, tc.wantNum)
			}
		})
	}
}

func TestSaveToDisk(t *testing.T) {

	tests := []struct {
		name           string
		initialData    map[int]map[string]string
		lastNum        int
		writeErr       error
		wantErr        bool
		wantPath       string
		wantJSONSubstr []string
	}{
		{
			name: "successful write",
			initialData: map[int]map[string]string{
				1: {"google.com": "available"},
			},
			lastNum:  1,
			writeErr: nil,
			wantErr:  false,
			wantPath: "test.json",
			wantJSONSubstr: []string{
				`"last_link_num": 1`,
				`"google.com": "available"`,
			},
		},
		{
			name:        "file write error",
			initialData: map[int]map[string]string{},
			lastNum:     0,
			writeErr:    errors.New("disk full"),
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// Arrange
			mock := &mockReaderWriter{
				err: tc.writeErr,
			}

			st, _ := NewStorage("test.json", mock)
			st.Data = tc.initialData
			st.LastLinkNum = tc.lastNum

			// Act
			err := st.SaveToDisk()

			// Assert
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if mock.path != tc.wantPath {
				t.Errorf("wrong path: got %s, want %s", mock.path, tc.wantPath)
			}

			content := string(mock.data)
			for _, substr := range tc.wantJSONSubstr {
				if !strings.Contains(content, substr) {
					t.Errorf("JSON does not contain '%s': %s", substr, content)
				}
			}
		})
	}
}
