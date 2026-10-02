# asset-refcount

资源引用计数：重复加载只算一次、归零才卸载、每帧卸载预算与 LRU，只用 Go 标准库。

```
go run ./cmd/asset-refcount --sample=load
go run ./cmd/asset-refcount --sample=release
go run ./cmd/asset-refcount --sample=budget
go run ./cmd/asset-refcount --sample=lru
go run ./cmd/asset-refcount --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **加载**：同一个资源重复 `Load` 只加载一次，`Loads` 只在第一次加一；
  每次 `Load` 让引用计数加一并返回一个新句柄。
- **释放**：`Release` 让对应资源的引用计数减一；只有计数归零时资源才进入待卸载队列；
  无效句柄返回 false 且不改变任何计数。
- **卸载预算**：`Tick` 一帧最多卸载 `budget` 个资源，返回本帧卸载数量；
  没卸完的留到下一帧。
- **卸载顺序**：待卸载资源按最后一次使用时间（`Use`）最旧的先卸，
  时间相同按名字升序；不许按释放顺序。
- **不变量**：`RefCount` 等于该资源未释放的句柄数；引用计数不会为负数；
  同一串加载、使用、释放、Tick 序列重复执行结果相同。
- **代价**：`RefCount` 不许扫描全部资源，`Scanned` 不随资源数线性增长。
- **规模**：10 万资源、每秒 10 万次查询，单次摊还 O(1)。

## 输出契约（不改格式）

```
loads=..
loaded=..
unloaded=%d
loaded_a=.. loaded_c=..
scanned=..
```

## 目录

```
assets.go              加载、引用计数、卸载预算与 LRU
assets_test.go         单元测试
cmd/asset-refcount/    命令行入口
```
