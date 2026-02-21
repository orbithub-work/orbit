# Satellite应用端口配置指南

## 配置文件位置

```
~/.smartarchive/config.json  (macOS/Linux)
%APPDATA%\SmartArchive\config.json  (Windows)
```

## 配置文件格式

```json
{
  "core_port": 8848,
  "satellite_base": 8850,
  "satellite_size": 50,
  "available_ports": [
    8850,
    8853,
    8854,
    8855,
    ...
  ]
}
```

**说明**：
- `core_port`: Go内核的HTTP端口
- `satellite_base`: Satellite端口段起始
- `satellite_size`: 端口段大小
- `available_ports`: 当前可用的端口列表（已排除被占用的）

---

## Satellite应用如何使用

### Go语言示例

```go
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	CorePort       int   `json:"core_port"`
	AvailablePorts []int `json:"available_ports"`
}

func main() {
	// 1. 读取配置
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	// 2. 尝试绑定端口
	port, listener := findAndBindPort(config.AvailablePorts)
	if listener == nil {
		panic("No available ports")
	}
	defer listener.Close()

	fmt.Printf("✅ Satellite started on port %d\n", port)

	// 3. 注册到Go内核
	registerToCore(config.CorePort, port)

	// 4. 启动HTTP服务
	http.Serve(listener, yourHandler)
}

func loadConfig() (*Config, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".smartarchive", "config.json")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func findAndBindPort(ports []int) (int, net.Listener) {
	for _, port := range ports {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			return port, listener
		}
	}
	return 0, nil
}

func registerToCore(corePort, myPort int) {
	url := fmt.Sprintf("http://localhost:%d/api/plugins/register", corePort)
	
	payload := map[string]interface{}{
		"plugin_id": "com.smartarchive.dock",
		"name":      "剪辑助手 Dock",
		"version":   "1.0.0",
		"mode":      "network_service",
		"endpoint":  fmt.Sprintf("http://localhost:%d", myPort),
	}

	// POST to core...
	_ = payload
	_ = url
}
```

---

### Python示例

```python
import json
import socket
from pathlib import Path
from http.server import HTTPServer, BaseHTTPRequestHandler

def load_config():
    config_path = Path.home() / '.smartarchive' / 'config.json'
    with open(config_path) as f:
        return json.load(f)

def find_available_port(ports):
    for port in ports:
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.bind(('127.0.0.1', port))
            sock.close()
            return port
        except OSError:
            continue
    return None

def main():
    # 1. 读取配置
    config = load_config()
    
    # 2. 找到可用端口
    port = find_available_port(config['available_ports'])
    if not port:
        raise Exception("No available ports")
    
    print(f"✅ Satellite started on port {port}")
    
    # 3. 注册到Go内核
    # register_to_core(config['core_port'], port)
    
    # 4. 启动HTTP服务
    server = HTTPServer(('127.0.0.1', port), YourHandler)
    server.serve_forever()

if __name__ == '__main__':
    main()
```

---

### JavaScript/Node.js示例

```javascript
const fs = require('fs');
const path = require('path');
const http = require('http');
const os = require('os');

function loadConfig() {
  const configPath = path.join(os.homedir(), '.smartarchive', 'config.json');
  return JSON.parse(fs.readFileSync(configPath, 'utf8'));
}

function findAvailablePort(ports) {
  return new Promise((resolve) => {
    const tryPort = (index) => {
      if (index >= ports.length) {
        resolve(null);
        return;
      }

      const port = ports[index];
      const server = http.createServer();
      
      server.once('error', () => {
        tryPort(index + 1);
      });
      
      server.once('listening', () => {
        server.close();
        resolve(port);
      });
      
      server.listen(port, '127.0.0.1');
    };
    
    tryPort(0);
  });
}

async function main() {
  // 1. 读取配置
  const config = loadConfig();
  
  // 2. 找到可用端口
  const port = await findAvailablePort(config.available_ports);
  if (!port) {
    throw new Error('No available ports');
  }
  
  console.log(`✅ Satellite started on port ${port}`);
  
  // 3. 注册到Go内核
  // await registerToCore(config.core_port, port);
  
  // 4. 启动HTTP服务
  const server = http.createServer(yourHandler);
  server.listen(port, '127.0.0.1');
}

main();
```

---

## 端口冲突处理

### 场景1：配置文件中的端口被占用
```
Satellite启动 → 按顺序尝试available_ports列表
→ 如果都被占用 → 报错并提示用户重启Go内核
```

### 场景2：其他软件占用了端口段
```
Go内核启动 → 扫描端口段 → 跳过被占用的端口
→ 写入available_ports（只包含可用的）
```

### 场景3：端口段完全被占用
```
Go内核启动 → 发现没有可用端口
→ 警告用户：请检查端口占用情况
→ 建议修改satellite_base配置
```

---

## 最佳实践

1. **Satellite应用启动顺序**：
   - 先读取配置
   - 按顺序尝试端口
   - 成功绑定后立即注册到Go内核

2. **错误处理**：
   - 如果所有端口都被占用，提示用户重启Go内核
   - 记录日志，方便排查

3. **端口释放**：
   - Satellite退出时，端口自动释放
   - 下次启动可以重用

4. **配置更新**：
   - 如果端口冲突严重，用户可以手动修改config.json
   - 或者重启Go内核，自动重新扫描
