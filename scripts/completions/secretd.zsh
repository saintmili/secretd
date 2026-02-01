#compdef secretd

_arguments -C \
  '1:command:(
    init
    unlock
    add
    list
    show
    update
    delete
    change-master-password
    generate
    export
    doctor
    help
  )' \
  '*::arg:->args'

case $state in
args)
  case $words[2] in
    add|update)
      _arguments \
        '--generate[generate password]:length' \
        '--reveal[show password]'
      ;;
    show)
      _arguments \
        '--reveal[show password]'
      ;;
    export)
      _values 'format' json csv
      ;;
  esac
  ;;
esac

