**使用:**


```golang
package main

import (
	"github.com/GamerCode/logx"
)

func main() {
	conf := `
		{
			"encoding": "json",
			"level": "info",
			"filename": "logx.log",
			"maxsize": 128,
			"maxbackups": 30,
			"maxage": 30,
			"console": true,
			"initialFields": {
				"app": "LogX",
				"version": "1.0.0"
			}
		}
	`

	cfg := logx.ParseConfig(conf)

	//
	// cfg := logx.NewDefaultConfig()

	l := logx.NewZapLogger(cfg)
	l.With(logx.String("host", "10.8.8.8")).Info("info message")

}

```

**输出：**

```
{"level":"info","time":"2021-07-20T15:36:20.537+0800","caller":"logx/logx_test.go:30","msg":"info message","app":"LogX","version":"1.0.0","host":"10.8.8.8"}
```



**配置说明:**

```json
{
    // 格式: json|console
    "encoding": "json",
    // 等级：debug|info|warn|error|fatal
    "level": "info",
    // 日志文件路径
    "filename": "logx.log",
    // 每个日志文件保存的最大尺寸 单位：M
    "maxsize": 128,
    // 日志文件最多保存多少个备份
    "maxbackups": 30,
    // 文件最多保存多少天
    "maxage": 30,
    // 是否打印屏幕
    "console": false,
    // 默认附加字段
    "initialFields": {
        "app": "LogX",
        "version": 1.0
    }
}
```

