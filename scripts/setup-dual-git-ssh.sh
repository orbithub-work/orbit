#!/bin/bash

# 配置 GitHub 和 Gitee 双SSH密钥

echo "🔑 配置 GitHub 和 Gitee 双SSH密钥"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 1. 生成GitHub密钥（如果不存在）
if [ ! -f ~/.ssh/github_ed25519 ]; then
    echo ""
    echo "1️⃣  生成GitHub SSH密钥..."
    ssh-keygen -t ed25519 -C "zengchang42@gmail.com" -f ~/.ssh/github_ed25519 -N ""
else
    echo ""
    echo "1️⃣  GitHub密钥已存在，跳过生成"
fi

# 2. 生成Gitee密钥（如果不存在）
if [ ! -f ~/.ssh/gitee_ed25519 ]; then
    echo ""
    echo "2️⃣  生成Gitee SSH密钥..."
    ssh-keygen -t ed25519 -C "zengchang42@gmail.com" -f ~/.ssh/gitee_ed25519 -N ""
else
    echo ""
    echo "2️⃣  Gitee密钥已存在，跳过生成"
fi

# 3. 启动ssh-agent
echo ""
echo "3️⃣  启动ssh-agent..."
eval "$(ssh-agent -s)"

# 4. 添加密钥到ssh-agent
echo ""
echo "4️⃣  添加密钥到ssh-agent..."
ssh-add ~/.ssh/github_ed25519
ssh-add ~/.ssh/gitee_ed25519

# 5. 配置SSH config
echo ""
echo "5️⃣  配置SSH..."

# 备份现有配置
if [ -f ~/.ssh/config ]; then
    cp ~/.ssh/config ~/.ssh/config.backup.$(date +%Y%m%d_%H%M%S)
fi

# 写入新配置
cat > ~/.ssh/config << 'EOF'
# GitHub
Host github.com
  HostName github.com
  User git
  IdentityFile ~/.ssh/github_ed25519
  AddKeysToAgent yes

# Gitee
Host gitee.com
  HostName gitee.com
  User git
  IdentityFile ~/.ssh/gitee_ed25519
  AddKeysToAgent yes
EOF

# 6. 配置Git（全局默认）
echo ""
echo "6️⃣  配置Git全局信息..."
git config --global user.name "zengchangwt"
git config --global user.email "zengchang42@gmail.com"

# 7. 配置Git签名（使用GitHub密钥作为默认）
echo ""
echo "7️⃣  配置Git签名（默认使用GitHub密钥）..."
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/github_ed25519.pub
git config --global commit.gpgsign true

# 8. 显示公钥
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 配置完成！"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 GitHub 公钥（添加到 https://github.com/settings/keys）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat ~/.ssh/github_ed25519.pub
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Gitee 公钥（添加到 https://gitee.com/profile/sshkeys）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat ~/.ssh/gitee_ed25519.pub
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 添加步骤："
echo ""
echo "【GitHub】"
echo "1. 打开 https://github.com/settings/keys"
echo "2. 点击 'New SSH key'"
echo "3. Title: OrbitHub"
echo "4. Key type: Authentication Key"
echo "5. 粘贴上面的 GitHub 公钥"
echo "6. 再添加一次，Key type 选 Signing Key"
echo ""
echo "【Gitee】"
echo "1. 打开 https://gitee.com/profile/sshkeys"
echo "2. 点击 '添加公钥'"
echo "3. 标题: OrbitHub"
echo "4. 粘贴上面的 Gitee 公钥"
echo ""
echo "🧪 测试连接："
echo "   ssh -T git@github.com"
echo "   ssh -T git@gitee.com"
echo ""
echo "💡 提示："
echo "   - GitHub 和 Gitee 会自动使用各自的密钥"
echo "   - 提交签名默认使用 GitHub 密钥"
echo "   - 如需为特定仓库使用 Gitee 密钥签名："
echo "     cd your-gitee-repo"
echo "     git config user.signingkey ~/.ssh/gitee_ed25519.pub"
echo ""
