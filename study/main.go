package main

import "fmt"

const MAX int = 3

func main() {
	//arrayTest()
	//var numbers = make([]int,3,5)
	//printSlice(numbers)
	//rangeTest()
	mapTest()
}

func arrayTest() {
	a := []int{10,100,200}
	var i int
	var ptr[MAX]* int

	for  i = 0; i < MAX; i++ {
		ptr[i] = &a[i] /* 整数地址赋值给指针数组 */
	}

	for  i = 0; i < MAX; i++ {
		fmt.Printf("a[%d] = %d\n", i,*ptr[i] )
	}
	a = append(a, 1)
	printSlice(a)
	// 声明数组的同时快速初始化数组
	//   balance := [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
}

func printSlice(numbers []int){
	fmt.Printf("len=%d cap=%d slice=%v\n",len(numbers),cap(numbers),numbers)
}

func rangeTest() {
	//这是我们使用range去求一个slice的和。使用数组跟这个很类似
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum:", sum)
	//在数组上使用range将传入index和值两个变量。上面那个例子我们不需要使用该元素的序号，所以我们使用空白符"_"省略了。有时侯我们确实需要知道它的索引。
	for i, num := range nums {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}
	//range也可以用在map的键值对上。
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}
	//range也可以用来枚举Unicode字符串。第一个参数是字符的索引，第二个是字符（Unicode的值）本身。
	for i, c := range "go" {
		fmt.Println(i, c)
	}
}

func mapTest() {
	var countryCapitalMap map[string]string /*创建集合 */
	countryCapitalMap = make(map[string]string)

	/* map插入key - value对,各个国家对应的首都 */
	countryCapitalMap [ "France" ] = "巴黎"
	countryCapitalMap [ "Italy" ] = "罗马"
	countryCapitalMap [ "Japan" ] = "东京"
	countryCapitalMap [ "India " ] = "新德里"

	/*使用键输出地图值 */
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是", countryCapitalMap [country])
	}

	/*查看元素在集合中是否存在 */
	if capital, ok := countryCapitalMap [ "American" ]; ok {
		fmt.Println("American 的首都是", capital)
	} else {
		fmt.Println("American 的首都不存在")
	}
	delete(countryCapitalMap, "France")
	/*打印地图*/
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是", countryCapitalMap [ country ])
	}
}