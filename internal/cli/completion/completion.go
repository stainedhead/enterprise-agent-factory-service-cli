package completion

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

func init() {
	registrations = append(registrations, func(r *cli.Registry) {
		r.AddLocal(&cli.Command{Service: "completion", Summary: "Print a shell completion script for bash, zsh, fish or PowerShell.", Local: run,
			Override: cli.Override{Examples: []cli.Example{
				{Title: "Enable completion in the current bash session", Command: "source <(afs completion bash)"},
				{Title: "Install the zsh script", Command: "afs completion zsh > \"${fpath[1]}/_afs\""},
				{Title: "Install the fish script", Command: "afs completion fish > ~/.config/fish/completions/afs.fish"},
				{Title: "Load the PowerShell script in a profile", Command: "afs completion powershell | Out-String | Invoke-Expression"},
			}}})
	})
}

// Shells are the shells a script exists for.
var Shells = []string{"bash", "zsh", "fish", "powershell"}

func run(_ context.Context, inv *cli.Invocation) int {
	if len(inv.P.Positional) != 1 {
		_, _ = fmt.Fprintf(inv.Env.Stderr, "usage: afs completion <%s>\n\nafs: error: give the shell to print a completion script for\n", strings.Join(Shells, "|"))
		return cli.ExitUsage
	}
	if err := Generate(inv.Env.Stdout, inv.P.Positional[0], inv.Reg); err != nil {
		_, _ = fmt.Fprintf(inv.Env.Stderr, "usage: afs completion <%s>\n\nafs: error: %v\n", strings.Join(Shells, "|"), err)
		return cli.ExitUsage
	}
	return cli.ExitOK
}

// model is everything a script needs, read from the command table (D-C27).
type model struct {
	globals    []string            // every global option, with dashes
	globalEnum map[string][]string // --output, --color
	takesValue []string            // global options that take a value
	services   []service
}

type service struct {
	name string
	ops  []operation
}

type operation struct {
	name string
	opts []option
}

type option struct {
	flag string
	enum []string
}

func build(r *cli.Registry) model {
	m := model{globalEnum: map[string][]string{"--output": cli.OutputFormats, "--color": cli.ColorModes}}
	for _, g := range cli.GlobalHelp() {
		m.globals = append(m.globals, g[0])
		switch g[0] {
		case "--no-paginate", "--generate-cli-skeleton", "--debug", "--no-cli-pager", "--version":
		default:
			m.takesValue = append(m.takesValue, g[0])
		}
	}
	for _, s := range r.Services() {
		sv := service{name: s}
		for _, c := range r.Commands(s) {
			if c.Name == "" {
				continue
			}
			op := operation{name: c.Name}
			for _, o := range c.Options {
				op.opts = append(op.opts, option{flag: o.Flag, enum: o.Enum})
				if o.Kind == cli.Bool {
					op.opts = append(op.opts, option{flag: "--no-" + o.Name()})
				}
			}
			sort.Slice(op.opts, func(i, j int) bool { return op.opts[i].flag < op.opts[j].flag })
			sv.ops = append(sv.ops, op)
		}
		if s == "completion" {
			for _, sh := range Shells {
				sv.ops = append(sv.ops, operation{name: sh})
			}
		}
		m.services = append(m.services, sv)
	}
	return m
}

func (o operation) flags() []string {
	out := make([]string, len(o.opts))
	for i, x := range o.opts {
		out[i] = x.flag
	}
	return out
}

// Generate writes the completion script of a shell for the registry's command table.
func Generate(w io.Writer, shell string, r *cli.Registry) error {
	m := build(r)
	var s string
	switch shell {
	case "bash":
		s = m.bash()
	case "zsh":
		s = m.zsh()
	case "fish":
		s = m.fish()
	case "powershell":
		s = m.powershell()
	default:
		return fmt.Errorf("unknown shell %q: choose from %s", shell, strings.Join(Shells, ", "))
	}
	_, err := io.WriteString(w, s)
	return err
}

func words(l []string) string { return strings.Join(l, " ") }

func (m model) serviceNames() []string {
	out := make([]string, len(m.services))
	for i, s := range m.services {
		out[i] = s.name
	}
	return out
}

func (s service) opNames() []string {
	out := make([]string, len(s.ops))
	for i, o := range s.ops {
		out[i] = o.name
	}
	return out
}

func (m model) bash() string {
	var b strings.Builder
	b.WriteString("# bash completion for afs; load with: source <(afs completion bash)\n_afs() {\n")
	b.WriteString("  local cur prev svc op i w\n  COMPREPLY=()\n  cur=\"${COMP_WORDS[COMP_CWORD]}\"\n  prev=\"${COMP_WORDS[COMP_CWORD-1]}\"\n  svc=\"\"; op=\"\"\n")
	b.WriteString("  for ((i=1; i<COMP_CWORD; i++)); do\n    w=\"${COMP_WORDS[i]}\"\n    case \"$w\" in\n      -*) ;;\n")
	b.WriteString("      *) case \"${COMP_WORDS[i-1]}\" in\n           " + strings.Join(m.takesValue, "|") + ") ;;\n           *) if [ -z \"$svc\" ]; then svc=\"$w\"; elif [ -z \"$op\" ]; then op=\"$w\"; fi ;;\n         esac ;;\n    esac\n  done\n")
	b.WriteString("  case \"$svc $op $prev\" in\n")
	for _, s := range m.services {
		for _, o := range s.ops {
			for _, x := range o.opts {
				if len(x.enum) > 0 {
					fmt.Fprintf(&b, "    \"%s %s %s\") COMPREPLY=($(compgen -W \"%s\" -- \"$cur\")); return ;;\n", s.name, o.name, x.flag, words(x.enum))
				}
			}
		}
	}
	for _, f := range []string{"--output", "--color"} {
		fmt.Fprintf(&b, "    *\" %s\") COMPREPLY=($(compgen -W \"%s\" -- \"$cur\")); return ;;\n", f, words(m.globalEnum[f]))
	}
	b.WriteString("  esac\n  if [[ \"$cur\" == -* ]]; then\n    case \"$svc $op\" in\n")
	for _, s := range m.services {
		for _, o := range s.ops {
			fmt.Fprintf(&b, "      \"%s %s\") COMPREPLY=($(compgen -W \"%s %s\" -- \"$cur\")) ;;\n", s.name, o.name, words(o.flags()), words(m.globals))
		}
	}
	fmt.Fprintf(&b, "      *) COMPREPLY=($(compgen -W \"%s\" -- \"$cur\")) ;;\n    esac\n    return\n  fi\n", words(m.globals))
	b.WriteString("  if [ -z \"$svc\" ]; then\n    COMPREPLY=($(compgen -W \"" + words(m.serviceNames()) + "\" -- \"$cur\"))\n  elif [ -z \"$op\" ]; then\n    case \"$svc\" in\n")
	for _, s := range m.services {
		fmt.Fprintf(&b, "      %s) COMPREPLY=($(compgen -W \"%s\" -- \"$cur\")) ;;\n", s.name, words(s.opNames()))
	}
	b.WriteString("    esac\n  fi\n}\ncomplete -F _afs afs\n")
	return b.String()
}

func (m model) zsh() string {
	var b strings.Builder
	b.WriteString("#compdef afs\n# zsh completion for afs; save as _afs in a directory of $fpath\n_afs() {\n  local svc op prev\n  svc=\"${words[2]}\"; op=\"${words[3]}\"; prev=\"${words[CURRENT-1]}\"\n")
	b.WriteString("  case \"$svc $op $prev\" in\n")
	for _, s := range m.services {
		for _, o := range s.ops {
			for _, x := range o.opts {
				if len(x.enum) > 0 {
					fmt.Fprintf(&b, "    \"%s %s %s\") compadd -- %s; return ;;\n", s.name, o.name, x.flag, words(x.enum))
				}
			}
		}
	}
	for _, f := range []string{"--output", "--color"} {
		fmt.Fprintf(&b, "    *\" %s\") compadd -- %s; return ;;\n", f, words(m.globalEnum[f]))
	}
	b.WriteString("  esac\n  if [[ \"${words[CURRENT]}\" == -* ]]; then\n    case \"$svc $op\" in\n")
	for _, s := range m.services {
		for _, o := range s.ops {
			fmt.Fprintf(&b, "      \"%s %s\") compadd -- %s %s ;;\n", s.name, o.name, words(o.flags()), words(m.globals))
		}
	}
	fmt.Fprintf(&b, "      *) compadd -- %s ;;\n    esac\n    return\n  fi\n", words(m.globals))
	b.WriteString("  if (( CURRENT == 2 )); then\n    compadd -- " + words(m.serviceNames()) + "\n  elif (( CURRENT == 3 )); then\n    case \"$svc\" in\n")
	for _, s := range m.services {
		fmt.Fprintf(&b, "      %s) compadd -- %s ;;\n", s.name, words(s.opNames()))
	}
	b.WriteString("    esac\n  fi\n}\ncompdef _afs afs\n")
	return b.String()
}

func (m model) fish() string {
	var b strings.Builder
	b.WriteString("# fish completion for afs; save as ~/.config/fish/completions/afs.fish\ncomplete -c afs -f\n")
	for _, g := range m.globals {
		if e, ok := m.globalEnum[g]; ok {
			fmt.Fprintf(&b, "complete -c afs -l %s -xa '%s'\n", strings.TrimPrefix(g, "--"), words(e))
		} else {
			fmt.Fprintf(&b, "complete -c afs -l %s\n", strings.TrimPrefix(g, "--"))
		}
	}
	fmt.Fprintf(&b, "complete -c afs -n '__fish_use_subcommand' -a '%s'\n", words(m.serviceNames()))
	for _, s := range m.services {
		fmt.Fprintf(&b, "complete -c afs -n '__fish_seen_subcommand_from %s; and not __fish_seen_subcommand_from %s' -a '%s'\n", s.name, words(s.opNames()), words(s.opNames()))
		for _, o := range s.ops {
			for _, x := range o.opts {
				cond := fmt.Sprintf("__fish_seen_subcommand_from %s; and __fish_seen_subcommand_from %s", s.name, o.name)
				if len(x.enum) > 0 {
					fmt.Fprintf(&b, "complete -c afs -n '%s' -l %s -xa '%s'\n", cond, strings.TrimPrefix(x.flag, "--"), words(x.enum))
				} else {
					fmt.Fprintf(&b, "complete -c afs -n '%s' -l %s\n", cond, strings.TrimPrefix(x.flag, "--"))
				}
			}
		}
	}
	return b.String()
}

func quoteList(l []string) string {
	q := make([]string, len(l))
	for i, s := range l {
		q[i] = "'" + s + "'"
	}
	return strings.Join(q, ",")
}

func (m model) powershell() string {
	var b strings.Builder
	b.WriteString("# PowerShell completion for afs; load with: afs completion powershell | Out-String | Invoke-Expression\n")
	fmt.Fprintf(&b, "$afsValueFlags = @(%s)\n$afsGlobals = @(%s)\n$afsEnums = @{ '--output' = @(%s); '--color' = @(%s) }\n", quoteList(m.takesValue), quoteList(m.globals), quoteList(m.globalEnum["--output"]), quoteList(m.globalEnum["--color"]))
	b.WriteString("$afsCommands = @{\n")
	for _, s := range m.services {
		fmt.Fprintf(&b, "  '%s' = @{\n", s.name)
		for _, o := range s.ops {
			fmt.Fprintf(&b, "    '%s' = @(%s)\n", o.name, quoteList(o.flags()))
		}
		b.WriteString("  }\n")
	}
	b.WriteString("}\n$afsOptionEnums = @{\n")
	for _, s := range m.services {
		for _, o := range s.ops {
			for _, x := range o.opts {
				if len(x.enum) > 0 {
					fmt.Fprintf(&b, "  '%s %s %s' = @(%s)\n", s.name, o.name, x.flag, quoteList(x.enum))
				}
			}
		}
	}
	b.WriteString("}\n")
	b.WriteString(`Register-ArgumentCompleter -Native -CommandName afs -ScriptBlock {
  param($wordToComplete, $commandAst, $cursorPosition)
  $w = @($commandAst.CommandElements | ForEach-Object { $_.ToString() })
  $positional = @()
  for ($i = 1; $i -lt $w.Count; $i++) {
    if ($w[$i] -notlike '-*' -and -not ($i -gt 1 -and $w[$i-1] -in $afsValueFlags)) { $positional += $w[$i] }
  }
  if ($wordToComplete -ne '' -and $positional.Count -gt 0) { $positional = @($positional | Select-Object -SkipLast 1) }
  $prev = if ($wordToComplete -eq '') { $w[-1] } elseif ($w.Count -gt 1) { $w[-2] } else { '' }
  $candidates = @()
  if ($positional.Count -ge 2 -and $afsOptionEnums.ContainsKey("$($positional[0]) $($positional[1]) $prev")) {
    $candidates = $afsOptionEnums["$($positional[0]) $($positional[1]) $prev"]
  } elseif ($afsEnums.ContainsKey($prev)) {
    $candidates = $afsEnums[$prev]
  } elseif ($wordToComplete -like '-*') {
    $candidates = $afsGlobals
    if ($positional.Count -ge 2 -and $afsCommands.ContainsKey($positional[0]) -and $afsCommands[$positional[0]].ContainsKey($positional[1])) {
      $candidates = $afsCommands[$positional[0]][$positional[1]] + $afsGlobals
    }
  } elseif ($positional.Count -eq 0) {
    $candidates = $afsCommands.Keys
  } elseif ($positional.Count -eq 1 -and $afsCommands.ContainsKey($positional[0])) {
    $candidates = $afsCommands[$positional[0]].Keys
  }
  $candidates | Where-Object { $_ -like "$wordToComplete*" } | Sort-Object | ForEach-Object {
    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
  }
}
`)
	return b.String()
}
