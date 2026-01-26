#!/bin/bash

# 批量更新所有 Lambda 函数
# 使用方法: ./update-all-functions.sh

set -e

echo "=========================================="
echo "批量更新 XZ Wallet Lambda 函数"
echo "=========================================="
echo ""

# 确保已经构建
if [ ! -d ".aws-sam/build" ]; then
    echo "❌ 错误: 请先运行 'sam build'"
    exit 1
fi

# 函数列表
FUNCTIONS=(
    "GetBalanceFunction:xz-wallet-backend-GetBalanceFunction-1Aq0Hy0Hy0Hy"
    "CreateTaskFunction:xz-wallet-backend-CreateTaskFunction-1Aq0Hy0Hy0Hy"
    "ListTasksFunction:xz-wallet-backend-ListTasksFunction-1Aq0Hy0Hy0Hy"
    "GetTaskFunction:xz-wallet-backend-GetTaskFunction-1AGdjuz2dekr"
    "BidTaskFunction:xz-wallet-backend-BidTaskFunction-1Aq0Hy0Hy0Hy"
    "SelectBidderFunction:xz-wallet-backend-SelectBidderFunction-1Aq0Hy0Hy0Hy"
    "ApproveWorkFunction:xz-wallet-backend-ApproveWorkFunction-1Aq0Hy0Hy0Hy"
    "CancelTaskFunction:xz-wallet-backend-CancelTaskFunction-OPpwD7mkgToS"
    "SubmitWorkFunction:xz-wallet-backend-SubmitWorkFunction-XXXXXXXXXX"
)

# 获取实际的函数名称
echo "🔍 获取实际的 Lambda 函数名称..."
ACTUAL_FUNCTIONS=$(aws lambda list-functions --query "Functions[?starts_with(FunctionName, 'xz-wallet-backend-')].FunctionName" --output text)

# 更新每个函数
for FUNC_NAME in $ACTUAL_FUNCTIONS; do
    # 提取函数类型
    if [[ $FUNC_NAME == *"GetBalance"* ]]; then
        BUILD_DIR="GetBalanceFunction"
    elif [[ $FUNC_NAME == *"CreateTask"* ]]; then
        BUILD_DIR="CreateTaskFunction"
    elif [[ $FUNC_NAME == *"ListTasks"* ]]; then
        BUILD_DIR="ListTasksFunction"
    elif [[ $FUNC_NAME == *"GetTask"* ]]; then
        BUILD_DIR="GetTaskFunction"
    elif [[ $FUNC_NAME == *"BidTask"* ]]; then
        BUILD_DIR="BidTaskFunction"
    elif [[ $FUNC_NAME == *"SelectBidder"* ]]; then
        BUILD_DIR="SelectBidderFunction"
    elif [[ $FUNC_NAME == *"ApproveWork"* ]]; then
        BUILD_DIR="ApproveWorkFunction"
    elif [[ $FUNC_NAME == *"CancelTask"* ]]; then
        BUILD_DIR="CancelTaskFunction"
    elif [[ $FUNC_NAME == *"SubmitWork"* ]]; then
        BUILD_DIR="SubmitWorkFunction"
    else
        echo "⚠️  跳过未知函数: $FUNC_NAME"
        continue
    fi

    echo ""
    echo "📦 更新 $BUILD_DIR -> $FUNC_NAME"
    
    # 创建 zip 文件
    cd .aws-sam/build/$BUILD_DIR
    if [ -f "$BUILD_DIR.zip" ]; then
        rm "$BUILD_DIR.zip"
    fi
    zip -q "$BUILD_DIR.zip" bootstrap
    cd ../../..
    
    # 更新函数
    aws lambda update-function-code \
        --function-name "$FUNC_NAME" \
        --zip-file "fileb://.aws-sam/build/$BUILD_DIR/$BUILD_DIR.zip" \
        --output text \
        --query 'FunctionName' > /dev/null
    
    echo "✅ $BUILD_DIR 更新成功"
    
    # 等待一下避免 API 限流
    sleep 1
done

echo ""
echo "=========================================="
echo "✅ 所有函数更新完成！"
echo "=========================================="
echo ""
echo "等待函数激活..."
sleep 5

# 检查所有函数状态
echo ""
echo "📊 函数状态检查:"
for FUNC_NAME in $ACTUAL_FUNCTIONS; do
    STATUS=$(aws lambda get-function --function-name "$FUNC_NAME" --query 'Configuration.LastUpdateStatus' --output text)
    if [ "$STATUS" == "Successful" ]; then
        echo "  ✅ $FUNC_NAME: $STATUS"
    else
        echo "  ⏳ $FUNC_NAME: $STATUS"
    fi
done

echo ""
echo "🎉 完成！所有函数已更新。"
