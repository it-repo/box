package ac

type (
	// Type of AC
	Type uint8
	// Logical of AC
	Logical uint8
)

// enum
const (
	TypeRole = iota
	TypePerm

	LogicalOR = iota
	LogicalAND
)

// Rule for Access Control
type Rule struct {
	rules []string
	both  bool
}

// NewRule Builder
func NewRule(rules []string, both bool) *Rule {
	return &Rule{
		rules: rules,
		both:  both,
	}
}

// Check AC
func (r *Rule) Check(rules []string) bool {
	if r.both {
		return r.checkAND(rules)
	}
	return r.checkOR(rules)
}

func (r *Rule) checkAND(rules []string) bool {
MUST:
	for _, must := range r.rules {
		for _, have := range rules {
			if must == have {
				continue MUST
			}
		}
		return false
	}
	return true
}

func (r *Rule) checkOR(rules []string) bool {
	for _, may := range r.rules {
		for _, have := range rules {
			if may == have {
				return true
			}
		}
	}
	return false
}
