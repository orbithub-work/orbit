# 内网插件 API 规范 (v1)

## 概述
智归档OS 提供一套内网 HTTP API，允许外部插件（如剪映助手）与基座进行交互。

- **基础 URL**: `http://localhost:32000/api`
- **默认端口**: `32000` (若被占用会自动递增)

## 系统接口

### 健康检查
`GET /health`

**响应示例**:
```json
{
    "status": "ok"
}
```

### 获取系统信息
`GET /system/info`

**响应示例**:
```json
{
    "os": "windows",
    "arch": "amd64",
    "pid": 1234,
    "port": 32000,
    "start_time": 1739260800
}
```

## 资产接口

### 索引文件
`POST /assets/index`

**请求参数**:
```json
{
    "path": "C:\\path\\to\\file.mp4",
    "project_id": "proj_1"
}
```

**响应示例**:
```json
{
    "success": true,
    "data": {
        "id": "asset_1",
        "name": "file.mp4"
    }
}
```

## 项目接口

### 获取项目列表
`GET /projects`

**响应示例**:
```json
{
    "success": true,
    "data": [
        {
            "id": "default",
            "name": "默认项目"
        }
    ]
}
```

## 资产接口

### 搜索资产
`GET /assets/search?q={query}`

**参数**:
- `q`: 搜索关键词

**响应示例**:
```json
{
    "query": "test",
    "results": [],
    "message": "Search implementation pending integration with SearchService"
}
```
