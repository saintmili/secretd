# bash completion for secretd

_secretd() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    commands="init unlock add list show update delete change-master-password generate export doctor help"

    case "$COMP_CWORD" in
        1)
            COMPREPLY=( $(compgen -W "$commands" -- "$cur") )
            return 0
            ;;
    esac

    case "${COMP_WORDS[1]}" in
        add|update)
            if [[ "$cur" == --* ]]; then
                COMPREPLY=( $(compgen -W "--generate --reveal" -- "$cur") )
            fi
            ;;
        show)
            if [[ "$cur" == --* ]]; then
                COMPREPLY=( $(compgen -W "--reveal" -- "$cur") )
            fi
            ;;
        export)
            COMPREPLY=( $(compgen -W "json csv" -- "$cur") )
            ;;
    esac
}

complete -F _secretd secretd

