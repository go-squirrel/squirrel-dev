#!/bin/bash

echo "开始执行测试脚本"
for i in {1..10}
do
    echo "hello $i" >> /tmp/test.out
    sleep 5
done
echo "测试脚本执行完成"
