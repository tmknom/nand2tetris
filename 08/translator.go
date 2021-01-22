package main

import (
	"fmt"
	"path/filepath"
	"sort"
)

type Translators struct {
	translators []*Translator
	moduleName  string
}

func NewTranslators(filename string) *Translators {
	moduleName := filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
	return &Translators{translators: []*Translator{}, moduleName: moduleName}
}

func (ts *Translators) Add(command *Command) {
	const uninitializedPC = -1
	translator := NewTranslator(uninitializedPC, command.commandType, command.arg1, command.arg2, &ts.moduleName)
	ts.translators = append(ts.translators, translator)
}

func (ts *Translators) TranslatorAll() []string {
	ti := &TranslatorInitializer{}
	result := ti.initializeHeader()

	for _, translator := range ts.translators {
		translator.setPC(len(result))
		assembler := translator.Translate()
		result = append(result, assembler...)
	}

	return append(result, ti.initializeFooter()...)
}

type Translator struct {
	pc          int
	commandType CommandType
	arg1        string
	arg2        *int
	moduleName  *string
}

const (
	basePointerAddress = 3
	baseTempAddress    = 5
	baseStaticAddress  = 16
)

func NewTranslator(pc int, commandType CommandType, arg1 string, arg2 *int, moduleName *string) *Translator {
	return &Translator{pc: pc, commandType: commandType, arg1: arg1, arg2: arg2, moduleName: moduleName}
}

func (t *Translator) setPC(pc int) {
	t.pc = pc
}

func (t *Translator) Translate() []string {
	switch t.commandType {
	case CommandArithmetic:
		return t.arithmetic()
	case CommandPush:
		return t.push()
	case CommandPop:
		return t.pop()
	case CommandLabel:
		return t.label()
	case CommandGoto:
		return t.labelGoto()
	case CommandIf:
		return t.ifGoto()
	default:
		return []string{}
	}
}

func (t *Translator) arithmetic() []string {
	switch t.arg1 {
	case "add":
		return t.add()
	case "sub":
		return t.sub()
	case "neg":
		return t.neg()
	case "eq":
		return t.eq()
	case "lt":
		return t.lt()
	case "gt":
		return t.gt()
	case "and":
		return t.and()
	case "or":
		return t.or()
	case "not":
		return t.not()
	default:
		return []string{}
	}
}

func (t *Translator) add() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を加算
	return append(t.binaryFunction("M=M+D"), t.incrementSP()...)
}

func (t *Translator) sub() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	return append(t.binaryFunction("M=M-D"), t.incrementSP()...)
}

func (t *Translator) neg() []string {
	// スタック領域の先頭の値（第一引数）の反転
	return append(t.unaryFunction("M=-M"), t.incrementSP()...)
}

func (t *Translator) eq() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := t.binaryFunction("D=M-D")
	// 減算結果がゼロよりゼロと等しければtrueをセット、ゼロ以外ならfalseをセット
	jumpStep := t.jumpTruth("JEQ", "JNE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, t.incrementSP()...)
	return append(t.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (t *Translator) lt() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := t.binaryFunction("D=M-D")
	// 減算結果がゼロより小さければtrueをセット、ゼロ以上ならfalseをセット
	jumpStep := t.jumpTruth("JLT", "JGE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, t.incrementSP()...)
	return append(t.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (t *Translator) gt() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := t.binaryFunction("D=M-D")
	// 減算結果がゼロより大きければtrueをセット、ゼロ以下ならfalseをセット
	jumpStep := t.jumpTruth("JGT", "JLE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, t.incrementSP()...)
	return append(t.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (t *Translator) and() []string {
	// スタック領域の先頭の値（第一引数）とDレジスタの値（第二引数）の論理積
	return append(t.binaryFunction("M=D&M"), t.incrementSP()...)
}

func (t *Translator) or() []string {
	// スタック領域の先頭の値（第一引数）とDレジスタの値（第二引数）の論理和
	return append(t.binaryFunction("M=D|M"), t.incrementSP()...)
}

func (t *Translator) not() []string {
	// スタック領域の先頭の値（第一引数）の否定
	return append(t.unaryFunction("M=!M"), t.incrementSP()...)
}

// 2変数関数
func (t *Translator) binaryFunction(step string) []string {
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
func (t *Translator) unaryFunction(step string) []string {
	return []string{
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		step,     // 演算
	}
}

func (t *Translator) jumpTruth(trueMnemonic string, falseMnemonic string) []string {
	trueJump := fmt.Sprintf("D;%s", trueMnemonic)
	falseJump := fmt.Sprintf("D;%s", falseMnemonic)
	return []string{
		"@TRUE",   // AレジスタにTRUEラベルをセット
		trueJump,  // trueMnemonicに合致したらTRUEラベルにジャンプ
		"@FALSE",  // AレジスタにFALSEラベルをセット
		falseJump, // falseMnemonicに合致したらFALSEラベルにジャンプ
	}
}

func (t *Translator) returnAddress(afterStepCount int) []string {
	// ステップ数の微調整
	const tweakStepCount = 2
	// リターンアドレスは後続のステップ数を加味して算出
	returnAddressInt := tweakStepCount + afterStepCount + t.pc
	returnAddress := fmt.Sprintf("@%d", returnAddressInt)
	return []string{
		returnAddress, // Aレジスタにリターンアドレスをセット
		"D=A",         // Dレジスタにリターンアドレスをセット
		"@R15",        // AレジスタにアドレスR15をセット
		"M=D",         // R15にリターンアドレスをセット
	}
}

func (t *Translator) push() []string {
	switch t.arg1 {
	case "constant":
		return t.pushConstant()
	case "local":
		return t.pushLocal()
	case "argument":
		return t.pushArgument()
	case "this":
		return t.pushThis()
	case "that":
		return t.pushThat()
	case "temp":
		return t.pushTemp()
	case "pointer":
		return t.pushPointer()
	case "static":
		return t.pushStatic()
	default:
		return []string{}
	}
}

func (t *Translator) pushConstant() []string {
	constant := fmt.Sprintf("@%d", *t.arg2)
	result := []string{
		constant, // Aレジスタに定数をセット
		"D=A",    // Dレジスタへ、Aレジスタの値（直前でセットした定数）をセット
	}

	// スタックにDレジスタの値を積む
	result = append(result, t.dRegisterToStack()...)
	// スタックポインタのインクリメント
	result = append(result, t.incrementSP()...)
	return result
}

func (t *Translator) pushLocal() []string {
	return t.pushLabel("LCL")
}

func (t *Translator) pushArgument() []string {
	return t.pushLabel("ARG")
}

func (t *Translator) pushThis() []string {
	return t.pushLabel("THIS")
}

func (t *Translator) pushThat() []string {
	return t.pushLabel("THAT")
}

func (t *Translator) pushTemp() []string {
	return t.pushAddress(baseTempAddress)
}

func (t *Translator) pushPointer() []string {
	return t.pushAddress(basePointerAddress)
}

func (t *Translator) pushStatic() []string {
	return t.pushAddress(baseStaticAddress)
}

func (t *Translator) pushAddress(baseAddress int) []string {
	address := fmt.Sprintf("@%d", *t.arg2+baseAddress)
	result := []string{
		address, // Aレジスタにアドレスをセット
		"D=M",   // 指定したアドレスから取得した値をDレジスタにセット
	}

	// スタックにDレジスタの値を積む
	result = append(result, t.dRegisterToStack()...)
	// スタックポインタのインクリメント
	result = append(result, t.incrementSP()...)
	return result
}

func (t *Translator) pushLabel(label string) []string {
	// 取得先アドレスを算出して、取得した値をDレジスタにセット
	index := fmt.Sprintf("@%d", *t.arg2)
	baseAddress := fmt.Sprintf("@%s", label)
	result := []string{
		index,       // インデックスをAレジスタにセット
		"D=A",       // Dレジスタへ、Aレジスタの値（インデックス）をセット
		baseAddress, // Aレジスタにベースアドレスをセット
		"A=D+M",     // 取得先アドレス（インデックス+ベースアドレス）を算出してAレジスタにセット
		"D=M",       // 取得した値をDレジスタにセット
	}

	// スタックにDレジスタの値を積む
	result = append(result, t.dRegisterToStack()...)
	// スタックポインタのインクリメント
	result = append(result, t.incrementSP()...)
	return result
}

func (t *Translator) pop() []string {
	switch t.arg1 {
	case "local":
		return t.popLocal()
	case "argument":
		return t.popArgument()
	case "this":
		return t.popThis()
	case "that":
		return t.popThat()
	case "temp":
		return t.popTemp()
	case "pointer":
		return t.popPointer()
	case "static":
		return t.popStatic()
	default:
		return []string{}
	}
}

func (t *Translator) popLocal() []string {
	return t.popLabel("LCL")
}

func (t *Translator) popArgument() []string {
	return t.popLabel("ARG")
}

func (t *Translator) popThis() []string {
	return t.popLabel("THIS")
}

func (t *Translator) popThat() []string {
	return t.popLabel("THAT")
}

func (t *Translator) popTemp() []string {
	return t.popAddress(baseTempAddress)
}

func (t *Translator) popPointer() []string {
	return t.popAddress(basePointerAddress)
}

func (t *Translator) popStatic() []string {
	return t.popAddress(baseStaticAddress)
}

func (t *Translator) popAddress(baseAddress int) []string {
	address := fmt.Sprintf("@%d", *t.arg2+baseAddress)

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

func (t *Translator) popLabel(label string) []string {
	index := fmt.Sprintf("@%d", *t.arg2)
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
func (t *Translator) incrementSP() []string {
	return []string{
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M+1", // SPの値をインクリメント
	}
}

// スタックにDレジスタの値を積む
func (t *Translator) dRegisterToStack() []string {
	return []string{
		"@SP", // AレジスタにアドレスSPをセット
		"A=M", // AレジスタにSPの値をセット
		"M=D", // スタック領域へ、Dレジスタの値（最初にセットした定数）をセット
	}
}

func (t *Translator) label() []string {
	label := fmt.Sprintf("(%s$%s)", *t.moduleName, t.arg1)
	return []string{label}
}

func (t *Translator) labelGoto() []string {
	label := fmt.Sprintf("@%s$%s", *t.moduleName, t.arg1)
	return []string{
		label,
		"0;JMP",
	}
}

func (t *Translator) ifGoto() []string {
	label := fmt.Sprintf("@%s$%s", *t.moduleName, t.arg1)
	return []string{
		// スタック領域の先頭の値をDレジスタにセット
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		// スタックの先頭の値を使って分岐
		label, // Aレジスタにラベルをセット
		"D;JNE", // Dレジスタの値（スタックの先頭の値）がゼロ以外ならラベルにジャンプ
	}
}

type TranslatorInitializer struct{}

func (ti *TranslatorInitializer) initializeHeader() []string {
	return ti.initializeLabels()
}

func (ti *TranslatorInitializer) initializeLabels() []string {
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
		result = append(result, ti.initializeLabel(labels[address], address)...)
	}
	return result
}

func (ti *TranslatorInitializer) initializeLabel(name string, address int) []string {
	constant := fmt.Sprintf("@%d", address)
	label := fmt.Sprintf("@%s", name)
	return []string{
		constant, // Aレジスタにアドレスを定数としてセット
		"D=A",    // Aレジスタの値をDレジスタにセット
		label,    // Aレジスタにラベルをセット
		"M=D",    // 指定したラベルにDレジスタの値（アドレス）をセット
	}
}

func (ti *TranslatorInitializer) initializeFooter() []string {
	endStep := ti.initializeEndStep()
	endLabel := ti.initializeEND()
	trueLabel := ti.initializeTRUE()
	falseLabel := ti.initializeFALSE()

	result := []string{}
	result = append(result, endStep...)
	result = append(result, trueLabel...)
	result = append(result, falseLabel...)

	// この処理は最後に追加する
	result = append(result, endLabel...)

	return result
}

func (ti *TranslatorInitializer) initializeEndStep() []string {
	return []string{
		"@END",
		"0;JMP",
	}
}

func (ti *TranslatorInitializer) initializeEND() []string {
	return []string{
		"(END)", // ENDラベル以降は何もしない
	}
}

func (ti *TranslatorInitializer) initializeTRUE() []string {
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

func (ti *TranslatorInitializer) initializeFALSE() []string {
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
