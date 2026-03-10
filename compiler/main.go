package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Instruction 零点语言指令
type Instruction struct {
	Op   int
	A    int
	B    int
	Dst  int
}

// Parser 解析器
type Parser struct {
	instructions []Instruction
}

// NewParser 创建新的解析器
func NewParser() *Parser {
	return &Parser{
		instructions: []Instruction{},
	}
}

// Parse 解析零点语言代码
func (p *Parser) Parse(code string) error {
	lines := strings.Split(code, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 4 {
			return fmt.Errorf("第%d行错误: 指令格式错误: %s", i+1, line)
		}

		var instr Instruction
		for j, part := range parts {
			val, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("第%d行错误: 指令包含非数字: %s", i+1, line)
			}
			switch j {
			case 0:
				instr.Op = val
			case 1:
				instr.A = val
			case 2:
				instr.B = val
			case 3:
				instr.Dst = val
			}
		}

		// 验证指令码
		if instr.Op < 1 || instr.Op > 26 || (instr.Op > 16 && instr.Op < 19) {
			return fmt.Errorf("第%d行错误: 指令码无效: %d", i+1, instr.Op)
		}

		// 验证寄存器编号
		switch instr.Op {
		case 1: // 立即数赋值
			if instr.Dst < 0 || instr.Dst > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效: %d", i+1, instr.Dst)
			}
		case 2: // 寄存器复制
			if instr.A < 0 || instr.A > 15 || instr.Dst < 0 || instr.Dst > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效", i+1)
			}
		case 3, 4, 5, 6, 7: // 算术运算
			if instr.A < 0 || instr.A > 15 || instr.B < 0 || instr.B > 15 || instr.Dst < 0 || instr.Dst > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效", i+1)
			}
		case 8, 9, 10, 11, 12, 13: // 比较运算
			if instr.A < 0 || instr.A > 15 || instr.B < 0 || instr.B > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效", i+1)
			}
		case 14, 15: // 逻辑与、逻辑或
			if instr.A < 0 || instr.A > 15 || instr.B < 0 || instr.B > 15 || instr.Dst < 0 || instr.Dst > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效", i+1)
			}
		case 16: // 逻辑非
			if instr.A < 0 || instr.A > 15 || instr.Dst < 0 || instr.Dst > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效", i+1)
			}
		case 19: // 写内存/硬件
			if instr.A < 0 || instr.A > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效", i+1)
			}
		case 20: // 读内存/硬件
			if instr.Dst < 0 || instr.Dst > 15 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效", i+1)
			}
		case 21: // 无条件跳转
			if instr.Dst != 0 {
				return fmt.Errorf("第%d行错误: 跳转指令目标寄存器必须为0", i+1)
			}
		case 22, 23: // 为0跳转、不为0跳转
			if instr.B < 0 || instr.B > 15 || instr.Dst != 0 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效或目标寄存器不为0", i+1)
			}
		case 24: // 调用函数/系统API
			if instr.Dst != 0 {
				return fmt.Errorf("第%d行错误: 调用指令目标寄存器必须为0", i+1)
			}
		case 25: // 函数返回
			if instr.A < 0 || instr.A > 15 || instr.Dst != 0 {
				return fmt.Errorf("第%d行错误: 寄存器编号无效或目标寄存器不为0", i+1)
			}
		case 26: // 退出程序
			if instr.A != 0 || instr.B != 0 || instr.Dst != 0 {
				return fmt.Errorf("第%d行错误: 退出指令参数必须为0", i+1)
			}
		}

		p.instructions = append(p.instructions, instr)
	}

	return nil
}

// GetInstructions 获取解析后的指令
func (p *Parser) GetInstructions() []Instruction {
	return p.instructions
}

// Optimize 优化指令
func (p *Parser) Optimize() {
	var optimized []Instruction
	
	for i := 0; i < len(p.instructions); i++ {
		instr := p.instructions[i]
		
		// 优化1: 移除连续的空操作
		if instr.Op == 0 && i > 0 && p.instructions[i-1].Op == 0 {
			continue
		}
		
		// 优化2: 移除不必要的加载立即数（如果下一条指令也是加载相同的寄存器）
		if instr.Op == 1 && i < len(p.instructions)-1 {
			nextInstr := p.instructions[i+1]
			if nextInstr.Op == 1 && nextInstr.Dst == instr.Dst {
				continue
			}
		}
		
		optimized = append(optimized, instr)
	}
	
	p.instructions = optimized
}

// CodeGenerator 代码生成器
type CodeGenerator interface {
	Generate(instructions []Instruction) ([]byte, error)
}

// WindowsX8664Generator Windows x86_64代码生成器
type WindowsX8664Generator struct {}

// NewWindowsX8664Generator 创建Windows x86_64代码生成器
func NewWindowsX8664Generator() *WindowsX8664Generator {
	return &WindowsX8664Generator{}
}

// Generate 生成Windows x86_64机器码
func (g *WindowsX8664Generator) Generate(instructions []Instruction) ([]byte, error) {
	// 简单的x86_64机器码生成实现
	// 生成一个最小的Windows可执行文件
	// 注意：这是一个简化的实现，实际的PE文件格式更加复杂
	

	
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
			// 调用ExitProcess
			machineCode = append(machineCode, []byte{
				0x48, 0xC7, 0xC0, 0x00, 0x00, 0x00, 0x00, // mov rax, 0
				0x48, 0xB8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // mov rax, ExitProcess
				0xFF, 0xD0, // call rax
			}...)
		}
	}
	
	// 函数尾声
	machineCode = append(machineCode, []byte{
		0x48, 0x83, 0xC4, 0x28, // add rsp, 28h
		0xC3, // ret
	}...)
	
	// MZ头部
	mzHeader := []byte{
		0x4D, 0x5A, // MZ
		0x90, 0x00, // 重定位表偏移
		0x03, 0x00, // 文件块数
		0x00, 0x00, // 重定位项数
		0x04, 0x00, // 头部大小（以段落为单位）
		0x00, 0x00, // 最小额外段落
		0xFF, 0xFF, // 最大额外段落
		0x00, 0x00, // 初始SS值
		0x00, 0x00, // 初始SP值
		0x00, 0x00, // 校验和
		0x00, 0x00, // 初始IP值
		0x00, 0x00, // 初始CS值
		0x40, 0x00, // 重定位表偏移
		0x00, 0x00, // 覆盖号
		0x00, 0x00, 0x00, 0x00, // 保留
		0x00, 0x00, 0x00, 0x00, // 保留
		0x00, 0x00, 0x00, 0x00, // 保留
		0x00, 0x00, 0x00, 0x00, // 保留
		0x00, 0x00, 0x00, 0x00, // 保留
		0x00, 0x00, 0x00, 0x00, // 保留
		0x80, 0x00, 0x00, 0x00, // 扩展头部偏移
	}
	
	// PE头部
	peHeader := []byte{
		0x50, 0x45, 0x00, 0x00, // PE
		0x64, 0x86, 0x00, 0x00, // 机器类型 (x86_64)
		0x01, 0x00, // 节数
		0x00, 0x00, 0x00, 0x00, // 时间戳
		0x00, 0x00, 0x00, 0x00, // 符号表偏移
		0x00, 0x00, 0x00, 0x00, // 符号数
		0xE0, 0x00, 0x00, 0x00, // 可选头部大小
		0x02, 0x00, 0x00, 0x00, // 特性
	}
	
	// 可选头部
	optionalHeader := []byte{
		0x0B, 0x02, // 魔术字 (PE32+)
		0x00, 0x03, // 主版本号
		0x00, 0x00, // 次版本号
		0x00, 0x00, 0x00, 0x00, // 代码大小
		0x00, 0x00, 0x00, 0x00, // 已初始化数据大小
		0x00, 0x00, 0x00, 0x00, // 未初始化数据大小
		0x10, 0x10, 0x00, 0x00, // 入口点
		0x10, 0x10, 0x00, 0x00, // 代码基址
	}
	
	// 节表
	sectionTable := []byte{
		0x2E, 0x74, 0x65, 0x78, 0x74, 0x00, 0x00, 0x00, // .text
		0x00, 0x10, 0x00, 0x00, // 虚拟大小
		0x10, 0x10, 0x00, 0x00, // 虚拟地址
		0x20, 0x00, 0x00, 0x00, // 原始大小
		0x20, 0x00, 0x00, 0x00, // 原始偏移
		0x00, 0x00, 0x00, 0x00, // 重定位表偏移
		0x00, 0x00, 0x00, 0x00, // 行号表偏移
		0x00, 0x00, // 重定位项数
		0x00, 0x00, // 行号项数
		0x20, 0x00, 0x00, 0x00, // 特性
	}
	
	// 数据段
	dataSection := []byte{
		0x30, 0x30, 0x30, 0x0A, // "000\n"
		0x00, 0x00, 0x00, 0x00, // 缓冲区
	}
	
	// 合并所有部分
	var code []byte
	code = append(code, mzHeader...)
	code = append(code, make([]byte, 0x40-len(mzHeader))...)
	code = append(code, peHeader...)
	code = append(code, optionalHeader...)
	code = append(code, sectionTable...)
	code = append(code, make([]byte, 0x100-len(code))...)
	code = append(code, machineCode...)
	code = append(code, dataSection...)
	
	return code, nil
}

const version = "1.0.0"

func main() {
	// 解析命令行参数
	inputFile := ""
	platform := "windows-x86_64" // 默认平台
	outputFile := "output.exe"
	debug := false
	optimize := false
	showHelp := false
	showVersion := false

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-h", "--help":
			showHelp = true
		case "-v", "--version":
			showVersion = true
		case "--platform":
			if i+1 < len(os.Args) {
				platform = os.Args[i+1]
				i++
			}
		case "-o":
			if i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			}
		case "--debug":
			debug = true
		case "--optimize":
			optimize = true
		default:
			if inputFile == "" {
				inputFile = os.Args[i]
			}
		}
	}

	// 显示版本信息
	if showVersion {
		fmt.Printf("zero-compiler version %s\n", version)
		os.Exit(0)
	}

	// 显示帮助信息
	if showHelp {
		fmt.Println("零点语言编译器")
		fmt.Println("用法: zero-compiler [选项] <input_file>")
		fmt.Println("选项:")
		fmt.Println("  -h, --help            显示此帮助信息")
		fmt.Println("  -v, --version         显示版本信息")
		fmt.Println("  --platform <platform> 指定目标平台")
		fmt.Println("  -o <output_file>      指定输出文件")
		fmt.Println("  --debug               启用调试模式")
		fmt.Println("  --optimize            启用优化")
		fmt.Println("支持的平台:")
		fmt.Println("  windows-x86_64: Windows 64位")
		fmt.Println("  linux-x86_64: Linux 64位")
		fmt.Println("  arm64: ARM64 (Android, 树莓派等)")
		fmt.Println("  cortex-m: ARM Cortex-M 单片机")
		os.Exit(0)
	}

	// 检查输入文件
	if inputFile == "" {
		fmt.Println("错误: 请指定输入文件")
		fmt.Println("使用 --help 查看用法")
		os.Exit(1)
	}

	// 创建编译器
	compiler, err := NewCompiler(platform)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 编译代码
	err = compiler.Compile(inputFile, outputFile)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 调试信息
	if debug {
		fmt.Println("调试模式:")
		fmt.Printf("  输入文件: %s\n", inputFile)
		fmt.Printf("  目标平台: %s\n", platform)
		fmt.Printf("  输出文件: %s\n", outputFile)
		fmt.Printf("  优化: %t\n", optimize)
	}
}

