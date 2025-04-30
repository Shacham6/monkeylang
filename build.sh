case "${1}" in
  test) earthly +test ;;
  test.l) make test ;;
  build) earthly +build ;;
esac
