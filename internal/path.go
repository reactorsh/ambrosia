package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

func genOutPath(cmd string, args []string) string {
	if len(args) > 1 {
		return args[1]
	}

	inPath := args[0]
	dir := filepath.Dir(inPath)
	ext := filepath.Ext(inPath)
	filename := filepath.Base(inPath)
	filename = strings.TrimSuffix(filename, ext)

	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", filename, cmd, ext))
}

func searchPrefix(infilePath string) string {
	infileBase := filepath.Base(infilePath)
	infileExt := filepath.Ext(infileBase)
	infileName := strings.TrimSuffix(infileBase, infileExt)
	return fmt.Sprintf("%s_", infileName)
}

func prefixPathTmpl(infilePath string) string {
	infileBase := filepath.Base(infilePath)
	infileExt := filepath.Ext(infileBase)
	infileName := strings.TrimSuffix(infileBase, infileExt)
	return filepath.Join(filepath.Dir(infilePath), fmt.Sprintf("%s_psort_%%c%s", infileName, infileExt))
}
