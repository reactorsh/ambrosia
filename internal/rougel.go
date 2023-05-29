package internal

func longestCommonSubsequence(s1, s2 []string) int {
	len1, len2 := len(s1), len(s2)
	dp := make([]int, len2+1)

	for i := 1; i <= len1; i++ {
		prev := dp[0]
		for j := 1; j <= len2; j++ {
			temp := dp[j]
			if s1[i-1] == s2[j-1] {
				dp[j] = prev + 1
			} else if dp[j] < dp[j-1] {
				dp[j] = dp[j-1]
			}
			prev = temp
		}
	}

	return dp[len2]
}

func rougeL(s1, s2 []string) float64 {
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	lcs := longestCommonSubsequence(s1, s2)
	return 2.0 * float64(lcs) / float64(len(s1)+len(s2))
}
