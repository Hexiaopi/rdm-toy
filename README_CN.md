
# rdm-toy

[![Go Report Card](https://goreportcard.com/badge/github.com/hexiaopi/rdm-toy)](https://goreportcard.com/report/github.com/hexiaopi/rdm-toy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hexiaopi/rdm-toy)](https://github.com/Hexiaopi/rdm-toy)

redis 可视化管理工具

## 功能

- redis 管理
- db 管理
- key 管理
- 支持命令行

## 🛠 技术栈

### 前端

- vue3
- vite
- pina
- element-plus

### 后端

- go
- redis
- gin

## 安装

### 前端

```bash
  cd web
  npm install
```

### 后端

```bash
  go mod tidy
```
    
## 本地运行

### 前端

```bash
  cd web
  npm run dev
```

### 后端

```bash
  go run cmd/root.go
```

## 截图

![client](./docs/images/client.png)
![conn-base](./docs/images/conn-base.png)
![conn-config](./docs/images/conn-config.png)
![conn-terminal](./docs/images/conn-terminal.png)
![conn-slowlog](./docs/images/conn-slowlog.png)
![conn-clients](./docs/images/conn-clients.png)
![conn-echart](./docs/images/conn-echart.png)
![db](./docs/images/db-1.png)
![key](./docs/images/key-1.png)

## 收藏历史

[![Star History Chart](https://api.star-history.com/svg?repos=hexiaopi/rdm-toy&type=Date)](https://star-history.com/#hexiaopi/rdm-toy&Date)


## 许可证

[MIT](https://choosealicense.com/licenses/mit/)

