import unittest
from zero_lang.generator import generate_output, generate_add, generate_sub, generate_mul, generate_div, generate_input_output

class TestGenerator(unittest.TestCase):
    def test_generate_output(self):
        """测试生成输出代码"""
        code = generate_output(100)
        expected = "1 100 0 1\n10 1 0 0\n15 0 0 0"
        self.assertEqual(code, expected)
    
    def test_generate_add(self):
        """测试生成加法代码"""
        code = generate_add(5, 3)
        expected = "1 5 0 1\n1 3 0 2\n2 1 2 3\n10 3 0 0\n15 0 0 0"
        self.assertEqual(code, expected)
    
    def test_generate_sub(self):
        """测试生成减法代码"""
        code = generate_sub(10, 4)
        expected = "1 10 0 1\n1 4 0 2\n3 1 2 3\n10 3 0 0\n15 0 0 0"
        self.assertEqual(code, expected)
    
    def test_generate_mul(self):
        """测试生成乘法代码"""
        code = generate_mul(6, 7)
        expected = "1 6 0 1\n1 7 0 2\n4 1 2 3\n10 3 0 0\n15 0 0 0"
        self.assertEqual(code, expected)
    
    def test_generate_div(self):
        """测试生成除法代码"""
        code = generate_div(12, 3)
        expected = "1 12 0 1\n1 3 0 2\n5 1 2 3\n10 3 0 0\n15 0 0 0"
        self.assertEqual(code, expected)
    
    def test_generate_input_output(self):
        """测试生成输入输出代码"""
        code = generate_input_output()
        expected = "11 1 0 0\n10 1 0 0\n15 0 0 0"
        self.assertEqual(code, expected)

if __name__ == "__main__":
    unittest.main()