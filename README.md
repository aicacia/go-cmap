# concurrent map

wraps `sync.Map` adds types and count

```go
m := cmap.New[string, string]()

isNew := m.Set("hello", "world")
value, ok := m.Get("hello")
```
