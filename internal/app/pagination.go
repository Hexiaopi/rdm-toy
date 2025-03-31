package app

type Page struct {
	// 页码
	PageNum int `json:"page_num"`
	// 每页数量
	PageSize int `json:"page_size"`
}

func CorrectPage(size, num int) Page {
	page := Page{
		PageNum:  num,
		PageSize: size,
	}
	if size <= 0 {
		page.PageSize = 10
	}
	if size > 10 {
		page.PageSize = 10
	}
	if num <= 0 {
		page.PageNum = 1
	}
	return page
}

func GetPageOffset(pageNum, pageSize int) int {
	result := 0
	if pageNum > 0 {
		result = (pageNum - 1) * pageSize
	}

	return result
}
