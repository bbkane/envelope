# https://github.com/charmbracelet/vhs

Require enventory

Output demo.gif

Set Shell zsh
Set PlaybackSpeed 0.5
Set TypingSpeed 75ms

Set FontSize 22
Set Height 650
Set Width 1250
Set Padding 25

Type 'eval "$(enventory shell zsh init)"  # add this to ~/.zshrc'
Sleep 200ms
Enter
Sleep 1s

Type "enventory env create --name demo_env"
Sleep 200ms
Enter
Sleep 1s

Type "enventory var create --env demo_env --name MYENVVAR --value 'hello enventory'"
Sleep 200ms
Enter
Sleep 1s

Type "enventory env show --name demo_env"
Sleep 200ms
Enter
Sleep 1s

Type "export-env demo_env"
Sleep 200ms
Enter
Sleep 1s


Type "echo $MYENVVAR"
Sleep 200ms
Enter
Sleep 1s

Sleep 3s
