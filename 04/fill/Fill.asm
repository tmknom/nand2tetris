// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// RAM[0]に初期値「32767」をセット
@32767 // Aレジスタに定数「32767」をセット
D=A // Aレジスタの値(=32767)をDレジスタに格納
@0 // Aレジスタにアドレス0をセット
M=D // R0=RAM[0]にDレジスタの値をセット

// RAM[1]に初期値「8192」をセット
@8192 // Aレジスタに定数「8192」をセット
D=A // Aレジスタの値(=8192)をDレジスタに格納
@1 // Aレジスタにアドレス1をセット
M=D // R1=RAM[1]にDレジスタの値をセット

// RAM[10]に初期値「0」をセット
@10 // Aレジスタにアドレス10をセット
M=0 // RAM[10]に0をセット

// 初期値をセットしたらLOOPへ移動
@LOOP // AレジスタにLOOPラベルをセット
0;JMP // LOOPへ移動

(FILL)
    // FILLの終了条件のチェック
    @10 // Aレジスタにアドレス10をセット
    D=M // DレジスタにR10の値を格納
    @R1 // AレジスタにR1シンボルをセット
    D=M-D // R1の値からDレジスタの値(R10)を減算して、結果をDレジスタに格納
    @LOOP // AレジスタにLOOPラベルをセット
    D;JLE // Dレジスタの値がゼロ以下ならLOOPへ移動

    // 塗りつぶす対象のアドレスを計算
    @10 // Aレジスタにアドレス10をセット
    D=M // RAM[10]の値をDレジスタに格納
    @SCREEN // AレジスタにSCREENのアドレスをセット
    D=D+A // Dレジスタの値を「RAM[10]の値+SCREENのアドレス」に変更(=塗りつぶす先のアドレス)
    @11 // Aレジスタにアドレス11をセット
    M=D // RAM[11]にDレジスタの値(=塗りつぶす先のアドレス)をセット

    // 塗りつぶす
    @R0 // AレジスタにR0シンボルをセット
    D=M // DレジスタにR0の値(=32767)をセット
    @11 // Aレジスタにアドレス11をセット
    A=M // Aレジスタに塗りつぶす先のアドレスをセット
    M=D // 塗りつぶす先のアドレスに、Dレジスタの値をセット

    // RAM[10]をインクリメント
    @10 // Aレジスタにアドレス10をセット
    M=M+1

    // 処理が完了したら、FILLに戻る
    @FILL // AレジスタにFILLラベルをセット
    0;JMP // FILLへ移動

(CLEAR)
    @SCREEN // AレジスタにSCREENシンボルをセット
    M=0 // RAM[SCREEN]に0をセットしてクリア

    // 処理が完了したら、LOOPに戻る
    @LOOP // AレジスタにLOOPラベルをセット
    0;JMP // LOOPへ移動

(LOOP)
    // キーボードの値をDレジスタにセット
    @KBD // AレジスタにKBDシンボルをセット
    D=M // DレジスタにRAM[KBD]の値をセット

    // キーボード入力があればFILLへジャンプ
    @FILL // AレジスタにFILLラベルをセット
    D;JNE // Dレジスタがゼロ以外ならFILLへジャンプ

    // キーボード入力がなければCLEARへジャンプ
    @CLEAR // AレジスタにCLEARラベルをセット
    0;JMP // 問答無用でCLEARへジャンプ

(END)
    @END
    0;JMP
