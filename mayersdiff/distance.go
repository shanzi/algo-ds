package mayersdiff

func upzero(a int) int {
	if a < 0 {
		return 0
	}
	return a
}

func getV(V []int, idx int) int {
	if idx >= 0 {
		return V[idx]
	} else {
		return V[len(V)+idx]
	}
}

func setV(V []int, idx int, val int) {
	if idx >= 0 {
		V[idx] = val
	} else {
		V[len(V)+idx] = val
	}
}

func EditingDistanceStrings(a, b string) int {
	return EditingDistance(([]rune)(a), ([]rune)(b))
}

func EditingDistance(a, b []rune) int {
	M := len(a)
	N := len(b)
	MAX := M + N

	V := make([]int, MAX+2)
	for i := 0; i < len(V); i++ {
		V[i] = MAX
	}

	setV(V, 1, 0)
	for D := 0; D <= MAX; D++ {
		var x int
		for k := -(D - 2*upzero(D-M)); k <= D-2*upzero(D-N); k += 2 {
			if k == -D || (k != D && getV(V, k-1) < getV(V, k+1)) {
				x = getV(V, k+1)
			} else {
				x = getV(V, k-1) + 1
			}
			y := x - k
			for x < M && y < N && a[x] == b[y] {
				x += 1
				y += 1
			}

			setV(V, k, x)

			if x == M && y == N {
				return D
			}
		}
	}
	return MAX
}
