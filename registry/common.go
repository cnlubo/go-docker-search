package registry

import (
	"strings"
)

const (
	HostURL = "https://index.docker.io"
)

// // 通用排序
// // 结构体排序，必须重写数组Len() Swap() Less()函数
// type bodyWrapper struct {
// 	Bodys [] interface{}
// 	by    func(p, q *interface{}) bool // 内部Less()函数会用到
// }
// type SortBodyBy func(p, q *interface{}) bool // 定义一个函数类型
//
// // 数组长度Len()
// func (acw bodyWrapper) Len() int {
// 	return len(acw.Bodys)
// }
//
// // 元素交换
// func (acw bodyWrapper) Swap(i, j int) {
// 	acw.Bodys[i], acw.Bodys[j] = acw.Bodys[j], acw.Bodys[i]
// }
//
// // 比较函数，使用外部传入的by比较函数
// func (acw bodyWrapper) Less(i, j int) bool {
// 	return acw.by(&acw.Bodys[i], &acw.Bodys[j])
// }
//
// // 自定义排序字段，参考SortBodyByCreateTime中的传入函数
// func SortBody(bodys [] interface{}, by SortBodyBy) {
// 	sort.Sort(bodyWrapper{bodys, by})
// }
//
// func SortBodyByStarCount(bodys [] interface{}) {
// 	fmt.Println(bodys)
// 	sort.Sort(bodyWrapper{bodys, func(p, q *interface{}) bool {
// 		v := reflect.ValueOf(*p)
// 		i := v.FieldByName("StarCount")
// 		v = reflect.ValueOf(*q)
// 		j := v.FieldByName("StarCount")
// 		// return  i.String() > j.String()
// 		return i.String() > j.String()
// 	}})
// }

// Removes any "/" prefix & suffix from the given string.
func fixSuffixPrefix(s string) string {

	sep := "/"

	s = strings.TrimPrefix(s, sep)
	if !strings.HasSuffix(s, sep) {
		return s + sep
	}
	return s
}

// Official Docker repos run under the "_" user, so make sure that if a given repo doesn't contain a "/", then "_/" needs to be added to the repo name.
// Example:
// input: mongo => _/mongo
func fixOfficialRepos(s string) string {
	if !strings.Contains(s, "/") {
		return "library/" + s
	}
	return s
}