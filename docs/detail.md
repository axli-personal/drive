## 路径存储

### 方案一: 深度 + 完整路径

```text
0/
1/home
2/home/dir
3/home/dir/file
```

优点: 便于获取路径.

> [Baidu File System Design](https://github.com/baidu/bfs/blob/master/docs/cn/design.md#namespace)

### 方案二: 文件名 + 父级指针

优点: 便于操作目录.
