// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)

//
// 初期値の設定
//

// R0(=RAM[0])に初期値「3」をセット
@3 // Aレジスタに定数「3」をセット
D=A // Aレジスタの値(=3)をDレジスタに格納
@0 // Aレジスタにアドレス0をセット
M=D // RAM[A]=RAM[0]にDレジスタの値をセット

// R1に初期値「2」をセット
@2 // Aレジスタに定数「2」をセット
D=A // Aレジスタ(=2)の値をDレジスタに格納
@1 // Aレジスタにアドレス1をセット
M=D // RAM[A]=RAM[1]にDレジスタの値をセット

// R2に初期値「-1」をセット
@1 // Aレジスタに定数「1」をセット
D=-A // Aレジスタ(=-1)の値をDレジスタに格納
@2 // Aレジスタにアドレス2をセット
M=D // RAM[2]にDレジスタの値をセット

//
// 初期化処理
//

// R2の初期化
@2 // Aレジスタにアドレス2をセット
D=M // DレジスタにRAM[2]の値をセット
@INIT // AレジスタにINITラベルをセット
D;JLT // if D<0 goto INIT

(INIT)
    // 外から「-1」がセットされたらRAM[2]をゼロにクリアする
    @0 // Aレジスタに定数「0」をセット
    D=A // Aレジスタ(=0)の値をDレジスタに格納
    @2 // Aレジスタにアドレス2をセット
    M=D // RAM[2]にDレジスタの値をセット

//
// 計算処理
//

(LOOP)
    // while文の終了条件をチェック
    // R1の値をDレジスタにロードし、Dレジスタの値が0以下ならENDへ
    @1 // Aレジスタにアドレス1をセット
    D=M // DレジスタにRAM[A]=RAM[1]の値をセット
    @END // AレジスタにENDラベルをセット
    D;JLE // if D=<0 goto END

    // R2 = R2 + R0を計算
    @0 // Aレジスタにアドレス0をセット
    D=M // DレジスタにRAM[0]の値をセット
    @2 // Aレジスタにアドレス2をセット
    M=D+M // RAM[0](=Dレジスタ)とM(RAM[2])の値を加算して、M(RAM[2])へ保存

    // R1をデクリメント
    @1 // Aレジスタにアドレス1をセット
    M=M-1 // RAM[1]の値をデクリメント

    // LOOP内の処理がワンセット完了したら、LOOPの先頭に戻る
    @LOOP // AレジスタにLOOPラベルをセット
    0;JMP // LOOPへ移動

(END)
    @END
    0;JMP
