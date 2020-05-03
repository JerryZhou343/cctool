# cctool
字幕生成工具


# 设计
1. 调用百度api 翻译英文到中文

## 使用

./cctool translate -h
翻译字幕

Usage:
   translate [flags]

Flags:
  -f, --from string        源语言 (default "en")
  -h, --help               help for translate
  -m, --merge              双语字幕
  -s, --source strings     源文件
  -t, --to string          目标语言 (default "zh")
      --transtool string   翻译器 (default "baidu")


###  使用示例
```json
cctool translate -f en -t zh -m -s e2Engish.srt
```