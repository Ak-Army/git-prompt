# go-git-prompt

Informative and fast Git prompt for any shell (Bash, Zsh, and PowerShell). 
Inspired by zsh-git-prompt and ubnt-intrepid/go-git-prompt

## Usage
Bash:
```bash
PS1='\w \$(git-prompt) % '
```

Zsh:
```zsh
PROMPT='%~ $(git-prompt) %% '
```

Fish:
```fish
function fish_prompt
   echo (git-prompt)" % "
end
```

PowerShell:
```ps1
function prompt {
  write-host "$(pwd) " -nonewline
  write-host (git-prompt) -nonewline
  return "`n> "
}
```

## Install
```shell-session
$ go get -v github.com/Ak-Army/git-prompt
```

## License
This software is released under the MIT license.
See [LICENSE](LICENSE) for details.
