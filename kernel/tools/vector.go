package tools

import (
	"fmt"
	"math"
	"math/rand"
)

type (
	Vector2 struct {
		X float64
		Y float64
	}
	Vector3 struct {
		X float64
		//Y float64	// 对于服务器来说Y没用
		Z float64
	}
)

func NewVector3(x, z float64) *Vector3 {
	return &Vector3{X: x, Z: z}
}

func CopyVector3(v *Vector3) *Vector3 {
	return NewVector3(v.X, v.Z)
}

func NewVector3Zero() *Vector3 {
	return NewVector3(0, 0)
}

func ToFloatVector(v []int, minification float64) *Vector3 {
	if minification == 0 {
		minification = 1000
	}
	var divideOne = 1.0 / minification
	return &Vector3{X: float64(v[0]) * divideOne, Z: float64(v[1]) * divideOne}
}

func (v *Vector3) Pos() (float64, float64) {
	return v.X, v.Z
}

func (v *Vector3) SetPos(x, z float64) {
	v.X, v.Z = x, z
}

// Magnitude 返回向量的模长（长度）
func (v *Vector3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Z*v.Z)
}

/*
 * Neg
 *  @Description: 向量取负
 *  @return *Vector3
 */
func (v *Vector3) Neg() *Vector3 {
	return NewVector3(-v.X, -v.Z)
}

/*
 * MulNum
 *  @Description: 向量乘以一个数
 *  @param f 乘数
 *  @return *Vector3 结果向量
 */
func (v *Vector3) MulNum(f float64) *Vector3 {
	return NewVector3(v.X*f, v.Z*f)
}

// Add 计算两个向量的和
func (v *Vector3) Add(other *Vector3) *Vector3 {
	return NewVector3(v.X+other.X, v.Z+other.Z)
}

// Subtract 计算两个向量的差
func (v *Vector3) Subtract(other *Vector3) *Vector3 {
	return NewVector3(v.X-other.X, v.Z-other.Z)
}

// Dot 计算两个向量的点乘
func (v *Vector3) Dot(other *Vector3) float64 {
	return v.X*other.X + v.Z*other.Z
}

// Normalize 计算向量的单位向量
func (v *Vector3) Normalize() *Vector3 {
	magnitude := v.Magnitude()
	if magnitude < 1e-6 {
		return NewVector3(0, 0)
	}
	return NewVector3(v.X/magnitude, v.Z/magnitude)
}

func RangeRand(min, max int) int {
	rangeData := (max - min) + 1
	return rand.Intn(rangeData) + min
}

// Cross 计算两个向量的叉积（在XZ平面中返回标量）
func (v *Vector3) Cross(other *Vector3) float64 {
	// 在2D平面(XZ)中，叉积结果是一个标量: ax*bz - az*bx
	return v.X*other.Z - v.Z*other.X
}

// SqrMagnitude 返回向量模长的平方
func (v *Vector3) SqrMagnitude() float64 {
	return v.X*v.X + v.Z*v.Z
}

// 如果需要保持原有的垂直向量功能，可以添加一个新的方法
// Perpendicular 返回在XZ平面中垂直于原向量的向量
func (v *Vector3) Perpendicular() *Vector3 {
	// 返回一个在XZ平面中垂直于原向量的向量
	return NewVector3(-v.Z, v.X)
}

func GetIntersectPoint(a, b, c, d *Vector3) (bool, *Vector3) {
	ab := b.Subtract(a)
	ca := a.Subtract(c)
	cd := d.Subtract(c)

	// 使用修正后的Cross方法计算叉积标量
	v1 := ca.Cross(cd)
	if math.Abs(v1*ab.Dot(ab)) > 1e-6 { // 注意这里改为点乘
		return false, nil
	}

	abCrossCd := ab.Cross(cd) // 叉积标量
	if math.Abs(abCrossCd) <= 1e-6 {
		return false, nil
	}

	ad := d.Subtract(a)
	cb := b.Subtract(c)

	if math.Min(a.X, b.X) > math.Max(c.X, d.X) ||
		math.Max(a.X, b.X) < math.Min(c.X, d.X) ||
		math.Min(a.Z, b.Z) > math.Max(c.Z, d.Z) ||
		math.Max(a.Z, b.Z) < math.Min(c.Z, d.Z) {
		return false, nil
	}

	nca := ca.Neg()
	// 修改叉积计算方式
	if !(nca.Cross(ab)*ab.Cross(ad) > 0) ||
		!(ca.Cross(cd)*cd.Cross(cb) > 0) {
		return false, nil
	}

	var v2 = abCrossCd // 直接使用标量值
	var ratio = v1 / v2
	intersectPos := a.Add(ab.MulNum(ratio))
	return true, intersectPos
}

func (v *Vector3) String() string {
	return fmt.Sprintf("x:%v z:%v", v.X, v.Z)
}
