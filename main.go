package main

import (
	"context"
	"main/tool"
	"os"

	"github.com/yanun0323/pkg/logs"
)

func main() {
	fn := os.Getenv("FN")
	ctx, l := logs.New(logs.LevelInfo).WithField("FN", fn).Attach(context.Background())

	switch fn {
	case "YAML_REPLACER":
		tool.YamlReplacer(ctx)
	default:
		l.Error("Unsupported 'FN' value")
	}
}
