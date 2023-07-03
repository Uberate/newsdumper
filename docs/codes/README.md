# 编码

newsdumper 提供了一些简单的封装将杂乱无章编码流程化。辅助任何人快速扩充本项目（Fork）。

NewsDumper 提供的 `staged`如下：

[internal/logics/init.go](../../internal/logics/init.go)

提供了多个阶段分别处理数据：

- [StartStaged](staged/StartStaged.md)：初始化阶段
- [GetterStaged](staged/GetterStaged.md)：获取新闻
- [SplitStaged](staged/SplitStaged.md)：裁剪新闻
- [SendStaged](staged/SendStaged.md)：调用回调