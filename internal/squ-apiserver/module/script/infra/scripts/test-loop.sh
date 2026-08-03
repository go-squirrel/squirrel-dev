#!/bin/bash

echo "开始执行测试脚本"
echo "hello"
echo "hello" > /tmp/test.out
for i in {1..10}
do
    echo "world $i" >> /tmp/test.out
    echo "world $i"
    sleep 5
done
echo "测试脚本执行完成"
