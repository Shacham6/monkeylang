_help() {
  cat <<EOT
Usage: ./build.sh <command>

commands:
    test              run the tests (locally)
    test.d, testd     run the tests (containarized using raw docker)
    test.e, teste     run the tests (containarized using earthly)
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
  test) shift; xx go test ./... ${*} ;;
  testd|test.d) shift; xx docker build . --target=test ${*} ;;
  test.e|teste) shift; xx earthly +test ${*} ;;
  build) shift; xx earthly +build ${*} ;;
  build.l|buildl) shift; xx go build -o targets/ ${*} ;;
  *|help) _help ;;
esac
