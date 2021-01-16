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

(FILL)
    @SCREEN // AレジスタにSCREENシンボルをセット
    M=1 // RAM[SCREEN]に1をセットして黒く塗りつぶす

    // 処理が完了したら、LOOPに戻る
    @LOOP // AレジスタにLOOPラベルをセット
    0;JMP // LOOPへ移動

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
