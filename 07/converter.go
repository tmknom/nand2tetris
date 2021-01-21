package main

import (
	"fmt"
	"sort"
)

type Converters struct {
	converters []*Converter
}

func NewConverters() *Converters {
	return &Converters{converters: []*Converter{}}
}

func (cs *Converters) Add(command *Command) {
	const uninitializedPC = -1
	converter := NewConverter(uninitializedPC, command.commandType, command.arg1, command.arg2)
	cs.converters = append(cs.converters, converter)
}

func (cs *Converters) ConvertAll() []string {
	ci := &ConverterInitializer{}
	result := ci.initializeHeader()

	for _, converter := range cs.converters {
		converter.setPC(len(result))
		assembler := converter.Convert()
		result = append(result, assembler...)
	}

	return append(result, ci.initializeFooter()...)
}

type Converter struct {
	pc          int
	commandType CommandType
	arg1        string
	arg2        *int
}

const (
	basePointerAddress = 3
	baseTempAddress    = 5
	baseStaticAddress  = 16
)

func NewConverter(pc int, commandType CommandType, arg1 string, arg2 *int) *Converter {
	return &Converter{pc: pc, commandType: commandType, arg1: arg1, arg2: arg2}
}

func (c *Converter) setPC(pc int) {
	c.pc = pc
}

func (c *Converter) Convert() []string {
	result := []string{}
	switch c.commandType {
	case CommandArithmetic:
		result = c.arithmetic()
	case CommandPush:
		result = c.push()
	case CommandPop:
		result = c.pop()
	default:
		return result
		//return fmt.Errorf("convert failed: %s", command.raw)
	}

	return result
}

func (c *Converter) arithmetic() []string {
	switch c.arg1 {
	case "add":
		return c.add()
	case "sub":
		return c.sub()
	case "neg":
		return c.neg()
	case "eq":
		return c.eq()
	case "lt":
		return c.lt()
	case "gt":
		return c.gt()
	case "and":
		return c.and()
	case "or":
		return c.or()
	case "not":
		return c.not()
	default:
		return []string{}
	}
}

func (c *Converter) add() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を加算
	return append(c.binaryFunction("M=M+D"), c.incrementSP()...)
}

func (c *Converter) sub() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	return append(c.binaryFunction("M=M-D"), c.incrementSP()...)
}

func (c *Converter) neg() []string {
	// スタック領域の先頭の値（第一引数）の反転
	return append(c.unaryFunction("M=-M"), c.incrementSP()...)
}

func (c *Converter) eq() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := c.binaryFunction("D=M-D")
	// 減算結果がゼロよりゼロと等しければtrueをセット、ゼロ以外ならfalseをセット
	jumpStep := c.jumpTruth("JEQ", "JNE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, c.incrementSP()...)
	return append(c.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) lt() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := c.binaryFunction("D=M-D")
	// 減算結果がゼロより小さければtrueをセット、ゼロ以上ならfalseをセット
	jumpStep := c.jumpTruth("JLT", "JGE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, c.incrementSP()...)
	return append(c.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) gt() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := c.binaryFunction("D=M-D")
	// 減算結果がゼロより大きければtrueをセット、ゼロ以下ならfalseをセット
	jumpStep := c.jumpTruth("JGT", "JLE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, c.incrementSP()...)
	return append(c.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) and() []string {
	// スタック領域の先頭の値（第一引数）とDレジスタの値（第二引数）の論理積
	return append(c.binaryFunction("M=D&M"), c.incrementSP()...)
}

func (c *Converter) or() []string {
	// スタック領域の先頭の値（第一引数）とDレジスタの値（第二引数）の論理和
	return append(c.binaryFunction("M=D|M"), c.incrementSP()...)
}

func (c *Converter) not() []string {
	// スタック領域の先頭の値（第一引数）の否定
	return append(c.unaryFunction("M=!M"), c.incrementSP()...)
}

// 2変数関数
func (c *Converter) binaryFunction(step string) []string {
	return []string{
		// 第二引数を取得
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値（第二引数）をDレジスタにセット
		// 第一引数を取得＆算術演算
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		step,     // 演算
	}
}

// 1変数関数
func (c *Converter) unaryFunction(step string) []string {
	return []string{
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		step,     // 演算
	}
}

func (c *Converter) jumpTruth(trueMnemonic string, falseMnemonic string) []string {
	trueJump := fmt.Sprintf("D;%s", trueMnemonic)
	falseJump := fmt.Sprintf("D;%s", falseMnemonic)
	return []string{
		"@TRUE",   // AレジスタにTRUEラベルをセット
		trueJump,  // trueMnemonicに合致したらTRUEラベルにジャンプ
		"@FALSE",  // AレジスタにFALSEラベルをセット
		falseJump, // falseMnemonicに合致したらFALSEラベルにジャンプ
	}
}

func (c *Converter) returnAddress(afterStepCount int) []string {
	// ステップ数の微調整
	const tweakStepCount = 2
	// リターンアドレスは後続のステップ数を加味して算出
	returnAddressInt := tweakStepCount + afterStepCount + c.pc
	returnAddress := fmt.Sprintf("@%d", returnAddressInt)
	return []string{
		returnAddress, // Aレジスタにリターンアドレスをセット
		"D=A",         // Dレジスタにリターンアドレスをセット
		"@R15",        // AレジスタにアドレスR15をセット
		"M=D",         // R15にリターンアドレスをセット
	}
}

func (c *Converter) push() []string {
	switch c.arg1 {
	case "constant":
		return c.pushConstant()
	case "local":
		return c.pushLocal()
	case "argument":
		return c.pushArgument()
	case "this":
		return c.pushThis()
	case "that":
		return c.pushThat()
	case "temp":
		return c.pushTemp()
	case "pointer":
		return c.pushPointer()
	case "static":
		return c.pushStatic()
	default:
		return []string{}
	}
}

func (c *Converter) pushConstant() []string {
	constant := fmt.Sprintf("@%d", *c.arg2)
	result := []string{
		constant, // Aレジスタに定数をセット
		"D=A",    // Dレジスタへ、Aレジスタの値（直前でセットした定数）をセット
	}

	// スタックにDレジスタの値を積む
	result = append(result, c.dRegisterToStack()...)
	// スタックポインタのインクリメント
	result = append(result, c.incrementSP()...)
	return result
}

func (c *Converter) pushLocal() []string {
	return c.pushLabel("LCL")
}

func (c *Converter) pushArgument() []string {
	return c.pushLabel("ARG")
}

func (c *Converter) pushThis() []string {
	return c.pushLabel("THIS")
}

func (c *Converter) pushThat() []string {
	return c.pushLabel("THAT")
}

func (c *Converter) pushTemp() []string {
	return c.pushAddress(baseTempAddress)
}

func (c *Converter) pushPointer() []string {
	return c.pushAddress(basePointerAddress)
}

func (c *Converter) pushStatic() []string {
	return c.pushAddress(baseStaticAddress)
}

func (c *Converter) pushAddress(baseAddress int) []string {
	address := fmt.Sprintf("@%d", *c.arg2+baseAddress)
	result := []string{
		address, // Aレジスタにアドレスをセット
		"D=M",   // 指定したアドレスから取得した値をDレジスタにセット
	}

	// スタックにDレジスタの値を積む
	result = append(result, c.dRegisterToStack()...)
	// スタックポインタのインクリメント
	result = append(result, c.incrementSP()...)
	return result
}

func (c *Converter) pushLabel(label string) []string {
	// 取得先アドレスを算出して、取得した値をDレジスタにセット
	index := fmt.Sprintf("@%d", *c.arg2)
	baseAddress := fmt.Sprintf("@%s", label)
	result := []string{
		index,       // インデックスをAレジスタにセット
		"D=A",       // Dレジスタへ、Aレジスタの値（インデックス）をセット
		baseAddress, // Aレジスタにベースアドレスをセット
		"A=D+M",     // 取得先アドレス（インデックス+ベースアドレス）を算出してAレジスタにセット
		"D=M",       // 取得した値をDレジスタにセット
	}

	// スタックにDレジスタの値を積む
	result = append(result, c.dRegisterToStack()...)
	// スタックポインタのインクリメント
	result = append(result, c.incrementSP()...)
	return result
}

func (c *Converter) pop() []string {
	switch c.arg1 {
	case "local":
		return c.popLocal()
	case "argument":
		return c.popArgument()
	case "this":
		return c.popThis()
	case "that":
		return c.popThat()
	case "temp":
		return c.popTemp()
	case "pointer":
		return c.popPointer()
	case "static":
		return c.popStatic()
	default:
		return []string{}
	}
}

func (c *Converter) popLocal() []string {
	return c.popLabel("LCL")
}

func (c *Converter) popArgument() []string {
	return c.popLabel("ARG")
}

func (c *Converter) popThis() []string {
	return c.popLabel("THIS")
}

func (c *Converter) popThat() []string {
	return c.popLabel("THAT")
}

func (c *Converter) popTemp() []string {
	return c.popAddress(baseTempAddress)
}

func (c *Converter) popPointer() []string {
	return c.popAddress(basePointerAddress)
}

func (c *Converter) popStatic() []string {
	return c.popAddress(baseStaticAddress)
}

func (c *Converter) popAddress(baseAddress int) []string {
	address := fmt.Sprintf("@%d", *c.arg2+baseAddress)

	result := []string{
		// スタック領域の先頭の値をDレジスタにセット
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		// 指定アドレスにスタック領域の先頭の値をセット
		address, // Aレジスタにアドレスをセット
		"M=D",   // 保存先アドレスにDレジスタの値（スタック領域の先頭の値）をセット
	}
	return result
}

func (c *Converter) popLabel(label string) []string {
	index := fmt.Sprintf("@%d", *c.arg2)
	baseAddress := fmt.Sprintf("@%s", label)

	result := []string{
		// 保存先アドレスを算出して、R14に一時的にセット
		index,       // インデックスをAレジスタにセット
		"D=A",       // Dレジスタへ、Aレジスタの値（インデックス）をセット
		baseAddress, // Aレジスタにベースアドレスをセット
		"D=D+M",     // 保存先アドレス（インデックス+ベースアドレス）を算出してDレジスタにセット
		"@R14",      // AレジスタにアドレスR14をセット
		"M=D",       // R14に一時的に保存先アドレスをセット
		// スタック領域の先頭の値をDレジスタにセット
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		// 保存先アドレスにスタック領域の先頭の値をセット
		"@R14", // AレジスタにアドレスR14をセット
		"A=M",  // Aレジスタに保存先アドレスをセット
		"M=D",  // 保存先アドレスにDレジスタの値（スタック領域の先頭の値）をセット
	}
	return result
}

// スタックポインタのインクリメント
// スタックに値を積んだら忘れずに実施する
func (c *Converter) incrementSP() []string {
	return []string{
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M+1", // SPの値をインクリメント
	}
}

// スタックにDレジスタの値を積む
func (c *Converter) dRegisterToStack() []string {
	return []string{
		"@SP", // AレジスタにアドレスSPをセット
		"A=M", // AレジスタにSPの値をセット
		"M=D", // スタック領域へ、Dレジスタの値（最初にセットした定数）をセット
	}
}

type ConverterInitializer struct{}

func (ci *ConverterInitializer) initializeHeader() []string {
	return ci.initializeLabels()
}

func (ci *ConverterInitializer) initializeLabels() []string {
	labels := map[int]string{
		256:  "SP",
		300:  "LCL",
		400:  "ARG",
		3000: "THIS",
		3010: "THAT",
	}

	// テストコードの実行を安定させるため、意図的にmapに順序概念を追加
	addresses := []int{}
	for address := range labels {
		addresses = append(addresses, address)
	}
	sort.Ints(addresses)

	result := []string{}
	for _, address := range addresses {
		result = append(result, ci.initializeLabel(labels[address], address)...)
	}
	return result
}

func (ci *ConverterInitializer) initializeLabel(name string, address int) []string {
	constant := fmt.Sprintf("@%d", address)
	label := fmt.Sprintf("@%s", name)
	return []string{
		constant, // Aレジスタにアドレスを定数としてセット
		"D=A",    // Aレジスタの値をDレジスタにセット
		label,    // Aレジスタにラベルをセット
		"M=D",    // 指定したラベルにDレジスタの値（アドレス）をセット
	}
}

func (ci *ConverterInitializer) initializeFooter() []string {
	endStep := ci.initializeEndStep()
	endLabel := ci.initializeEND()
	trueLabel := ci.initializeTRUE()
	falseLabel := ci.initializeFALSE()

	result := []string{}
	result = append(result, endStep...)
	result = append(result, trueLabel...)
	result = append(result, falseLabel...)

	// この処理は最後に追加する
	result = append(result, endLabel...)

	return result
}

func (ci *ConverterInitializer) initializeEndStep() []string {
	return []string{
		"@END",
		"0;JMP",
	}
}

func (ci *ConverterInitializer) initializeEND() []string {
	return []string{
		"(END)", // ENDラベル以降は何もしない
	}
}

func (ci *ConverterInitializer) initializeTRUE() []string {
	return []string{
		"(TRUE)",
		"  @SP",   // AレジスタにアドレスSPをセット
		"  A=M",   // AレジスタにSPの値をセット
		"  M=-1",  // スタックの先頭の値にtrueをセット
		"  @R15",  // AレジスタにアドレスR15をセット
		"  A=M",   // Aレジスタにリターンアドレスをセット
		"  0;JMP", // リターンアドレスにジャンプ
	}
}

func (ci *ConverterInitializer) initializeFALSE() []string {
	return []string{
		"(FALSE)",
		"  @SP",   // AレジスタにアドレスSPをセット
		"  A=M",   // AレジスタにSPの値をセット
		"  M=0",   // スタックの先頭の値にfalseをセット
		"  @R15",  // AレジスタにアドレスR15をセット
		"  A=M",   // Aレジスタにリターンアドレスをセット
		"  0;JMP", // リターンアドレスにジャンプ
	}
}
