package process

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_processTask(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		testName string
		filename string
		paths    []string
		exp      string
		toMatch  int
		err      error
	}{
		{
			testName: "just a file",
			paths: []string{
				filepath.Join("0.txt"),
			},
			exp:     "0",
			toMatch: 1,
			err:     nil,
		},
		{
			testName: "file one deep",
			paths: []string{
				filepath.Join("1", "1.txt"),
			},
			toMatch: 1,
			exp:     "1",
		},
		{
			testName: "file two deep",
			paths: []string{
				filepath.Join("1", "2", "2.txt"),
			},
			toMatch: 1,
			exp:     "1",
		},
		{
			testName: "file three deep",
			paths: []string{
				filepath.Join("1", "2", "3", "3.txt"),
			},
			toMatch: 1,
			exp:     "1",
		},
		{
			testName: "files 0 and 1 deep",
			paths: []string{
				filepath.Join("this.txt"),
				filepath.Join("1", "2.txt"),
			},
			toMatch: 2,
			exp:     ".+",
		},
		{
			testName: "files at every level",
			paths: []string{
				filepath.Join("0.txt"),
				filepath.Join("1", "1.txt"),
				filepath.Join("1", "2", "2.txt"),
				filepath.Join("1", "2", "3", "3.txt"),
			},
			toMatch: 4,
			exp:     `\d+`,
		},
		{
			testName: "filter some with regex",
			paths: []string{
				filepath.Join("thiw.txt"),
				filepath.Join("1", "1.txt"),
				filepath.Join("1", "2", "2.txt"),
				filepath.Join("1", "2", "3", "3.txt"),
			},
			toMatch: 3,
			exp:     "1",
		},
		{
			testName: "filter some more",
			paths: []string{
				filepath.Join("thiw.txt"),
				filepath.Join("1", "1.txt"),
				filepath.Join("1", "2", "2.txt"),
				filepath.Join("1", "2", "3", "3.txt"),
				filepath.Join("2", "1.txt"),
				filepath.Join("2", "2", "2.txt"),
				filepath.Join("2", "2", "3", "3.txt"),
			},
			toMatch: 3,
			exp:     "1",
		},
		{
			testName: "all accepted",
			paths: []string{
				filepath.Join("thiw.txt"),
				filepath.Join("1", "1.txt"),
				filepath.Join("1", "2", "2.txt"),
				filepath.Join("1", "2", "3", "3.txt"),
				filepath.Join("2", "1.txt"),
				filepath.Join("2", "2", "2.txt"),
				filepath.Join("2", "2", "3", "3.txt"),
			},
			toMatch: 7,
			exp:     ".*",
		},
		{
			testName: "bogus regex",
			paths: []string{
				filepath.Join("this.txt"),
			},
			toMatch: 0,
			exp:     "(",
			err:     fmt.Errorf("error parsing regexp: missing closing ): `(`"),
		},
	}

	for _, tt := range tests {
		dir, err := os.MkdirTemp("", "*")
		require.NoError(t, err)
		defer os.RemoveAll(tt.filename)
		destination, err := os.MkdirTemp("", "*")
		require.NoError(t, err)
		defer os.RemoveAll(destination)

		for i, path := range tt.paths {
			if strings.Contains(path, "/") {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, filepath.Dir(path)), 0755))
			}
			require.NoError(t, os.WriteFile(filepath.Join(dir, path), []byte(fmt.Sprintf("%d.txt", i)), 0644))
		}

		got, err := GetMatches(ctx, Task{
			Name:                 tt.testName,
			SearchDirectory:      dir,
			SearchPattern:        tt.exp,
			DestinationDirectory: destination,
		})
		if err != nil {
			assert.Equal(t, tt.err.Error(), err.Error())
		}
		assert.Equal(t, tt.toMatch, len(got))
	}
}
