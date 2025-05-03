_help() {
  cat <<EOT
Usage: ./build.sh <command>

commands:
    test              run the tests (containarized using earthly)
    test.l, testl     run the tests (locally)
    build             build the artifacts (contanarized using earthly)
    build.l, buildl   build the artifacts (locally)
    help              print this message
EOT
}

xx() {
  echo ">>" ${*}
  ${*}
}

case "${1}" in
  test) shift; xx earthly +test ${*} ;;
  test.l|testl) shift; xx go test ./... ${*} ;;
  build) shift; xx earthly +build ${*} ;;
  build.l|buildl) shift; xx go build -o targets/ ${*} ;;
  *|help) _help ;;
esac
