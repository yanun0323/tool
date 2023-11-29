# Tool

A toolkit project provide some useful tool.

## Yaml Replacer
Change yaml files with custom rule in one command.

### Usage
Change the `tool/yaml_replacer.go` variables in the top of file.

Then run 
```bash 
$ make run.yaml.replacer 
```

#### Parameters
```go
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
```