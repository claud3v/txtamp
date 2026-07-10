package mainview

func visibleRange(selected, total, visible int) (int, int) {
	if total <= 0 || visible <= 0 {
		return 0, 0
	}
	if visible >= total {
		return 0, total
	}

	start := selected - visible + 1
	if start < 0 {
		start = 0
	}

	end := start + visible
	if end > total {
		end = total
		start = end - visible
	}

	return start, end
}
