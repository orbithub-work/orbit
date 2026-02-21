#!/bin/bash

# OrbitHub Git SSH 配置脚本

echo "🔑 配置 OrbitHub Git SSH 签名"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 1. 生成SSH密钥
echo ""
echo "1️⃣  生成SSH密钥..."
ssh-keygen -t ed25519 -C "zengchang42@gmail.com" -f ~/.ssh/orbithub_ed25519 -N ""

# 2. 启动ssh-agent
echo ""
echo "2️⃣  启动ssh-agent..."
eval "$(ssh-agent -s)"

# 3. 添加密钥到ssh-agent
echo ""
echo "3️⃣  添加密钥到ssh-agent..."
ssh-add ~/.ssh/orbithub_ed25519

# 4. 配置Git用户信息
echo ""
echo "4️⃣  配置Git用户信息..."
git config --global user.name "zengchangwt"
git config --global user.email "zengchang42@gmail.com"

# 5. 配置Git使用SSH签名
echo ""
echo "5️⃣  配置Git SSH签名..."
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/orbithub_ed25519.pub
git config --global commit.gpgsign true

# 6. 配置SSH
echo ""
echo "6️⃣  配置SSH..."
cat >> ~/.ssh/config << 'EOF'

# OrbitHub GitHub
Host github.com
  HostName github.com
  User git
  IdentityFile ~/.ssh/orbithub_ed25519
  AddKeysToAgent yes
EOF

# 7. 显示公钥
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 配置完成！"
echo ""
echo "📋 请复制下面的公钥，添加到GitHub："
echo ""
echo "👉 https://github.com/settings/keys"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat ~/.ssh/orbithub_ed25519.pub
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 步骤："
echo "1. 复制上面的公钥"
echo "2. 打开 https://github.com/settings/keys"
echo "3. 点击 'New SSH key'"
echo "4. Title: OrbitHub"
echo "5. Key type: Signing Key"
echo "6. 粘贴公钥"
echo "7. 点击 'Add SSH key'"
echo ""
echo "🧪 测试连接："
echo "   ssh -T git@github.com"
echo ""
