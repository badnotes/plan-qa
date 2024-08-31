#!/bin/bash

# 进程名
process_name="plan-qa"

go build

# 查找进程ID
pid=$(ps aux | grep "$process_name" | grep -v grep | awk '{print $2}')

# 检查是否找到进程
if [ -z "$pid" ]; then
  echo "没有找到进程 $process_name"
else
  # 终止进程
  echo "找到进程 $process_name，PID 为 $pid"
  kill $pid
  echo "已终止进程 $process_name"
fi

./$process_name &

pid=$!

echo "进程 $process_name 已启动，PID 为 $pid"

