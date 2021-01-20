package main

import "fmt"

type Converter struct {
	pc          int
	commandType CommandType
	arg1        string
	arg2        *int
}

func NewConverter(pc int, commandType CommandType, arg1 string, arg2 *int) *Converter {
	return &Converter{pc: pc, commandType: commandType, arg1: arg1, arg2: arg2}
}

func (c *Converter) Convert() []string {
	result := []string{}
	switch c.commandType {
	case CommandArithmetic:
		result = c.arithmetic()
	case CommandPush:
		result = c.push()
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
	case "eq":
		return c.eq()
	case "lt":
		return c.lt()
	case "gt":
		return c.gt()
	default:
		return []string{}
	}
}

func (c *Converter) add() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を加算
	return append(c.binaryArithmetic("M=M+D"), c.incrementSP()...)
}

func (c *Converter) sub() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	return append(c.binaryArithmetic("M=M-D"), c.incrementSP()...)
}

func (c *Converter) eq() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := c.binaryArithmetic("D=M-D")
	// 減算結果がゼロよりゼロと等しければtrueをセット、ゼロ以外ならfalseをセット
	jumpStep := c.jumpTruth("JEQ", "JNE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, c.incrementSP()...)
	return append(c.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) lt() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := c.binaryArithmetic("D=M-D")
	// 減算結果がゼロより小さければtrueをセット、ゼロ以上ならfalseをセット
	jumpStep := c.jumpTruth("JLT", "JGE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, c.incrementSP()...)
	return append(c.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) gt() []string {
	// スタック領域の先頭の値（第一引数）からDレジスタの値（第二引数）を減算
	arithmeticStep := c.binaryArithmetic("D=M-D")
	// 減算結果がゼロより大きければtrueをセット、ゼロ以下ならfalseをセット
	jumpStep := c.jumpTruth("JGT", "JLE")

	arithmeticStep = append(arithmeticStep, jumpStep...)
	arithmeticStep = append(arithmeticStep, c.incrementSP()...)
	return append(c.returnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) binaryArithmetic(arithmeticStep string) []string {
	return []string{
		// 第二引数を取得
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値（第二引数）をDレジスタにセット
		// 第一引数を取得＆算術演算
		"@SP",          // AレジスタにアドレスSPをセット
		"AM=M-1",       // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		arithmeticStep, // 算術演算
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
	if c.arg1 == "constant" {
		return c.pushConstant()
	}

	return []string{}
}

func (c *Converter) pushConstant() []string {
	acommand := fmt.Sprintf("@%d", *c.arg2)
	incrementSP := c.incrementSP()

	result := []string{
		acommand, // Aレジスタに定数をセット
		"D=A",    // Dレジスタへ、Aレジスタの値（直前でセットした定数）をセット
		"@SP",    // AレジスタにアドレスSPをセット
		"A=M",    // AレジスタにSPの値をセット
		"M=D",    // スタック領域へ、Dレジスタの値（最初にセットした定数）をセット
	}
	result = append(result, incrementSP...)
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

type ConverterInitializer struct{}

func (ci *ConverterInitializer) Initialize() []string {
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
