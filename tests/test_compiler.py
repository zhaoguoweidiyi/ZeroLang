import unittest
from zero_lang.compiler import parse_instruction, validate_instruction, compile_code

class TestCompiler(unittest.TestCase):
    def test_parse_instruction(self):
        """测试解析指令"""
        instr, param1, param2, target = parse_instruction("1 100 0 1")
        self.assertEqual(instr, 1)
        self.assertEqual(param1, 100)
        self.assertEqual(param2, 0)
        self.assertEqual(target, 1)
    
    def test_validate_instruction(self):
        """测试验证指令"""
        # 测试合法指令
        validate_instruction(1, 100, 0, 1)
        validate_instruction(2, 1, 2, 3)
        validate_instruction(10, 1, 0, 0)
        validate_instruction(15, 0, 0, 0)
        
        # 测试非法指令
        with self.assertRaises(ValueError):
            validate_instruction(16, 100, 0, 1)  # 指令码无效
        with self.assertRaises(ValueError):
            validate_instruction(1, 100, 0, 20)  # 寄存器编号无效
        with self.assertRaises(ValueError):
            validate_instruction(10, 20, 0, 0)  # 输出指令参数错误
    
    def test_compile_code(self):
        """测试编译代码"""
        code = """
        1 100 0 1
        10 1 0 0
        15 0 0 0
        """
        instructions = compile_code(code)
        self.assertEqual(len(instructions), 3)
        self.assertEqual(instructions[0], (1, 100, 0, 1))
        self.assertEqual(instructions[1], (10, 1, 0, 0))
        self.assertEqual(instructions[2], (15, 0, 0, 0))
    
    def test_compile_invalid_code(self):
        """测试编译无效代码"""
        # 格式错误
        invalid_code = "1 100 0"
        with self.assertRaises(ValueError):
            compile_code(invalid_code)
        
        # 语义错误
        invalid_code = "1 100 0 20"
        with self.assertRaises(ValueError):
            compile_code(invalid_code)

if __name__ == "__main__":
    unittest.main()