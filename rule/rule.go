package rule

import _ "embed"

//go:embed code_rule.tmpl
var CodeRuleTemplate string

//go:embed ask_rule.tmpl
var AskRuleTemplate string
