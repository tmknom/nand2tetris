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
	return []string{
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M-1", // スタック領域の先頭アドレスをデクリメント
		"A=M",   // Aレジスタにスタック領域の先頭アドレスをセット
		"D=M",   // スタック領域の先頭の値をDレジスタにセット
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M-1", // スタック領域の先頭アドレスをデクリメント
		"A=M",   // Aレジスタにスタック領域の先頭アドレスをセット
		"M=D+M", // スタック領域の先頭の値とDレジスタの値を加算
		"@SP",   // AレジスタにアドレスSPをセット
		"M=M+1", // スタック領域の先頭アドレスをデクリメント
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
	end := ci.initializeEND()
	eq := ci.initializeEQ()
	neq := ci.initializeNEQ()

	result := []string{}
	result = append(result, end...)
	result = append(result, eq...)
	result = append(result, neq...)

	return result
}

func (ci *ConverterInitializer) initializeEND() []string {
	return []string{
		"@END",
		"0;JMP",
		"(END)",
		"  @END",  // AレジスタにENDラベルをセット
		"  0;JMP", // ENDラベルにジャンプ
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
