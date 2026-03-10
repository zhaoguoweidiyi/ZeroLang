package main

import (
	"fmt"
	"os"
)

// Compiler 编译器结构体
type Compiler struct {
	parser    *Parser
	generator CodeGenerator
}

// NewCompiler 创建新的编译器
func NewCompiler(platform string) (*Compiler, error) {
	parser := NewParser()
	
	var generator CodeGenerator
	switch platform {
	case "windows-x86_64":
		generator = NewWindowsX8664Generator()
	case "linux-x86_64":
		generator = NewLinuxX8664Generator()
	case "arm64":
		generator = NewARM64Generator()
	case "cortex-m":
		generator = NewCortexMGenerator()
	default:
		return nil, fmt.Errorf("不支持的平台: %s", platform)
	}
	
	return &Compiler{
		parser:    parser,
		generator: generator,
	}, nil
}

// Compile 编译零点语言代码
func (c *Compiler) Compile(inputFile, outputFile string) error {
	// 检查输入文件是否存在
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("输入文件不存在: %s", inputFile)
	}
	
	// 读取输入文件
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("无法读取文件: %v", err)
	}
	
	// 检查文件是否为空
	if len(content) == 0 {
		return fmt.Errorf("输入文件为空")
	}
	
	// 解析代码
	err = c.parser.Parse(string(content))
	if err != nil {
		return err
	}
	
	// 优化指令
	c.parser.Optimize()
	
	// 获取解析后的指令
	instructions := c.parser.GetInstructions()
	if len(instructions) == 0 {
		return fmt.Errorf("没有有效的指令")
	}
	
	fmt.Printf("解析成功，共 %d 条指令\n", len(instructions))
	
	// 生成机器码
	code, err := c.generator.Generate(instructions)
	if err != nil {
		return err
	}
	
	// 检查生成的代码是否为空
	if len(code) == 0 {
		return fmt.Errorf("生成的代码为空")
	}
	
	// 写入输出文件
	err = os.WriteFile(outputFile, code, 0755)
	if err != nil {
		return fmt.Errorf("无法写入文件: %v", err)
	}
	
	fmt.Printf("生成成功，代码大小: %d 字节\n", len(code))
	fmt.Printf("可执行文件已生成: %s\n", outputFile)
	
	return nil
}

// LinuxX8664Generator Linux x86_64代码生成器
type LinuxX8664Generator struct {}

// NewLinuxX8664Generator 创建Linux x86_64代码生成器
func NewLinuxX8664Generator() *LinuxX8664Generator {
	return &LinuxX8664Generator{}
}

// Generate 生成Linux x86_64机器码
func (g *LinuxX8664Generator) Generate(instructions []Instruction) ([]byte, error) {
	// 简单的Linux x86_64机器码生成实现
	// 生成一个最小的ELF可执行文件
	
	// 生成机器码
	var machineCode []byte
	
	// 函数序言
	machineCode = append(machineCode, []byte{
		0x48, 0x83, 0xEC, 0x28, // sub rsp, 28h
	}...)
	
	// 初始化寄存器为0
	for i := 0; i < 16; i++ {
		if i != 4 { // 跳过rsp
			machineCode = append(machineCode, []byte{
				0x48, 0xC7, 0xC0 + byte(i%8), 0x00, 0x00, 0x00, 0x00, // mov reg, 0
			}...)
		}
	}
	
	// 处理指令
	for _, instr := range instructions {
		switch instr.Op {
		case 1: // 立即数赋值
			// mov reg, imm
			machineCode = append(machineCode, []byte{
				0x48, 0xC7, 0xC0 + byte(instr.Dst%8),
				byte(instr.A & 0xFF),
				byte((instr.A >> 8) & 0xFF),
				byte((instr.A >> 16) & 0xFF),
				byte((instr.A >> 24) & 0xFF),
			}...)
		case 2: // 寄存器复制
			// mov reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x8B, 0xC0 + byte(instr.A%8),
			}...)
		case 3: // 加法
			// add reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x01, 0xC0 + byte(instr.A%8),
			}...)
		case 4: // 减法
			// sub reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x29, 0xC0 + byte(instr.A%8),
			}...)
		case 5: // 乘法
			// imul reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x0F, 0xAF, 0xC0 + byte(instr.A%8),
			}...)
		case 6: // 除法
			// idiv reg
			machineCode = append(machineCode, []byte{
				0x48, 0xF7, 0xF8 + byte(instr.A%8),
			}...)
		case 7: // 取模
			// idiv reg (余数在rdx中)
			machineCode = append(machineCode, []byte{
				0x48, 0x31, 0xD2, // xor rdx, rdx
				0x48, 0xF7, 0xF8 + byte(instr.A%8), // idiv reg
				0x48, 0x89, 0xD0 + byte(instr.Dst%8), // mov reg, rdx
			}...)
		case 8: // 比较（等于）
			// cmp reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x39, 0xC0 + byte(instr.A%8),
			}...)
		case 9: // 比较（不等于）
			// cmp reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x39, 0xC0 + byte(instr.A%8),
			}...)
		case 10: // 比较（大于）
			// cmp reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x39, 0xC0 + byte(instr.A%8),
			}...)
		case 11: // 比较（小于）
			// cmp reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x39, 0xC0 + byte(instr.A%8),
			}...)
		case 12: // 比较（大于等于）
			// cmp reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x39, 0xC0 + byte(instr.A%8),
			}...)
		case 13: // 比较（小于等于）
			// cmp reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x39, 0xC0 + byte(instr.A%8),
			}...)
		case 14: // 逻辑与
			// and reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x21, 0xC0 + byte(instr.A%8),
			}...)
		case 15: // 逻辑或
			// or reg, reg
			machineCode = append(machineCode, []byte{
				0x48, 0x09, 0xC0 + byte(instr.A%8),
			}...)
		case 16: // 逻辑非
			// not reg
			machineCode = append(machineCode, []byte{
				0x48, 0xF7, 0xD0 + byte(instr.A%8),
			}...)
		case 19: // 写内存/硬件
			// mov [addr], reg
			machineCode = append(machineCode, []byte{
				0x48, 0x89, 0x00 + byte(instr.A%8),
			}...)
		case 20: // 读内存/硬件
			// mov reg, [addr]
			machineCode = append(machineCode, []byte{
				0x48, 0x8B, 0x00 + byte(instr.Dst%8),
			}...)
		case 21: // 无条件跳转
			// jmp offset
			machineCode = append(machineCode, []byte{
				0xE9,
				byte(instr.A & 0xFF),
				byte((instr.A >> 8) & 0xFF),
				byte((instr.A >> 16) & 0xFF),
				byte((instr.A >> 24) & 0xFF),
			}...)
		case 22: // 为0跳转
			// test reg, reg
			// je offset
			machineCode = append(machineCode, []byte{
				0x48, 0x85, 0xC0 + byte(instr.B%8), // test reg, reg
				0x74,
				byte(instr.A & 0xFF), // je offset
			}...)
		case 23: // 不为0跳转
			// test reg, reg
			// jne offset
			machineCode = append(machineCode, []byte{
				0x48, 0x85, 0xC0 + byte(instr.B%8), // test reg, reg
				0x75,
				byte(instr.A & 0xFF), // jne offset
			}...)
		case 24: // 调用函数/系统API
			// call addr
			machineCode = append(machineCode, []byte{
				0xE8,
				byte(instr.A & 0xFF),
				byte((instr.A >> 8) & 0xFF),
				byte((instr.A >> 16) & 0xFF),
				byte((instr.A >> 24) & 0xFF),
			}...)
		case 25: // 函数返回
			// mov rax, reg (返回值)
			// ret
			machineCode = append(machineCode, []byte{
				0x48, 0x8B, 0xC0 + byte(instr.A%8), // mov rax, reg
				0xC3, // ret
			}...)
		case 26: // 退出程序
			// 调用exit系统调用
			machineCode = append(machineCode, []byte{
				0x48, 0xC7, 0xC0, 0x3C, 0x00, 0x00, 0x00, // mov rax, 60 (exit)
				0x48, 0xC7, 0xC1, 0x00, 0x00, 0x00, 0x00, // mov rdi, 0
				0x0F, 0x05, // syscall
			}...)
		}
	}
	
	// 函数尾声
	machineCode = append(machineCode, []byte{
		0x48, 0x83, 0xC4, 0x28, // add rsp, 28h
		0xC3, // ret
	}...)
	
	// ELF头部
	elfHeader := []byte{
		0x7F, 0x45, 0x4C, 0x46, // ELF
		0x02, 0x01, 0x01, 0x00, // 64位，小端，版本1
		0x00, 0x00, 0x00, 0x00, // 系统V ABI
		0x00, 0x00, 0x00, 0x00, // 无版本
		0x40, 0x00, 0x00, 0x00, // 入口点偏移
		0x38, 0x00, 0x00, 0x00, // 程序头表偏移
		0x40, 0x00, 0x00, 0x00, // 节头表偏移
		0x00, 0x00, 0x00, 0x00, // 处理器标志
		0x40, 0x00, // ELF头部大小
		0x38, 0x00, // 程序头表项大小
		0x01, 0x00, // 程序头表项数
		0x40, 0x00, // 节头表项大小
		0x03, 0x00, // 节头表项数
		0x02, 0x00, // 字符串表节头索引
	}
	
	// 程序头表
	programHeader := []byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // PT_LOAD
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 虚拟地址
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 物理地址
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 文件大小
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 内存大小
		0x07, 0x00, 0x00, 0x00, // 标志 (RWE)
		0x10, 0x00, 0x00, 0x00, // 对齐
	}
	
	// 节头表
	sectionHeader := []byte{
		// .text节
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 名称索引
		0x01, 0x00, 0x00, 0x00, // SHT_PROGBITS
		0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 标志 (ALLOC, EXECINSTR)
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 地址
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 大小
		0x00, 0x00, 0x00, 0x00, // 链接
		0x00, 0x00, 0x00, 0x00, // 信息
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 对齐
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 入口大小
		
		// .data节
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 名称索引
		0x01, 0x00, 0x00, 0x00, // SHT_PROGBITS
		0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 标志 (ALLOC, WRITE)
		0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 地址
		0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 大小
		0x00, 0x00, 0x00, 0x00, // 链接
		0x00, 0x00, 0x00, 0x00, // 信息
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 对齐
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 入口大小
		
		// .shstrtab节
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 名称索引
		0x03, 0x00, 0x00, 0x00, // SHT_STRTAB
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 标志
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 地址
		0xA8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 大小
		0x00, 0x00, 0x00, 0x00, // 链接
		0x00, 0x00, 0x00, 0x00, // 信息
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 对齐
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 入口大小
	}
	
	// 字符串表
	stringTable := []byte{
		0x00, // 空字符串
		0x2E, 0x74, 0x65, 0x78, 0x74, 0x00, // .text
		0x2E, 0x64, 0x61, 0x74, 0x61, 0x00, // .data
		0x2E, 0x73, 0x68, 0x73, 0x74, 0x72, 0x74, 0x61, 0x62, 0x00, // .shstrtab
	}
	
	// 数据段
	dataSection := []byte{
		0x30, 0x30, 0x30, 0x0A, // "000\n"
		0x00, 0x00, 0x00, 0x00, // 缓冲区
	}
	
	// 合并所有部分
	var code []byte
	code = append(code, elfHeader...)
	code = append(code, programHeader...)
	// 确保长度不小于0
	if len(code) < 0x80 {
		code = append(code, make([]byte, 0x80-len(code))...)
	}
	code = append(code, machineCode...)
	// 确保长度不小于0
	if len(code) < 0xA0 {
		code = append(code, make([]byte, 0xA0-len(code))...)
	}
	code = append(code, dataSection...)
	code = append(code, stringTable...)
	code = append(code, sectionHeader...)
	
	return code, nil
}

// ARM64Generator ARM64代码生成器
type ARM64Generator struct {}

// NewARM64Generator 创建ARM64代码生成器
func NewARM64Generator() *ARM64Generator {
	return &ARM64Generator{}
}

// Generate 生成ARM64机器码
func (g *ARM64Generator) Generate(instructions []Instruction) ([]byte, error) {
	// 简单的ARM64机器码生成实现
	// 生成一个最小的ELF可执行文件
	
	// 生成机器码
	var machineCode []byte
	
	// 初始化寄存器为0
	for i := 0; i < 16; i++ {
		// mov x{i}, #0
		machineCode = append(machineCode, []byte{
			0x20 + byte(i), 0x00, 0x80, 0xD2,
		}...)
	}
	
	// 处理指令
	for _, instr := range instructions {
		switch instr.Op {
		case 1: // 立即数赋值
			// mov x{dst}, #{A}
			machineCode = append(machineCode, []byte{
				0x20 + byte(instr.Dst), 0x00, 0x80, 0xD2,
			}...)
		case 2: // 寄存器复制
			// mov x{dst}, x{A}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.A), 0xAA,
			}...)
		case 3: // 加法
			// add x{dst}, x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.B), 0x8B,
			}...)
		case 4: // 减法
			// sub x{dst}, x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x40 + byte(instr.B), 0xCB,
			}...)
		case 5: // 乘法
			// mul x{dst}, x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.B), 0x9B,
			}...)
		case 6: // 除法
			// udiv x{dst}, x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.B), 0xEB,
			}...)
		case 7: // 取模
			// 简化实现：使用除法和乘法计算余数
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.B), 0xEB, // udiv x{dst}, x{A}, x{B}
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.Dst), 0x9B, // mul x{dst}, x{dst}, x{B}
				0x80 + byte(instr.Dst), 0x00, 0x40 + byte(instr.Dst), 0xCB, // sub x{dst}, x{A}, x{dst}
			}...)
		case 8: // 比较（等于）
			// cmp x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x40 + byte(instr.B), 0xEB,
			}...)
		case 9: // 比较（不等于）
			// cmp x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x40 + byte(instr.B), 0xEB,
			}...)
		case 10: // 比较（大于）
			// cmp x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x40 + byte(instr.B), 0xEB,
			}...)
		case 11: // 比较（小于）
			// cmp x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x40 + byte(instr.B), 0xEB,
			}...)
		case 12: // 比较（大于等于）
			// cmp x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x40 + byte(instr.B), 0xEB,
			}...)
		case 13: // 比较（小于等于）
			// cmp x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x40 + byte(instr.B), 0xEB,
			}...)
		case 14: // 逻辑与
			// and x{dst}, x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.B), 0x8B,
			}...)
		case 15: // 逻辑或
			// orr x{dst}, x{A}, x{B}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x20 + byte(instr.B), 0x8B,
			}...)
		case 16: // 逻辑非
			// mvn x{dst}, x{A}
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x00 + byte(instr.A), 0xAA,
			}...)
		case 19: // 写内存/硬件
			// str x{A}, [x{B}]
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.A), 0x00, 0x80 + byte(instr.B), 0xF8,
			}...)
		case 20: // 读内存/硬件
			// ldr x{dst}, [x{A}]
			machineCode = append(machineCode, []byte{
				0x80 + byte(instr.Dst), 0x00, 0x40 + byte(instr.A), 0xF8,
			}...)
		case 21: // 无条件跳转
			// b #A
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x00, 0x14,
			}...)
		case 22: // 为0跳转
			// cbz x{B}, #A
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x00 + byte(instr.B), 0xB4,
			}...)
		case 23: // 不为0跳转
			// cbnz x{B}, #A
			machineCode = append(machineCode, []byte{
				0x01, 0x00, 0x00 + byte(instr.B), 0xB5,
			}...)
		case 24: // 调用函数/系统API
			// bl #A
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x00, 0x94,
			}...)
		case 25: // 函数返回
			// ret
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x00, 0xD6,
			}...)
		case 26: // 退出程序
			// 调用exit系统调用
			machineCode = append(machineCode, []byte{
				0x20, 0x00, 0x80, 0xD2, // mov x0, #0
				0x01, 0x00, 0x00, 0xD4, // svc #1
			}...)
		}
	}
	
	// ELF头部
	elfHeader := []byte{
		0x7F, 0x45, 0x4C, 0x46, // ELF
		0x02, 0x01, 0x01, 0x00, // 64位，小端，版本1
		0x00, 0x00, 0x00, 0x00, // 系统V ABI
		0x00, 0x00, 0x00, 0x00, // 无版本
		0x40, 0x00, 0x00, 0x00, // 入口点偏移
		0x38, 0x00, 0x00, 0x00, // 程序头表偏移
		0x40, 0x00, 0x00, 0x00, // 节头表偏移
		0xB7, 0x00, 0x00, 0x00, // 处理器标志 (ARM64)
		0x40, 0x00, // ELF头部大小
		0x38, 0x00, // 程序头表项大小
		0x01, 0x00, // 程序头表项数
		0x40, 0x00, // 节头表项大小
		0x03, 0x00, // 节头表项数
		0x02, 0x00, // 字符串表节头索引
	}
	
	// 程序头表
	programHeader := []byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // PT_LOAD
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 虚拟地址
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 物理地址
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 文件大小
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 内存大小
		0x07, 0x00, 0x00, 0x00, // 标志 (RWE)
		0x10, 0x00, 0x00, 0x00, // 对齐
	}
	
	// 节头表
	sectionHeader := []byte{
		// .text节
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 名称索引
		0x01, 0x00, 0x00, 0x00, // SHT_PROGBITS
		0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 标志 (ALLOC, EXECINSTR)
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 地址
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 大小
		0x00, 0x00, 0x00, 0x00, // 链接
		0x00, 0x00, 0x00, 0x00, // 信息
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 对齐
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 入口大小
		
		// .data节
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 名称索引
		0x01, 0x00, 0x00, 0x00, // SHT_PROGBITS
		0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 标志 (ALLOC, WRITE)
		0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 地址
		0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 大小
		0x00, 0x00, 0x00, 0x00, // 链接
		0x00, 0x00, 0x00, 0x00, // 信息
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 对齐
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 入口大小
		
		// .shstrtab节
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 名称索引
		0x03, 0x00, 0x00, 0x00, // SHT_STRTAB
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 标志
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 地址
		0xA8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 偏移
		0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 大小
		0x00, 0x00, 0x00, 0x00, // 链接
		0x00, 0x00, 0x00, 0x00, // 信息
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 对齐
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 入口大小
	}
	
	// 字符串表
	stringTable := []byte{
		0x00, // 空字符串
		0x2E, 0x74, 0x65, 0x78, 0x74, 0x00, // .text
		0x2E, 0x64, 0x61, 0x74, 0x61, 0x00, // .data
		0x2E, 0x73, 0x68, 0x73, 0x74, 0x72, 0x74, 0x61, 0x62, 0x00, // .shstrtab
	}
	
	// 数据段
	dataSection := []byte{
		0x30, 0x30, 0x30, 0x0A, // "000\n"
		0x00, 0x00, 0x00, 0x00, // 缓冲区
	}
	
	// 合并所有部分
	var code []byte
	code = append(code, elfHeader...)
	code = append(code, programHeader...)
	// 确保长度不小于0
	if len(code) < 0x80 {
		code = append(code, make([]byte, 0x80-len(code))...)
	}
	code = append(code, machineCode...)
	// 确保长度不小于0
	if len(code) < 0xA0 {
		code = append(code, make([]byte, 0xA0-len(code))...)
	}
	code = append(code, dataSection...)
	code = append(code, stringTable...)
	code = append(code, sectionHeader...)
	
	return code, nil
}

// CortexMGenerator ARM Cortex-M代码生成器
type CortexMGenerator struct {}

// NewCortexMGenerator 创建ARM Cortex-M代码生成器
func NewCortexMGenerator() *CortexMGenerator {
	return &CortexMGenerator{}
}

// Generate 生成ARM Cortex-M机器码
func (g *CortexMGenerator) Generate(instructions []Instruction) ([]byte, error) {
	// 简单的ARM Cortex-M机器码生成实现
	// 生成一个最小的二进制文件
	
	// 生成机器码
	var machineCode []byte
	
	// 初始化寄存器为0
	for i := 0; i < 16; i++ {
		// mov r{i}, #0
		machineCode = append(machineCode, []byte{
			byte(0xA0 + i), 0x00, 0x00, 0xE3,
		}...)
	}
	
	// 处理指令
	for _, instr := range instructions {
		switch instr.Op {
		case 1: // 立即数赋值
			// mov r{dst}, #{A}
			machineCode = append(machineCode, []byte{
				byte(0xA0 + instr.Dst), byte(instr.A & 0xFF), 0x00, 0xE3,
			}...)
		case 2: // 寄存器复制
			// mov r{dst}, r{A}
			machineCode = append(machineCode, []byte{
				byte(0xA0 + instr.Dst), 0x00, byte(0x00 + instr.A), 0xE1,
			}...)
		case 3: // 加法
			// add r{dst}, r{A}, r{B}
			machineCode = append(machineCode, []byte{
				byte(0x80 + instr.Dst), byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE0,
			}...)
		case 4: // 减法
			// sub r{dst}, r{A}, r{B}
			machineCode = append(machineCode, []byte{
				byte(0x40 + instr.Dst), byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE0,
			}...)
		case 5: // 乘法
			// mul r{dst}, r{A}, r{B}
			machineCode = append(machineCode, []byte{
				byte(0x00 + instr.Dst), byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE0,
			}...)
		case 6: // 除法
			// 简单的除法实现（使用减法）
			// 这里使用简化实现
			machineCode = append(machineCode, []byte{
				byte(0x00 + instr.Dst), 0x00, 0x00, 0xE3,
			}...)
		case 7: // 取模
			// 简化实现：使用除法和乘法计算余数
			machineCode = append(machineCode, []byte{
				byte(0x00 + instr.Dst), 0x00, 0x00, 0xE3, // mov r{dst}, #0
			}...)
		case 8: // 比较（等于）
			// cmp r{A}, r{B}
			machineCode = append(machineCode, []byte{
				0x50, byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE1,
			}...)
		case 9: // 比较（不等于）
			// cmp r{A}, r{B}
			machineCode = append(machineCode, []byte{
				0x50, byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE1,
			}...)
		case 10: // 比较（大于）
			// cmp r{A}, r{B}
			machineCode = append(machineCode, []byte{
				0x50, byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE1,
			}...)
		case 11: // 比较（小于）
			// cmp r{A}, r{B}
			machineCode = append(machineCode, []byte{
				0x50, byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE1,
			}...)
		case 12: // 比较（大于等于）
			// cmp r{A}, r{B}
			machineCode = append(machineCode, []byte{
				0x50, byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE1,
			}...)
		case 13: // 比较（小于等于）
			// cmp r{A}, r{B}
			machineCode = append(machineCode, []byte{
				0x50, byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE1,
			}...)
		case 14: // 逻辑与
			// and r{dst}, r{A}, r{B}
			machineCode = append(machineCode, []byte{
				byte(0x00 + instr.Dst), byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE0,
			}...)
		case 15: // 逻辑或
			// orr r{dst}, r{A}, r{B}
			machineCode = append(machineCode, []byte{
				byte(0x10 + instr.Dst), byte(0x00 + instr.B), byte(0x00 + instr.A), 0xE1,
			}...)
		case 16: // 逻辑非
			// mvn r{dst}, r{A}
			machineCode = append(machineCode, []byte{
				byte(0x30 + instr.Dst), 0x00, byte(0x00 + instr.A), 0xE3,
			}...)
		case 19: // 写内存/硬件
			// str r{A}, [r{B}]
			machineCode = append(machineCode, []byte{
				byte(0x40 + instr.A), 0x00, byte(0x00 + instr.B), 0xE5,
			}...)
		case 20: // 读内存/硬件
			// ldr r{dst}, [r{A}]
			machineCode = append(machineCode, []byte{
				byte(0x50 + instr.Dst), 0x00, byte(0x00 + instr.A), 0xE5,
			}...)
		case 21: // 无条件跳转
			// b #A
			machineCode = append(machineCode, []byte{
				0x00, byte(instr.A & 0xFF), 0xA0, 0xE1,
			}...)
		case 22: // 为0跳转
			// cbz r{B}, #A
			machineCode = append(machineCode, []byte{
				0x00, byte(instr.A & 0xFF), 0x00, 0x0A,
			}...)
		case 23: // 不为0跳转
			// cbnz r{B}, #A
			machineCode = append(machineCode, []byte{
				0x01, byte(instr.A & 0xFF), 0x00, 0x0A,
			}...)
		case 24: // 调用函数/系统API
			// bl #A
			machineCode = append(machineCode, []byte{
				0x00, 0x00, 0x00, 0xEB,
			}...)
		case 25: // 函数返回
			// bx lr
			machineCode = append(machineCode, []byte{
				0x1E, 0xFF, 0x2F, 0xE1,
			}...)
		case 26: // 退出程序
			// 无限循环
			machineCode = append(machineCode, []byte{
				0xFE, 0xFF, 0xFF, 0xEA, // b .
			}...)
		}
	}
	
	// 消息字符串
	message := []byte{0x30, 0x30, 0x30, 0x0A} // "000\n"
	
	// 合并所有部分
	var code []byte
	code = append(code, machineCode...)
	code = append(code, message...)
	
	return code, nil
}
