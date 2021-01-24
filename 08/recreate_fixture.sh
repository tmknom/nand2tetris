#/bin/bash
#
# テスト用の比較ファイルを一括更新する
#
# 共通部分の更新時などに利用することを想定
# CPUEmulatorでテストも実行し、コミットしても問題ないか同時にチェックできる

set -ex

function recreate(){
  dir_name=${1}
  file_name="${dir_name}/${2}"

  # 新しいasmファイル作成
  go run . ${file_name}

  pushd ${dir_name} > /dev/null
  cmp=$(find . -name "*.asm.cmp") # cmpファイルの検索
  src=$(echo ${cmp} | sed -e 's/.cmp//g')
  cp ${src} ${cmp}
  popd > /dev/null
}

function test(){
  dir_name=${1}

  pushd ${dir_name} > /dev/null
  tst=$(find . -name "*.tst" | grep -v VME)
  tool_path=$(git rev-parse --show-toplevel)/suite/tools
  ${tool_path}/CPUEmulator.sh ${tst}
  popd > /dev/null
}

function recreateAndTest(){
  dir_name=${1}
  file_name=${2}
  recreate ${dir_name} ${file_name}
  test ${dir_name}
}

#recreateAndTest "FunctionCalls/StaticsTest"
#recreateAndTest "FunctionCalls/NestedCall"
recreateAndTest "FunctionCalls/FibonacciElement"

# SimpleFunctionのみテスト時の初期化処理が特殊で、テスト実行前に
# SimpleFunction.asmの 「(SimpleFunction.test)」から手前の初期化コードを事前に削除が必要
# sedコマンドで初期化コードを削除してテストを実行する
recreate "FunctionCalls/SimpleFunction" "SimpleFunction.vm"
sed -i '' "1,20d" "FunctionCalls/SimpleFunction/SimpleFunction.asm"
test "FunctionCalls/SimpleFunction"

recreateAndTest "ProgramFlow/FibonacciSeries" "FibonacciSeries.vm"
recreateAndTest "ProgramFlow/BasicLoop" "BasicLoop.vm"

recreateAndTest "MemoryAccess/StaticTest" "StaticTest.vm"
recreateAndTest "MemoryAccess/PointerTest" "PointerTest.vm"
recreateAndTest "MemoryAccess/BasicTest" "BasicTest.vm"

recreateAndTest "StackArithmetic/StackTest" "StackTest.vm"
recreateAndTest "StackArithmetic/SimpleAdd" "SimpleAdd.vm"

echo "All succeeded!"
