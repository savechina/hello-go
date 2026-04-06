# 阶段复习：基础部分（Review Basic）

这章是基础部分的阶段复习（Review Basic）。目标不是再引入新语法，而是把泛型（Generics）、包（Packages）、指针（Pointers）和日志记录（Logging）串成一个完整的小程序思路。很多人单独看每章都懂，但一到组合使用就会卡住；所以复习的价值，就是把这些知识放回同一个上下文里，看看它们如何协作。

先看泛型容器。`notebook[T comparable]` 用 `comparable` 限制元素可比较，确保可以去重。

```go
type notebook[T comparable] struct {
	items []T
}

func (n *notebook[T]) Add(item T) {
	if !n.Contains(item) {
		n.items = append(n.items, item)
	}
}
```

指针接收者（pointer receiver）则负责修改状态。`Finish` 会累加已完成数量，所以它必须操作同一个对象。

```go
type learner struct {
	Name      string
	Completed int
}

func (l *learner) Finish(topic string) string {
	if l == nil {
		return "nil learner"
	}
	l.Completed++
	return fmt.Sprintf("%s finished %s", l.Name, topic)
}
```

日志部分展示了 `slog` 的结构化输出。复习程序不仅要“算出结果”，还要“说清楚结果”。

```go
func buildStudyLog(name string, completed int) string {
	var buffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buffer, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("study summary", "learner", name, "completed", completed)
	return strings.TrimSpace(buffer.String())
}
```

最后，包组织把复习章节的辅助逻辑拆出去，例如 `tag.Prefix`，让主流程更清晰。

```go
func summarizeLearnerProgress(name string, topics []string) string {
	student := &learner{Name: name}
	updates := make([]string, 0, len(topics))
	for _, topic := range topics {
		updates = append(updates, student.Finish(topic))
	}
	return fmt.Sprintf("%s %s", tag.Prefix("progress"), strings.Join(updates, " | "))
}
```

`summaryExample` 把泛型、指针和日志整合到一起，正是阶段复习要达到的效果。[源码](../../internal/basic/review/review.go)

## 复习题

1. 问：为什么 `notebook` 要用 `comparable`？
   答：因为它要判断元素是否已存在。
2. 问：为什么 `Finish` 不能随便改成值接收者？
   答：因为它要更新 `Completed`，值接收者只会改副本。
3. 问：复习章节最大的意义是什么？
   答：把分散知识点合并成可运行、可观察、可维护的小程序。
