package githubtrending


import (
"testing"
)

func TestGithubTrendingContent(t *testing.T) {
	for _, one := range []string{""} {
		for _, since := range []string{""} {
			TrendingContent(one, since)
		}
	}

}