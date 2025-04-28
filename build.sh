case "${1}" in
  test) earthly +test ;;

  build) earthly +build ;;
esac
