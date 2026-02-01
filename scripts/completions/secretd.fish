# fish completion for secretd

complete -c secretd -n '__fish_use_subcommand' -a '
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
'

# add / update flags
complete -c secretd -n '__fish_seen_subcommand_from add update' -l generate -d 'Generate password'
complete -c secretd -n '__fish_seen_subcommand_from add show' -l reveal -d 'Reveal password'

# export formats
complete -c secretd -n '__fish_seen_subcommand_from export' -a 'json csv'

