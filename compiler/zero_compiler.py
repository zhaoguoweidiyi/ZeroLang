#!/usr/bin/env python3
"""
零点语言原生编译器
输入：纯数字4字段指令（op a b dst）
输出：各平台原生机器码
"""

import argparse
import struct
import os

class ZeroCompiler:
    """零点语言编译器"""
    
    def __init__(self, target_arch="x86_64", target_os="windows"):
        """初始化编译器"""
        self.target_arch = target_arch
        self.target_os = target_os
        self.instructions = []
        self.labels = {}
        self.label_counter = 0
    
    def parse_file(self, input_file):
        """解析输入文件"""
        with open(input_file, 'r', encoding='utf-8') as f:
            lines = f.readlines()
        
        for i, line in enumerate(lines, 1):
            line = line.strip()
            if not line:
                continue
            
            parts = line.split()
            if len(parts) != 4:
                raise ValueError(f"第{i}行错误: 指令格式错误: {line}")
            
            try:
                op = int(parts[0])
                a = int(parts[1])
                b = int(parts[2])
                dst = int(parts[3])
            except ValueError:
                raise ValueError(f"第{i}行错误: 指令包含非数字: {line}")
            
            # 验证指令码
            if op < 0 or op > 15:
                raise ValueError(f"第{i}行错误: 指令码无效: {op}")
            
            # 验证寄存器编号
            if op == 1:  # 赋值
                if dst < 0 or dst > 15:
                    raise ValueError(f"第{i}行错误: 寄存器编号无效: {dst}")
            elif op in (2, 3, 4, 5):  # 算术运算
                if a < 0 or a > 15 or b < 0 or b > 15 or dst < 0 or dst > 15:
                    raise ValueError(f"第{i}行错误: 寄存器编号无效")
            elif op in (6, 7):  # 比较