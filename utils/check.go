package utils

type Check struct {
	Name   string
	Ok     bool
	Weight float64
}

type Result struct {
	Score        float64 // 0 → 1
	Passed       bool
	TotalWeight  float64
	FailedWeight float64
	Issues       []string
}

// Options allow future extension without breaking API
type Options struct {
	MinScore float64 // default: 1 (all must pass)
}

func EvaluateChecks(checks []Check, opts *Options) Result {
	if opts == nil {
		opts = &Options{MinScore: 1}
	}

	var total, failed float64
	issues := make([]string, 0, len(checks))

	for _, c := range checks {
		// Ignore invalid weights
		if c.Weight <= 0 {
			continue
		}

		total += c.Weight

		if !c.Ok {
			failed += c.Weight
			issues = append(issues, c.Name)
		}
	}

	var score float64
	if total == 0 {
		score = 1 // define empty as perfect (or change to 0 if stricter)
	} else {
		score = 1 - (failed / total)
	}

	return Result{
		Score:        score,
		Passed:       score >= opts.MinScore,
		TotalWeight:  total,
		FailedWeight: failed,
		Issues:       issues,
	}
}
