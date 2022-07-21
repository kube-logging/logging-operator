package config

import (
	"strings"

	"github.com/siliconbrain/go-seqs/seqs"
)

func versionStmt(version string) Renderer {
	return Line(Formatted("@version: %s", version))
}

func includeStmt(file string) Renderer {
	return Line(Formatted("@include %q", file))
}

func globalOptionsDefStmt(def globalOptionsDef) Renderer {
	return braceDefStmt("options", "", nil) // TODO
}

type globalOptionsDef struct {
	// TODO
}

func flagsOption(flags []string) Renderer {
	return optionExpr("flags", flags...)
}

func optionExpr[V any](key string, values ...V) Renderer {
	return AllOf(String(key), String("("), SpaceSeparated(seqs.ToSlice(seqs.Map(seqs.Map(seqs.FromSlice(values), eraseType[V]), literal))...), String(")"))
}

func braceDefStmt(kind string, name string, body Renderer) Renderer {
	return AllOf(
		Line(AllOf(String(kind), Space, If(name != "", AllOf(String(name), Space)), String("{"))),
		Indented(body),
		Line(String("};")),
	)
}

func parenDefStmt(kind string, args ...Renderer) Renderer {
	return Line(AllOf(String(kind), String("("), SpaceSeparated(args...), String(");")))
}

func literal(v any) Renderer {
	switch v := v.(type) {
	case string:
		if strings.ContainsAny(v, " ${}") {
			return Quoted(v)
		}
		return String(v)
	case bool:
		if v {
			return String("yes")
		}
		return String("no")
	default:
		return Formatted("%v", v)
	}
}

func eraseType[V any](v V) any {
	return v
}
