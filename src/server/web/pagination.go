package web

import "strconv"

func parsePagination(page string, size string, defaultPage int, defaultSize int) (pageInt int, sizeInt int) {
	if len(page) < 1 {
		page = strconv.Itoa(defaultPage)
	}

	if len(size) < 1 {
		size = strconv.Itoa(defaultSize)
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = defaultPage
	}

	sizeInt, err = strconv.Atoi(size)
	if err != nil {
		sizeInt = defaultSize
	}

	return pageInt, sizeInt
}
