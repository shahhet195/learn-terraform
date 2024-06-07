| Module | Description |
|------|------|
{{- range $key, $mods := .ByFirstLevel }}
| **{{ $key | trimPrefix "modules/" | trimSuffix "test" | title }}** | |
{{- range $i, $mod := $mods }}
| [{{ $mod.RelDirectory | trimPrefix $key | trimPrefix "/" }}]({{ $mod.RelDirectory }}) | {{ $mod.ShortDescription | trimAll " " | abbrev 80 |  }} |
{{- end }}
{{- end }}
