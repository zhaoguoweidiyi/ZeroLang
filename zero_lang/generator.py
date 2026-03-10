import sys

def generate_output(value):
    """生成输出指定值的代码"""
    return f"1 {value} 0 1\n10 1 0 0\n15 0 0 0"

def generate_add(a, b):
    """生成计算a+b并输出结果的代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n2 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_sub(a, b):
    """生成计算a-b并输出结果的代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n3 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_mul(a, b):
    """生成计算a*b并输出结果的代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n4 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_div(a, b):
    """生成计算a/b并输出结果的代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n5 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_input_output():
    """生成输入一个值并输出的代码"""
    return "11 1 0 0\n10 1 0 0\n15 0 0 0"

def generate_logic_and(a, b):
    """生成逻辑与运算代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n16 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_logic_or(a, b):
    """生成逻辑或运算代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n17 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_logic_not(a):
    """生成逻辑非运算代码"""
    return f"1 {a} 0 1\n18 1 0 2\n10 2 0 0\n15 0 0 0"

def generate_bit_and(a, b):
    """生成位与运算代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n19 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_bit_or(a, b):
    """生成位或运算代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n20 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_bit_xor(a, b):
    """生成位异或运算代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n21 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_left_shift(a, b):
    """生成左移运算代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n22 1 2 3\n10 3 0 0\n15 0 0 0"

def generate_right_shift(a, b):
    """生成右移运算代码"""
    return f"1 {a} 0 1\n1 {b} 0 2\n23 1 2 3\n10 3 0 0\n15 0 0 0"

def main():
    """代码生成器主函数"""
    if len(sys.argv) < 2:
        print("用法:")
        print("  zero-generator output <value>    # 输出指定值")
        print("  zero-generator add <a> <b>       # 计算a+b并输出")
        print("  zero-generator sub <a> <b>       # 计算a-b并输出")
        print("  zero-generator mul <a> <b>       # 计算a*b并输出")
        print("  zero-generator div <a> <b>       # 计算a/b并输出")
        print("  zero-generator input-output      # 输入并输出")
        print("  zero-generator logic-and <a> <b> # 计算逻辑与并输出")
        print("  zero-generator logic-or <a> <b>  # 计算逻辑或并输出")
        print("  zero-generator logic-not <a>     # 计算逻辑非并输出")
        print("  zero-generator bit-and <a> <b>   # 计算位与并输出")
        print("  zero-generator bit-or <a> <b>    # 计算位或并输出")
        print("  zero-generator bit-xor <a> <b>   # 计算位异或并输出")
        print("  zero-generator left-shift <a> <b> # 计算左移并输出")
        print("  zero-generator right-shift <a> <b> # 计算右移并输出")
        sys.exit(1)
    
    command = sys.argv[1]
    
    if command == "output":
        if len(sys.argv) != 3:
            print("用法: zero-generator output <value>")
            sys.exit(1)
        value = int(sys.argv[2])
        code = generate_output(value)
    elif command == "add":
        if len(sys.argv) != 4:
            print("用法: zero-generator add <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_add(a, b)
    elif command == "sub":
        if len(sys.argv) != 4:
            print("用法: zero-generator sub <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_sub(a, b)
    elif command == "mul":
        if len(sys.argv) != 4:
            print("用法: zero-generator mul <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_mul(a, b)
    elif command == "div":
        if len(sys.argv) != 4:
            print("用法: zero-generator div <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_div(a, b)
    elif command == "input-output":
        code = generate_input_output()
    elif command == "logic-and":
        if len(sys.argv) != 4:
            print("用法: zero-generator logic-and <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_logic_and(a, b)
    elif command == "logic-or":
        if len(sys.argv) != 4:
            print("用法: zero-generator logic-or <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_logic_or(a, b)
    elif command == "logic-not":
        if len(sys.argv) != 3:
            print("用法: zero-generator logic-not <a>")
            sys.exit(1)
        a = int(sys.argv[2])
        code = generate_logic_not(a)
    elif command == "bit-and":
        if len(sys.argv) != 4:
            print("用法: zero-generator bit-and <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_bit_and(a, b)
    elif command == "bit-or":
        if len(sys.argv) != 4:
            print("用法: zero-generator bit-or <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_bit_or(a, b)
    elif command == "bit-xor":
        if len(sys.argv) != 4:
            print("用法: zero-generator bit-xor <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_bit_xor(a, b)
    elif command == "left-shift":
        if len(sys.argv) != 4:
            print("用法: zero-generator left-shift <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_left_shift(a, b)
    elif command == "right-shift":
        if len(sys.argv) != 4:
            print("用法: zero-generator right-shift <a> <b>")
            sys.exit(1)
        a = int(sys.argv[2])
        b = int(sys.argv[3])
        code = generate_right_shift(a, b)
    else:
        print(f"未知命令: {command}")
        sys.exit(1)
    
    print("生成的零点语言代码:")
    print(code)

if __name__ == "__main__":
    main()