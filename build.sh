_help() {
  cat <<EOT
Usage: ./build.sh <command>

commands:
    test              run the tests (locally)
    test.e, teste     run the tests (containarized using earthly)
    build             build the artifacts (locally)
    build.e, builde   build the artifacts (contanarized using earthly)
    help              print this message
EOT
}

xx() {
  echo ">>" ${*}
  ${*}
}

case "${1}" in
  test) shift; xx go test ./... ${*} ;;
  test.e|teste) shift; xx earthly +test ${*} ;;
  build) shift; xx go build -o targets/ ${*} ;;
  build.e|builde) shift; xx earthly +build ${*} ;;
  *|help) _help ;;
esac
