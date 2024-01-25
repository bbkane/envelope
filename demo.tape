# https://github.com/charmbracelet/vhs

Require namedenv

Output demo.gif

Set Shell zsh
Set PlaybackSpeed 0.5
Set TypingSpeed 75ms

Set FontSize 22
Set Height 650
Set Width 1650
Set Padding 50

# https://github.com/charmbracelet/vhs/issues/39#issuecomment-1296334921
Hide
Sleep 3s
Show

Type "namedenv env create --name demo_env"
Sleep 200ms
Enter
Sleep 1s

Type "namedenv env localvar create  --env-name demo_env --name MYENVVAR --value 'hello namedenv'"
Sleep 200ms
Enter
Sleep 1s

Type "namedenv env show --name demo_env"
Sleep 200ms
Enter
Sleep 1s

Type "source <(namedenv env print-script --name demo_env)"
Sleep 200ms
Enter
Sleep 1s

Type "echo $MYENVVAR"
Sleep 200ms
Enter
Sleep 1s

Sleep 3s
