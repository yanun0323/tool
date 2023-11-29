package tool

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/yanun0323/pkg/logs"
)

var (
	_dir                      = "/Users/hqcc-user14/Documents/Project/Helm/"              /* 目標專案的 absolute path */
	_includeRelativeDirRegexp = regexp.MustCompile("^esc-*")                              /* (包含)目標專案第一層子資料夾 */
	_excludeRelativeDirRegexp = regexp.MustCompile("(esc-bo-gateway|esc-web|esc-wallet)") /* (排除)目標專案第一層子資料夾 */
	_includeFilePathRegexp    = regexp.MustCompile(".*test.yaml$")                        /* (包含)檔案路徑 */
	_excludeFilePathRegexp    = regexp.MustCompile(".*helmignore$")                       /* (排除)檔案路徑 */
	_targetParentKey          = "db:"                                                     /* 要更改的 Key 的前綴 */
	_targetKeys               = "maxIdleConns,maxOpenConns"                               /* 要更改的 Key (支援多個，使用','隔開) */
	_targetValues             = "10,10"                                                   /* 對應上方順序，要更改的數值 (支援多個，使用','隔開) */
)

func YamlReplacer(ctx context.Context) {
	l := logs.Get(ctx).WithFunc("HelmReplacer")

	if err := validateParameters(); err != nil {
		l.WithError(err).Error("invalid parameters")
		return
	}

	outerErr := filepath.WalkDir(_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errors.New(fmt.Sprintf("walk dir, err: %+v", err))
		}

		relativePath := strings.TrimPrefix(path, _dir)
		if !isTargetRelativeDir(relativePath) {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return errors.New(fmt.Sprintf("open file, err: %+v", err))
		}
		defer f.Close()

		if !isTargetFile(f.Name()) {
			return nil
		}

		l.Infof("read file %s", f.Name())

		data, err := io.ReadAll(f)
		if err != nil {
			return errors.New(fmt.Sprintf("read all data from file, err: %+v", err))
		}

		targetKeys := strings.Split(_targetKeys, ",")
		if len(targetKeys) == 0 {
			return errors.New("empty target keys")
		}

		targetValues := strings.Split(_targetValues, ",")
		if len(targetValues) != len(targetKeys) {
			return errors.New("target keys count does not match target keys")
		}

		targetKeysFound := make([]bool, len(targetKeys))
		targetParentKey := _targetParentKey
		if !strings.HasSuffix(targetParentKey, ":") {
			targetParentKey = targetParentKey + ":"
		}

		keyParentFound := false
		keyParentTab := 0

		rows := strings.Split(string(data), "\n")
		for i, row := range rows {
			tab := countTab(row)
			r := strings.TrimSpace(row)
			if tab < keyParentTab {
				keyParentFound = false
			}

			if keyParentFound {
				for j, found := range targetKeysFound {
					if found {
						continue
					}
					if strings.HasPrefix(r, targetKeys[j]) {
						l.Debugf("change %s value into %s", r, targetValues[j])
						sp := strings.Split(row, ":")
						sp[1] = targetValues[j]
						rows[i] = strings.Join(sp, ": ")
						targetKeysFound[j] = true
					}
				}

				if tab > keyParentTab {
					continue
				}
			}

			if r == targetParentKey {
				keyParentFound = true
			}

			keyParentTab = tab
		}

		changed := false
		for _, found := range targetKeysFound {
			changed = changed || found
		}

		if !changed {
			l.Warn("keywords not found")
			return nil
		}

		{
			f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return errors.New(fmt.Sprintf("open file for edit, err: %+v", err))
			}
			defer f.Close()

			l.Debug("write data")
			if _, err := f.WriteString(strings.Join(rows, "\n")); err != nil {
				return errors.New(fmt.Sprintf("write data, err: %+v", err))
			}
		}

		l.Info("data changed")

		return nil
	})
	if outerErr != nil {
		l.WithError(outerErr).Fatalf("walk dir, err: %+v", outerErr)
	}
}

func validateParameters() error {
	if len(_targetParentKey) == 0 {
		return errors.New("require _targetKeyParent variable")
	}

	if len(_targetKeys) == 0 {
		return errors.New("require _targetKeys variable")
	}

	if len(_targetValues) == 0 {
		return errors.New("require _targetValue variable")
	}

	return nil
}

func isTargetRelativeDir(relativeDir string) bool {
	return _includeRelativeDirRegexp.MatchString(relativeDir) &&
		!_excludeRelativeDirRegexp.MatchString(relativeDir)
}

func isTargetFile(filename string) bool {
	return _includeFilePathRegexp.MatchString(filename) &&
		!_excludeFilePathRegexp.MatchString(filename)

}

func countTab(s string) int {
	c := 0
	for _, b := range s {
		if b != ' ' {
			return c
		}
		c++
	}
	return c
}
