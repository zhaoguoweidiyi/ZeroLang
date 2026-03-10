import sys

def parse_instruction(line):
    """解析一条指令"""
    parts = line.strip().split()
    if len(parts) != 4:
        raise ValueError(f"指令格式错误: {line}")
    try:
        instr = int(parts[0])
        param1 = int(parts[1])
        param2 = int(parts[2])
        target = int(parts[3])
    except ValueError:
        raise ValueError(f"指令包含非数字: {line}")
    return instr, param1, param2, target

def validate_instruction(instr, param1, param2, target):
    """验证指令是否合法"""
    if not 0 <= instr <= 23:
        raise ValueError(f"指令码无效: {instr}")
    
    if instr == 1:  # 赋值
        if not 0 <= target <= 15:
            raise ValueError(f"寄存器编号无效: {target}")
    elif instr in (2, 3, 4, 5, 16, 17, 19, 20, 21, 22, 23):  # 算术运算和位运算
        if not (0 <= param1 <= 15 and 0 <= param2 <= 15 and 0 <= target <= 15):
            raise ValueError(f"寄存器编号无效")
    elif instr in (6, 7):  # 比较
        if not (0 <= param1 <= 15 and 0 <= param2 <= 15):
            raise ValueError(f"寄存器编号无效")
    elif instr == 8:  # 无条件跳转
        if param2 != 0 or target != 0:
            raise ValueError(f"跳转指令参数错误")
    elif instr == 9:  # 条件跳转
        if param2 != 0 or target != 0:
            raise ValueError(f"条件跳转指令参数错误")
    elif instr == 10:  # 输出
        if not 0 <= param1 <= 15 or param2 != 0 or target != 0:
            raise ValueError(f"输出指令参数错误")
    elif instr == 11:  # 输入
        if not 0 <= param1 <= 15 or param2 != 0 or target != 0:
            raise ValueError(f"输入指令参数错误")
    elif instr == 12:  # 读内存
        if param2 != 0 or not 0 <= target <= 15:
            raise ValueError(f"读内存指令参数错误")
    elif instr == 13:  # 写内存
        if not 0 <= param1 <= 15 or target != 0:
            raise ValueError(f"写内存指令参数错误")
    elif instr == 14:  # 函数调用
        if param2 != 0 or target != 0:
            raise ValueError(f"函数调用指令参数错误")
    elif instr == 15:  # 程序结束
        if param1 != 0 or param2 != 0 or target != 0:
            raise ValueError(f"程序结束指令参数错误")
    elif instr == 18:  # 逻辑非
        if not 0 <= param1 <= 15 or param2 != 0 or not 0 <= target <= 15:
            raise ValueError(f"逻辑非指令参数错误")

def compile_code(code):
    """编译零点语言代码"""
    instructions = []
    lines = code.split('\n')  # 减少一次strip操作
    for i, line in enumerate(lines, 1):
        line = line.strip()
        if not line:
            continue
        try:
            parts = line.split()
            if len(parts) != 4:
                raise ValueError(f"指令格式错误: {line}")
            instr = int(parts[0])
            param1 = int(parts[1])
            param2 = int(parts[2])
            target = int(parts[3])
            # 内联验证，减少函数调用
            if not 0 <= instr <= 23:
                raise ValueError(f"指令码无效: {instr}")
            if instr == 1:
                if not 0 <= target <= 15:
                    raise ValueError(f"寄存器编号无效: {target}")
            elif instr in (2, 3, 4, 5, 16, 17, 19, 20, 21, 22, 23):
                if not (0 <= param1 <= 15 and 0 <= param2 <= 15 and 0 <= target <= 15):
                    raise ValueError(f"寄存器编号无效")
            elif instr in (6, 7):
                if not (0 <= param1 <= 15 and 0 <= param2 <= 15):
                    raise ValueError(f"寄存器编号无效")
            elif instr == 8:
                if param2 != 0 or target != 0:
                    raise ValueError(f"跳转指令参数错误")
            elif instr == 9:
                if param2 != 0 or target != 0:
                    raise ValueError(f"条件跳转指令参数错误")
            elif instr == 10:
                if not 0 <= param1 <= 15 or param2 != 0 or target != 0:
                    raise ValueError(f"输出指令参数错误")
            elif instr == 11:
                if not 0 <= param1 <= 15 or param2 != 0 or target != 0:
                    raise ValueError(f"输入指令参数错误")
            elif instr == 12:
                if param2 != 0 or not 0 <= target <= 15:
                    raise ValueError(f"读内存指令参数错误")
            elif instr == 13:
                if not 0 <= param1 <= 15 or target != 0:
                    raise ValueError(f"写内存指令参数错误")
            elif instr == 14:
                if param2 != 0 or target != 0:
                    raise ValueError(f"函数调用指令参数错误")
            elif instr == 15:
                if param1 != 0 or param2 != 0 or target != 0:
                    raise ValueError(f"程序结束指令参数错误")
            elif instr == 18:
                if not 0 <= param1 <= 15 or param2 != 0 or not 0 <= target <= 15:
                    raise ValueError(f"逻辑非指令参数错误")
            instructions.append((instr, param1, param2, target))
        except ValueError as e:
            raise ValueError(f"第{i}行错误: {e}")
    return instructions

def main():
    """编译器主函数"""
    if len(sys.argv) != 2:
        print("用法: zero-compiler <input_file>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            code = f.read()
        instructions = compile_code(code)
        print(f"编译成功! 共 {len(instructions)} 条指令")
        for i, (instr, param1, param2, target) in enumerate(instructions, 1):
            print(f"指令{i}: {instr} {param1} {param2} {target}")
    except Exception as e:
        print(f"编译错误: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()