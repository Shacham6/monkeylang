_help() {
  cat <<EOT
Usage: ./build.sh <command>

commands:
    test              run the tests (locally)
    test.d, testd     run the tests (containarized using raw docker)
    test.e, teste     run the tests (containarized using earthly)
    build             build the artifacts (contanarized using earthly)
    build.d, buildd   build the artifacts (contanarized using raw docker)
    build.l, buildl   build the artifacts (locally)
    help              print this message
EOT
}

xx() {
  echo ">>" ${*}
  ${*}
}

case "${1}" in
  test) shift; xx go test ./... ${*} ;;
  testd|test.d) shift; xx docker build . --target=test ${*} ;;
  test.e|teste) shift; xx earthly +test ${*} ;;
  build) shift; xx go build -o targets/ ${*} ;;
  build.d|buildd) shift; xx docker build -v $(pwd)/targets:/app/targets --target build . ;;
  build.e|builde) shift; xx earthly +build ${*} ;;
  *|help) _help ;;
esac
