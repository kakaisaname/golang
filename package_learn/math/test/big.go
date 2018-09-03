package main

import (
	"math/big"
	"fmt"
	"math"
)
//Big包实现了任意精度算术（大数）。支持以下数字类型：
//Int    signed integers
//Rat    rational numbers
//Float  floating-point numbers


//使用e的经典连续分数
// e = [1; 0,1,1,2,1,1，... 2n，1,1，......]
//即，对于第n个术语，请使用
// 1如果n mod 3！= 1
//（n-1）/ 3 * 2如果n mod 3 == 1
func recur(n, lim int64) *big.Rat {
	term := new(big.Rat)
	if n%3 != 1 {
		term.SetInt64(1)
	} else {
		term.SetInt64((n - 1) / 3 * 2)
	}

	if n > lim {
		return term
	}

	//直接将frac初始化为小数
	//recur结果的倒数
	frac := new(big.Rat).Inv(recur(n+1, lim))

	return term.Add(term, frac)
}

//这个实例演示了如何使用big.Rat来计算
//在理性收敛序列中的前15个项
//常数e（自然对数的基数
func main() {
	for i := 1; i <= 15; i++ {
		r := recur(0, int64(i))


		//将r打印为分数和浮点数。
		//由于big.Rat实现了fmt.Formatter，我们可以使用％-13s
		//得到分数的左对齐字符串表示。
		fmt.Printf("%-13s = %s\n", r, r.FloatString(8))
	}

}

//使用 big.Float 以200位的精度计算2的平方根，以及如何将结果打印为十进制数。

func a() {
	// 我们将做计算与200位精度在尾数。
	const prec = 200

	// 用牛顿法计算2的平方根。我们从
	// sqrt (2) 的初始估计值, 然后循环访问:
	//     x_{n+1} = 1/2 * ( x_n + (2.0 / x_n) )

	// 因为牛顿的方法加倍的正确数字在每个
	// 迭代, 我们至少需要 log_2 (prec) 步骤。
	steps := int(math.Log2(prec))

	// 初始化计算所需的值。
	two := new(big.Float).SetPrec(prec).SetInt64(2)
	half := new(big.Float).SetPrec(prec).SetFloat64(0.5)

	// 使用1作为初始估计值。
	x := new(big.Float).SetPrec(prec).SetInt64(1)

	// 我们使用 t 作为一个临时变量。没有必要设定它的精确度
	// 从大。浮动值 (== 0) 精度自动承担
	// 作为结果使用时参数的最大精度 (接收者)
	// 大。浮动操作。
	t := new(big.Float)

	// 迭代。
	for i := 0; i <= steps; i++ {
		t.Quo(two, x)  // t = 2.0 / x_n
		t.Add(x, t)    // t = x_n + (2.0 / x_n)
		x.Mul(half, t) // x_{n+1} = 0.5 * t
	}

	// 我们可以使用常规的裂变材料。自大以来的 Printf 动词。浮动实现了裂变材料。格式
	fmt.Printf("sqrt(2) = %.50f\n", x)

	// 打印2和之间的x*x错误。
	t.Mul(x, x) // t = x*x
	fmt.Printf("error = %e\n", t.Sub(two, t))

}