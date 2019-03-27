package yaml_test

import (
	"bytes"
	"os"

	. "gopkg.in/check.v1"
	"gopkg.in/niemeyer/ynext.v3"
	"strings"
)

var nodeTests = []struct {
	yaml string
	node yaml.Node
}{
	{
		"null\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "null",
				Tag:    "!!null",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"foo\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "foo",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"\"foo\"\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Style:  yaml.DoubleQuotedStyle,
				Value:  "foo",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"'foo'\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Style:  yaml.SingleQuotedStyle,
				Value:  "foo",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!!str 123\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Style:  yaml.TaggedStyle,
				Value:  "123",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		// Although the node isn't TaggedStyle, dropping the tag would change the value.
		"[encode]!!binary gIGC\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "gIGC",
				Tag:    "!!binary",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		// Same, but with strings we can just quote them.
		"[encode]\"123\"\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "123",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!tag:something 123\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Style:  yaml.TaggedStyle,
				Value:  "123",
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"[encode]!tag:something 123\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "123",
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!tag:something {}\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Style:  yaml.TaggedStyle | yaml.FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"[encode]!tag:something {}\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Style:  yaml.FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"!tag:something []\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Style:  yaml.TaggedStyle | yaml.FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"[encode]!tag:something []\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Style:  yaml.FlowStyle,
				Tag:    "!tag:something",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"''\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Style:  yaml.SingleQuotedStyle,
				Value:  "",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"|\n  foo\n  bar\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Style:  yaml.LiteralStyle,
				Value:  "foo\nbar\n",
				Tag:    "!!str",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"true\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "true",
				Tag:    "!!bool",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"-10\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "-10",
				Tag:    "!!int",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"4294967296\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "4294967296",
				Tag:    "!!int",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"0.1000\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "0.1000",
				Tag:    "!!float",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"-.inf\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "-.inf",
				Tag:    "!!float",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		".nan\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  ".nan",
				Tag:    "!!float",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"{}\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Style:  yaml.FlowStyle,
				Value:  "",
				Tag:    "!!map",
				Line:   1,
				Column: 1,
			}},
		},
	}, {
		"a: b c\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Value:  "",
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 1,
				}, {
					Kind:   yaml.ScalarNode,
					Value:  "b c",
					Tag:    "!!str",
					Line:   1,
					Column: 4,
				}},
			}},
		},
	}, {
		"a:\n  b: c\n  d: e\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Tag:    "!!map",
				Line:   1,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 1,
				}, {
					Kind:   yaml.MappingNode,
					Tag:    "!!map",
					Line:   2,
					Column: 3,
					Children: []*yaml.Node{{
						Kind:   yaml.ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 3,
					}, {
						Kind:   yaml.ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   2,
						Column: 6,
					}, {
						Kind:   yaml.ScalarNode,
						Value:  "d",
						Tag:    "!!str",
						Line:   3,
						Column: 3,
					}, {
						Kind:   yaml.ScalarNode,
						Value:  "e",
						Tag:    "!!str",
						Line:   3,
						Column: 6,
					}},
				}},
			}},
		},
	}, {
		"- a\n- b\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Value:  "",
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 3,
				}, {
					Kind:   yaml.ScalarNode,
					Value:  "b",
					Tag:    "!!str",
					Line:   2,
					Column: 3,
				}},
			}},
		},
	}, {
		"- a\n- - b\n  - c\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 3,
				}, {
					Kind:   yaml.SequenceNode,
					Tag:    "!!seq",
					Line:   2,
					Column: 3,
					Children: []*yaml.Node{{
						Kind:   yaml.ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 5,
					}, {
						Kind:   yaml.ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   3,
						Column: 5,
					}},
				}},
			}},
		},
	}, {
		"[a, b]\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Style:  yaml.FlowStyle,
				Value:  "",
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 2,
				}, {
					Kind:   yaml.ScalarNode,
					Value:  "b",
					Tag:    "!!str",
					Line:   1,
					Column: 5,
				}},
			}},
		},
	}, {
		"- a\n- [b, c]\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Tag:    "!!seq",
				Line:   1,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 3,
				}, {
					Kind:   yaml.SequenceNode,
					Tag:    "!!seq",
					Style:  yaml.FlowStyle,
					Line:   2,
					Column: 3,
					Children: []*yaml.Node{{
						Kind:   yaml.ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 4,
					}, {
						Kind:   yaml.ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   2,
						Column: 7,
					}},
				}},
			}},
		},
	}, {
		"a: &x 1\nb: &y 2\nc: *x\nd: *y\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   1,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Line:   1,
				Column: 1,
				Tag:    "!!map",
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Value:  "a",
					Tag:    "!!str",
					Line:   1,
					Column: 1,
				},
					saveNode("x", &yaml.Node{
						Kind:   yaml.ScalarNode,
						Value:  "1",
						Tag:    "!!int",
						Anchor: "x",
						Line:   1,
						Column: 4,
					}),
					{
						Kind:   yaml.ScalarNode,
						Value:  "b",
						Tag:    "!!str",
						Line:   2,
						Column: 1,
					},
					saveNode("y", &yaml.Node{
						Kind:   yaml.ScalarNode,
						Value:  "2",
						Tag:    "!!int",
						Anchor: "y",
						Line:   2,
						Column: 4,
					}),
					{
						Kind:   yaml.ScalarNode,
						Value:  "c",
						Tag:    "!!str",
						Line:   3,
						Column: 1,
					}, {
						Kind:   yaml.AliasNode,
						Value:  "x",
						Alias:  dropNode("x"),
						Line:   3,
						Column: 4,
					}, {
						Kind:   yaml.ScalarNode,
						Value:  "d",
						Tag:    "!!str",
						Line:   4,
						Column: 1,
					}, {
						Kind:   yaml.AliasNode,
						Value:  "y",
						Tag:    "",
						Alias:  dropNode("y"),
						Line:   4,
						Column: 4,
					}},
			}},
		},
	}, {

		"# One\n# Two\ntrue # Three\n# Four\n# Five\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   3,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "true",
				Tag:    "!!bool",
				Line:   3,
				Column: 1,
				Header: "# One\n# Two",
				Inline: "# Three",
				Footer: "# Four\n# Five",
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# H1\n# H2\ntrue # I\n# F1\n# F2\n\n# DF1\n\n# DF2\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   7,
			Column: 1,
			Header: "# DH1\n\n# DH2",
			Footer: "# DF1\n\n# DF2",
			Children: []*yaml.Node{{
				Kind:   yaml.ScalarNode,
				Value:  "true",
				Tag:    "!!bool",
				Line:   7,
				Column: 1,
				Header: "# H1\n# H2",
				Inline: "# I",
				Footer: "# F1\n# F2",
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# HA1\n# HA2\nka: va # IA\n# FA1\n# FA2\n\n# HB1\n# HB2\nkb: vb # IB\n# FB1\n# FB2\n\n# DF1\n\n# DF2\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   7,
			Column: 1,
			Header: "# DH1\n\n# DH2",
			Footer: "# DF1\n\n# DF2",
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Tag:    "!!map",
				Line:   7,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Line:   7,
					Column: 1,
					Tag:    "!!str",
					Value:  "ka",
					Header: "# HA1\n# HA2",
					Footer: "# FA1\n# FA2",
				}, {
					Kind:   yaml.ScalarNode,
					Line:   7,
					Column: 5,
					Tag:    "!!str",
					Value:  "va",
					Inline: "# IA",
				}, {
					Kind:   yaml.ScalarNode,
					Line:   13,
					Column: 1,
					Tag:    "!!str",
					Value:  "kb",
					Header: "# HB1\n# HB2",
					Footer: "# FB1\n# FB2",
				}, {
					Kind:   yaml.ScalarNode,
					Line:   13,
					Column: 5,
					Tag:    "!!str",
					Value:  "vb",
					Inline: "# IB",
				}},
			}},
		},
	}, {
		"# DH1\n\n# DH2\n\n# HA1\n# HA2\n- la # IA\n# FA1\n# FA2\n\n# HB1\n# HB2\n- lb # IB\n# FB1\n# FB2\n\n# DF1\n\n# DF2\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   7,
			Column: 1,
			Header: "# DH1\n\n# DH2",
			Footer: "# DF1\n\n# DF2",
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Tag:    "!!seq",
				Line:   7,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   7,
					Column: 3,
					Value:  "la",
					Header: "# HA1\n# HA2",
					Inline: "# IA",
					Footer: "# FA1\n# FA2",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   13,
					Column: 3,
					Value:  "lb",
					Header: "# HB1\n# HB2",
					Inline: "# IB",
					Footer: "# FB1\n# FB2",
				}},
			}},
		},
	}, {
		"# DH1\n\n- la # IA\n\n# HB1\n- lb\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   3,
			Column: 1,
			Header: "# DH1",
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Tag:    "!!seq",
				Line:   3,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   3,
					Column: 3,
					Value:  "la",
					Inline: "# IA",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 3,
					Value:  "lb",
					Header: "# HB1",
				}},
			}},
		},
	}, {
		"# DH1\n\n# HA1\nka:\n  # HB1\n  kb:\n  # HC1\n  # HC2\n  - lc # IC\n  # FC1\n  # FC2\n\n  # HD1\n  - ld # ID\n  # FD1\n\n# DF1\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   4,
			Column: 1,
			Header: "# DH1",
			Footer: "# DF1",
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Tag:    "!!map",
				Line:   4,
				Column: 1,
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   4,
					Column: 1,
					Value:  "ka",
					Header: "# HA1",
				}, {
					Kind:   yaml.MappingNode,
					Tag:    "!!map",
					Line:   6,
					Column: 3,
					Children: []*yaml.Node{{
						Kind:   yaml.ScalarNode,
						Tag:    "!!str",
						Line:   6,
						Column: 3,
						Value:  "kb",
						Header: "# HB1",
					}, {
						Kind:   yaml.SequenceNode,
						Line:   9,
						Column: 3,
						Tag:    "!!seq",
						Children: []*yaml.Node{{
							Kind:   yaml.ScalarNode,
							Tag:    "!!str",
							Line:   9,
							Column: 5,
							Value:  "lc",
							Header: "# HC1\n# HC2",
							Inline: "# IC",
							Footer: "# FC1\n# FC2",
						}, {
							Kind:   yaml.ScalarNode,
							Tag:    "!!str",
							Line:   14,
							Column: 5,
							Value:  "ld",
							Header: "# HD1",

							Inline: "# ID",
							Footer: "# FD1",
						}},
					}},
				}},
			}},
		},
	}, {
		"# H1\n[la, lb] # I\n# F1\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   2,
			Column: 1,
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Tag:    "!!seq",
				Style:  yaml.FlowStyle,
				Line:   2,
				Column: 1,
				Header: "# H1",
				Inline: "# I",
				Footer: "# F1",
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   2,
					Column: 2,
					Value:  "la",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   2,
					Column: 6,
					Value:  "lb",
				}},
			}},
		},
	}, {
		"# DH1\n\n# SH1\n[\n  # HA1\n  la, # IA\n  # FA1\n\n  # HB1\n  lb, # IB\n  # FB1\n]\n# SF1\n\n# DF1\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   4,
			Column: 1,
			Header: "# DH1",
			Footer: "# DF1",
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Tag:    "!!seq",
				Style:  yaml.FlowStyle,
				Line:   4,
				Column: 1,
				Header: "# SH1",
				Footer: "# SF1",
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 3,
					Value:  "la",
					Header: "# HA1",
					Inline: "# IA",
					Footer: "# FA1",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   10,
					Column: 3,
					Value:  "lb",
					Header: "# HB1",
					Inline: "# IB",
					Footer: "# FB1",
				}},
			}},
		},
	}, {
		"# DH1\n\n# SH1\n[\n  # HA1\n  la,\n  # FA1\n\n  # HB1\n  lb,\n  # FB1\n]\n# SF1\n\n# DF1\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   4,
			Column: 1,
			Header: "# DH1",
			Footer: "# DF1",
			Children: []*yaml.Node{{
				Kind:   yaml.SequenceNode,
				Tag:    "!!seq",
				Style:  yaml.FlowStyle,
				Line:   4,
				Column: 1,
				Header: "# SH1",
				Footer: "# SF1",
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 3,
					Value:  "la",
					Header: "# HA1",
					Footer: "# FA1",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   10,
					Column: 3,
					Value:  "lb",
					Header: "# HB1",
					Footer: "# FB1",
				}},
			}},
		},
	}, {
		"# DH1\n\n# MH1\n{\n  # HA1\n  ka: va, # IA\n  # FA1\n\n  # HB1\n  kb: vb, # IB\n  # FB1\n}\n# MF1\n\n# DF1\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   4,
			Column: 1,
			Header: "# DH1",
			Footer: "# DF1",
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Tag:    "!!map",
				Style:  yaml.FlowStyle,
				Line:   4,
				Column: 1,
				Header: "# MH1",
				Footer: "# MF1",
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 3,
					Value:  "ka",
					Header: "# HA1",
					Footer: "# FA1",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 7,
					Value:  "va",
					Inline: "# IA",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   10,
					Column: 3,
					Value:  "kb",
					Header: "# HB1",
					Footer: "# FB1",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   10,
					Column: 7,
					Value:  "vb",
					Inline: "# IB",
				}},
			}},
		},
	}, {
		"# DH1\n\n# MH1\n{\n  # HA1\n  ka: va,\n  # FA1\n\n  # HB1\n  kb: vb,\n  # FB1\n}\n# MF1\n\n# DF1\n",
		yaml.Node{
			Kind:   yaml.DocumentNode,
			Line:   4,
			Column: 1,
			Header: "# DH1",
			Footer: "# DF1",
			Children: []*yaml.Node{{
				Kind:   yaml.MappingNode,
				Tag:    "!!map",
				Style:  yaml.FlowStyle,
				Line:   4,
				Column: 1,
				Header: "# MH1",
				Footer: "# MF1",
				Children: []*yaml.Node{{
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 3,
					Value:  "ka",
					Header: "# HA1",
					Footer: "# FA1",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   6,
					Column: 7,
					Value:  "va",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   10,
					Column: 3,
					Value:  "kb",
					Header: "# HB1",
					Footer: "# FB1",
				}, {
					Kind:   yaml.ScalarNode,
					Tag:    "!!str",
					Line:   10,
					Column: 7,
					Value:  "vb",
				}},
			}},
		},
	},
}

func (s *S) TestNodeRoundtrip(c *C) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")
	for i, item := range nodeTests {
		c.Logf("test %d: %q", i, item.yaml)

		decode := true
		encode := true

		testYaml := item.yaml
		if s := strings.TrimPrefix(testYaml, "[decode]"); s != testYaml {
			encode = false
			testYaml = s
		}
		if s := strings.TrimPrefix(testYaml, "[encode]"); s != testYaml {
			decode = false
			testYaml = s
		}

		if decode {
			var node yaml.Node
			err := yaml.Unmarshal([]byte(testYaml), &node)
			c.Assert(err, IsNil)
			c.Assert(node, DeepEquals, item.node)
		}
		if encode {
			buf := bytes.Buffer{}
			enc := yaml.NewEncoder(&buf)
			enc.SetIndent(2)
			err := enc.Encode(&item.node)
			c.Assert(err, IsNil)
			err = enc.Close()
			c.Assert(err, IsNil)
			c.Assert(buf.String(), Equals, testYaml)
		}
	}
}

var savedNodes = make(map[string]*yaml.Node)

func saveNode(name string, node *yaml.Node) *yaml.Node {
	savedNodes[name] = node
	return node
}

func peekNode(name string) *yaml.Node {
	return savedNodes[name]
}

func dropNode(name string) *yaml.Node {
	node := savedNodes[name]
	delete(savedNodes, name)
	return node
}