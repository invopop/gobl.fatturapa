package fatturapa_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"encoding/xml"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGOBLToXMLExamples(t *testing.T) {
	schema, err := test.LoadSchema()
	require.NoError(t, err)

	path := test.GetDataPath() + "/" + test.PathGOBLFatturaPA

	var files []string
	err = filepath.Walk(path, func(path string, _ os.FileInfo, _ error) error {
		if filepath.Ext(path) == ".json" {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	require.NoError(t, err)

	for _, file := range files {
		fmt.Printf("processing file: %v\n", file)

		env := test.LoadTestFile(file)

		doc, err := test.ConvertFromGOBL(env, test.NewConverter())
		require.NoError(t, err)

		data, err := xml.MarshalIndent(doc, "", "\t")
		require.NoError(t, err)

		np := strings.TrimSuffix(file, filepath.Ext(file)) + ".xml"
		outPath := filepath.Join(test.GetDataPath(), "out", np)

		if *test.UpdateOut {
			errs := test.ValidateXML(schema, data)
			for _, e := range errs {
				assert.NoError(t, e)
			}
			if len(errs) > 0 {
				assert.Fail(t, "Invalid XML:\n"+string(data))
			}

			err = os.WriteFile(outPath, data, 0644)
			require.NoError(t, err, "writing file")
		}

		expected, err := os.ReadFile(outPath)

		require.False(t, os.IsNotExist(err), "output file %s missing, run tests with `--update` flag to create", filepath.Base(outPath))
		require.NoError(t, err)
		require.Equal(t, string(expected), string(data), "output file %s does not match, run tests with `--update` flag to update", filepath.Base(outPath))
	}
}

func TestXMLToGOBLExamples(t *testing.T) {
	var err error
	var files []string

	path := test.GetDataPath() + "/" + test.PathFatturaPAGOBL

	err = filepath.Walk(path, func(path string, _ os.FileInfo, _ error) error {
		if filepath.Ext(path) == ".xml" {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	require.NoError(t, err)

	for _, file := range files {
		fmt.Printf("processing file: %v\n", file)

		data, err := os.ReadFile(filepath.Join(test.GetDataPath(), file))
		require.NoError(t, err)

		env, err := test.ConvertToGOBL(data, test.NewConverter())
		require.NoError(t, err)
		require.NotNil(t, env)

		require.NoError(t, env.Calculate())

		data, err = json.MarshalIndent(env, "", "\t")
		require.NoError(t, err)

		np := strings.TrimSuffix(file, filepath.Ext(file)) + ".json"
		outPath := filepath.Join(test.GetDataPath(), "out", np)

		//if *test.UpdateOut {
		err = env.Validate()
		require.NoError(t, err)
		//assert.Fail(t, "Invalid GOBL:\n"+string(err.Error()))

		err = os.WriteFile(outPath, data, 0644)
		require.NoError(t, err, "writing file")
		//}

		expected, err := os.ReadFile(outPath)
		require.False(t, os.IsNotExist(err), "output file %s missing, run tests with `--update` flag to create", filepath.Base(outPath))
		require.NoError(t, err)
		require.Equal(t, string(expected), string(data), "output file %s does not match, run tests with `--update` flag to update", filepath.Base(outPath))
	}
}
