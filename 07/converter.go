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
		result = c.convertArithmetic()
	case CommandPush:
		result = c.convertPush()
	default:
		return result
		//return fmt.Errorf("convert failed: %s", command.raw)
	}

	return result
}

func (c *Converter) convertArithmetic() []string {
	switch c.arg1 {
	case "add":
		return c.convertAdd()
	case "eq":
		return c.convertEq()
	case "lt":
		return c.convertLt()
	case "gt":
		return c.convertGt()
	default:
		return []string{}
	}
}

func (c *Converter) convertAdd() []string {
	return []string{
		// addの第一引数を取得
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		// addの第一引数を取得＆第一引数を加算してスタック領域の先頭の値を更新
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"M=D+M",  // スタック領域の先頭の値とDレジスタの値を加算
		// スタック領域の先頭アドレスをインクリメント
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M+1", // SPの値をインクリメント
	}
}

func (c *Converter) convertEq() []string {
	arithmeticStep := []string{
		// eqの第一引数を取得
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		// eqの第二引数を取得＆第一引数と減算
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M-D",  // スタック領域の先頭の値から第一引数を減算してDレジスタにセット
		// true/falseをセット
		"@EQ",   // AレジスタにEQラベルをセット
		"D;JEQ", // Dレジスタの値（減算結果）がゼロと等しければEQラベルにジャンプ
		"@NEQ",  // AレジスタにNEQラベルをセット
		"D;JNE", // Dレジスタの値（減算結果）がゼロ以外と等しければNEQラベルにジャンプ
		// スタック領域の先頭アドレスをインクリメント
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M+1", // SPの値をインクリメント
	}

	return append(c.convertReturnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) convertLt() []string {
	arithmeticStep := []string{
		// eqの第一引数を取得
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		// eqの第二引数を取得＆第一引数と減算
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M-D",  // スタック領域の先頭の値から第一引数を減算してDレジスタにセット
		// true/falseをセット
		"@EQ",   // AレジスタにEQラベルをセット
		"D;JLT", // Dレジスタの値（減算結果）がゼロより小さければEQラベルにジャンプ
		"@NEQ",  // AレジスタにNEQラベルをセット
		"D;JGE", // Dレジスタの値（減算結果）がゼロ以上ならばNEQラベルにジャンプ
		// スタック領域の先頭アドレスをインクリメント
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M+1", // SPの値をインクリメント
	}

	return append(c.convertReturnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) convertGt() []string {
	arithmeticStep := []string{
		// eqの第一引数を取得
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M",    // スタック領域の先頭の値をDレジスタにセット
		// eqの第二引数を取得＆第一引数と減算
		"@SP",    // AレジスタにアドレスSPをセット
		"AM=M-1", // スタック領域の先頭アドレスをデクリメントしてAレジスタにセット
		"D=M-D",  // スタック領域の先頭の値から第一引数を減算してDレジスタにセット
		// true/falseをセット
		"@EQ",   // AレジスタにEQラベルをセット
		"D;JGT", // Dレジスタの値（減算結果）がゼロより大きければEQラベルにジャンプ
		"@NEQ",  // AレジスタにNEQラベルをセット
		"D;JLE", // Dレジスタの値（減算結果）がゼロ以下ならばNEQラベルにジャンプ
		// スタック領域の先頭アドレスをインクリメント
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M+1", // SPの値をインクリメント
	}

	return append(c.convertReturnAddress(len(arithmeticStep)), arithmeticStep...)
}

func (c *Converter) convertReturnAddress(afterStepCount int) []string {
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

func (c *Converter) convertPush() []string {
	if c.arg1 == "constant" {
		return c.convertPushConstant()
	}

	return []string{}
}

func (c *Converter) convertPushConstant() []string {
	acommand := fmt.Sprintf("@%d", *c.arg2)

	return []string{
		acommand, // Aレジスタに定数をセット
		"D=A",    // Dレジスタへ、Aレジスタの値（直前でセットした定数）をセット
		"@SP",    // AレジスタにアドレスSPをセット
		"A=M",    // AレジスタにSPの値をセット
		"M=D",    // スタック領域へ、Dレジスタの値（最初にセットした定数）をセット
		"@SP",    // AレジスタにアドレスSPをセット
		"M=M+1",  // SPの値をインクリメント
	}
}

type ConverterInitializer struct{}

func (ci *ConverterInitializer) Initialize() []string {
	endStep := ci.initializeEndStep()
	end := ci.initializeEND()
	eq := ci.initializeEQ()
	neq := ci.initializeNEQ()

	result := []string{}
	result = append(result, endStep...)
	result = append(result, eq...)
	result = append(result, neq...)

	// この処理は最後に追加する
	result = append(result, end...)

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

func (ci *ConverterInitializer) initializeEQ() []string {
	return []string{
		"(EQ)",
		"  @SP",   // AレジスタにアドレスSPをセット
		"  A=M",   // AレジスタにSPの値をセット
		"  M=-1",  // スタックの先頭の値にtrueをセット
		"  @R15",  // AレジスタにアドレスR15をセット
		"  A=M",   // Aレジスタにリターンアドレスをセット
		"  0;JMP", // リターンアドレスにジャンプ
	}
}

func (ci *ConverterInitializer) initializeNEQ() []string {
	return []string{
		"(NEQ)",
		"  @SP",   // AレジスタにアドレスSPをセット
		"  A=M",   // AレジスタにSPの値をセット
		"  M=0",   // スタックの先頭の値にfalseをセット
		"  @R15",  // AレジスタにアドレスR15をセット
		"  A=M",   // Aレジスタにリターンアドレスをセット
		"  0;JMP", // リターンアドレスにジャンプ
	}
}
