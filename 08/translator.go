package main

import (
	"fmt"
	"path/filepath"
	"sort"
)

type Translators struct {
	translators []*Translator
	moduleName  string
	hasInit     HasInit
}

func NewTranslators(filename string, hasInit HasInit) *Translators {
	moduleName := filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
	return &Translators{translators: []*Translator{}, moduleName: moduleName, hasInit: hasInit}
}

func (ts *Translators) Add(command *Command) {
	const uninitializedPC = -1
	translator := NewTranslator(uninitializedPC, command.commandType, command.arg1, command.arg2, &ts.moduleName)
	ts.translators = append(ts.translators, translator)
}

func (ts *Translators) TranslateAll() []string {
	ti := &TranslatorInitializer{hasInit: ts.hasInit}
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
	case CommandFunction:
		return t.function()
	case CommandReturn:
		return t.returnFunction()
	case CommandCall:
		return t.call()
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
	// 「x=y」ならtrueをセット、そうでなければfalseをセット
	return t.compareBinary("JEQ")
}

func (t *Translator) lt() []string {
	// 「x<y」ならtrueをセット、そうでなければfalseをセット
	return t.compareBinary("JLT")
}

func (t *Translator) gt() []string {
	// 「x>y」ならtrueをセット、そうでなければfalseをセット
	return t.compareBinary("JGT")
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

// 2値を比較し、比較結果(true/false)をスタックに積む
func (t *Translator) compareBinary(condition string) []string {
	arithmeticStep := []string{}
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	sub := t.binaryFunction("D=M-D")
	arithmeticStep = append(arithmeticStep, sub...)
	// Dレジスタに格納した減算結果と引数の条件を比較＆true/falseをDレジスタにセット
	jumpTruth := t.jumpTruth(condition)
	arithmeticStep = append(arithmeticStep, jumpTruth...)
	// Dレジスタに格納されたtrue/falseをスタックに積む
	arithmeticStep = append(arithmeticStep, t.dRegisterToStack()...)
	// スタックに値を積んだので、スタックポインタをインクリメントしておく
	arithmeticStep = append(arithmeticStep, t.incrementSP()...)

	// true/falseセット後のリターンアドレスの相対位置
	address := len(sub) + len(jumpTruth)

	result := []string{}
	result = append(result, t.returnFromJumpTruth(address)...)
	result = append(result, arithmeticStep...)
	return result
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

// Dレジスタの値を参照し、条件に合致したらtrue、そうでなければfalseをDレジスタにセットする
func (t *Translator) jumpTruth(condition string) []string {
	trueJump := fmt.Sprintf("D;%s", condition)
	return []string{
		"@TRUE",  // AレジスタにTRUEラベルをセット
		trueJump, // 条件に合致したらTRUEラベルにジャンプ
		"@FALSE", // AレジスタにFALSEラベルをセット
		"0;JMP",  // 条件に合致しなかったらFALSEラベルにジャンプ
	}
}

func (t *Translator) returnFromJumpTruth(afterStepCount int) []string {
	returnStep := []string{
		"D=A",  // Dレジスタにリターンアドレスをセット
		"@R14", // AレジスタにアドレスR14をセット
		"M=D",  // R14にリターンアドレスをセット
	}

	// リターンアドレスのセットに必要なステップ数
	// +1 をしているのはリターンアドレスのAレジスタへのセット処理がreturnStepに含まれていないため
	returnStepCount := len(returnStep) + 1

	// リターンアドレスは後続のステップ数を加味して算出
	returnAddressInt := t.pc + afterStepCount + returnStepCount
	returnAddress := fmt.Sprintf("@%d", returnAddressInt)

	result := []string{}
	result = append(result, returnAddress)
	result = append(result, returnStep...)
	return result
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

// TODO R14じゃなくてR13を使う
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
		label,   // Aレジスタにラベルをセット
		"D;JNE", // Dレジスタの値（スタックの先頭の値）がゼロ以外ならラベルにジャンプ
	}
}

func (t *Translator) function() []string {
	// 関数名のラベルを定義
	label := fmt.Sprintf("(%s)", t.arg1)
	result := []string{label}

	// 指定されたローカル変数の数だけ領域を確保して初期化
	for i := 0; i < *t.arg2; i++ {
		initLocal := []string{
			"@SP", // AレジスタにアドレスSPをセット
			"A=M", // // AレジスタにSPの値をセット
			"M=0", // スタックの先頭の値に0をセット
		}
		result = append(result, initLocal...)
		result = append(result, t.incrementSP()...)
	}

	return result
}

func (t *Translator) returnFunction() []string {
	// FRAME=LCL
	// R13にFRAMEの値を格納して参照できるようにしておく
	frame := []string{
		"@LCL", // AレジスタにアドレスLCLをセット
		"D=M",  // LCLの値をDレジスタにセット
		"@R13", // AレジスタにアドレスR13をセット
		"M=D",  // Dレジスタ（LCLの値）をR13にセット
	}

	// RET = *(FRAME-5)
	// R14にリターンアドレスを格納して最後に使う
	ret := t.restoreByFrame("R14", 5)

	// *ARG = pop()
	retValue := []string{
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		"@ARG",   // AレジスタにアドレスARGをセット
		"A=M",    // AレジスタにARGの値をセット
		"M=D",    // 返り値（スタック領域の先頭の値）をARGにセット
	}

	// SP = ARG+1 : 呼び出すもとのSPを戻り値の直後のアドレスに変更
	sp := []string{
		"D=A",   // ARGのアドレスをDレジスタにセット
		"@SP",   // AレジスタにアドレスSPをセット
		"M=D+1", // 「ARG+1」を計算してSPにセット
	}

	// THAT = *(FRAME-1)
	that := t.restoreByFrame("THAT", 1)

	// THIS = *(FRAME-2)
	this := t.restoreByFrame("THIS", 2)

	// ARG = *(FRAME-3)
	arg := t.restoreByFrame("ARG", 3)

	// LCL = *(FRAME-4)
	lcl := t.restoreByFrame("LCL", 4)

	// goto RET
	gotoRet := []string{
		"@R14",  // AレジスタにアドレスR14をセット
		"A=M",   // Aレジスタにリターンアドレスをセット
		"0;JMP", // リターンアドレスにジャンプ
	}

	result := []string{}
	result = append(result, frame...)
	result = append(result, ret...)
	result = append(result, retValue...)
	result = append(result, sp...)
	result = append(result, that...)
	result = append(result, this...)
	result = append(result, arg...)
	result = append(result, lcl...)
	result = append(result, gotoRet...)

	return result
}

// call <func-name> <arg-count>
// call Main.add 1
func (t *Translator) call() []string {
	functionName := t.arg1

	// push return-address
	label := fmt.Sprintf("RETURN-ADDRESS$%s$%s$%d", *t.moduleName, functionName, t.pc)
	returnAddress := fmt.Sprintf("@%s", label)
	ret := []string{
		returnAddress, // リターンアドレスをAレジスタにセット
		"D=A",         // リターンアドレスを取得してDレジスタにセット
	}
	ret = append(ret, t.dRegisterToStack()...)
	ret = append(ret, t.incrementSP()...)

	// push LCL
	callerLCL := t.storeCallerState("LCL")
	// push ARG
	callerARG := t.storeCallerState("ARG")
	// push THIS
	callerTHIS := t.storeCallerState("THIS")
	// push THAT
	callerTHAT := t.storeCallerState("THAT")

	// @ARG = SP-n-5
	argCount := fmt.Sprintf("@%d", *t.arg2)
	const callerStateCount = "@5" // 呼び出し元の関数の状態の数=RTN+LCL+ARG+THIS+THAT=5
	arg := []string{
		argCount,         // 関数の引数の数(n)をAレジスタにセット
		"D=A",            // 関数の引数の数(n)をAレジスタから取得してDレジスタにセット
		callerStateCount, // 呼び出し元の関数の状態の数(=5)をAレジスタにセット
		"D=D+A",          // 「n+5」を算出してDレジスタにセット
		"@SP",            // AレジスタにアドレスSPをセット
		"D=M-D",          // 「SP-n-5」を算出してDレジスタにセット
		"@ARG",           // AレジスタにアドレスARGをセット
		"M=D",            // 「SP-n-5」をARGにセット
	}

	// @LCL=SP
	lcl := []string{
		"@SP",  // AレジスタにアドレスSPをセット
		"D=M",  // SPの値をDレジスタにセット
		"@LCL", // AレジスタにアドレスARGをセット
		"M=D",  // SPの値をARGにセット
	}

	// goto f
	functionLabel := fmt.Sprintf("@%s", functionName)
	gotoFunction := []string{
		functionLabel,
		"0;JMP",
	}

	// (return-address)
	returnAddressLabel := []string{fmt.Sprintf("(%s)", label)}

	result := []string{}
	result = append(result, ret...)
	result = append(result, callerLCL...)
	result = append(result, callerARG...)
	result = append(result, callerTHIS...)
	result = append(result, callerTHAT...)
	result = append(result, arg...)
	result = append(result, lcl...)
	result = append(result, gotoFunction...)
	result = append(result, returnAddressLabel...)

	return result
}

func (t *Translator) storeCallerState(label string) []string {
	labelAddress := fmt.Sprintf("@%s", label)
	result := []string{
		labelAddress, // 指定したラベルのアドレスをAレジスタにセット
		"D=M",        // 取得した値をDレジスタにセット
	}
	result = append(result, t.dRegisterToStack()...)
	result = append(result, t.incrementSP()...)
	return result
}

func (t *Translator) restoreByFrame(definedLabel string, frameIndex int) []string {
	index := fmt.Sprintf("@%d", frameIndex)
	label := fmt.Sprintf("@%s", definedLabel)
	return []string{
		"@R13",  // AレジスタにアドレスR13（FRAMEのアドレス）をセット
		"D=M",   // FRAMEの値をDレジスタにセット
		index,   // Aレジスタに定数をセット
		"A=D-A", // 「FRAME-frameIndex」を計算してAレジスタにセット
		"D=M",   // *(FRAME-frameIndex)の値をDレジスタにセット
		label,   // Aレジスタに定義済みラベルをセット
		"M=D",   // Dレジスタの値をTHATにセット
	}
}

type TranslatorInitializer struct {
	hasInit HasInit
}

func (ti *TranslatorInitializer) initializeHeader() []string {
	//return []string{}
	result := []string{}
	result = append(result, ti.initializeLabels()...)
	result = append(result, ti.initializeSysInit()...) // これは必ず最後に追加する
	return result
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

// 初期化処理が終わったら最後にSys.initを実行する
func (ti *TranslatorInitializer) initializeSysInit() []string {
	// 以前のテストケースも動くように「function Sys.init 0」が定義されてるときだけSys.initを呼ぶ
	if ti.hasInit {
		return []string{"@Sys.init"}
	}
	return []string{}
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

// Dレジスタにtrueをセットする
func (ti *TranslatorInitializer) initializeTRUE() []string {
	return []string{
		"(TRUE)",
		"  D=-1",  // Dレジスタにtrueをセット
		"  @R14",  // AレジスタにアドレスR14をセット
		"  A=M",   // Aレジスタにリターンアドレスをセット
		"  0;JMP", // リターンアドレスにジャンプ
	}
}

// Dレジスタにfalseをセットする
func (ti *TranslatorInitializer) initializeFALSE() []string {
	return []string{
		"(FALSE)",
		"  D=0",   // Dレジスタにfalseをセット
		"  @R14",  // AレジスタにアドレスR14をセット
		"  A=M",   // Aレジスタにリターンアドレスをセット
		"  0;JMP", // リターンアドレスにジャンプ
	}
}
