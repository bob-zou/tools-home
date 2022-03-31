package pmath

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"tools-home/internal/service/utils"
)

type QType int

const (
	QTypeNormal QType = iota + 1
	QTypeAdvanced
)

var (
	_bit2Questions       = make(map[int][]*Bit2Question) // map[max]...
	_bit2StringQuestions = make(map[int]map[QType][]string)
	_bit3Questions       = make(map[int][]*Bit3Question)
	_bit3StringQuestions = make(map[int]map[QType][]string)
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}

func questionString(args ...interface{}) string {
	var buffer bytes.Buffer
	for _, a := range args {
		switch v := a.(type) {
		case int:
			buffer.WriteString(strconv.Itoa(v))
		case string:
			buffer.WriteString(v)
		}
	}
	return buffer.String()
}

type QBank interface {
	RandomQuestions(int, ...QType) [][]string
}

func NewQBank(bit int) QBank {
	var qBank QBank
	qBank = &Q2Bank{}

	if bit == 3 {
		qBank = &Q3Bank{}
		return qBank
	}

	return qBank
}

type Q2Bank struct{}

func (q *Q2Bank) RandomQuestions(max int, qts ...QType) [][]string {
	defer func(start time.Time) {
		fmt.Printf("Q2Bank cost: %dms\n", time.Since(start).Milliseconds())
	}(time.Now())
	var (
		column, row = 4, 25
		page        [][]string
		b2qs        = generateB2Questions(max)
		questions   []string
		count       = 0
		used        = make(map[int]struct{})
	)

	if len(qts) == 0 || utils.SliceContains(qts, QTypeNormal) {
		questions = append(questions, b2qs[QTypeNormal]...)
	}

	if len(qts) == 0 || utils.SliceContains(qts, QTypeAdvanced) {
		questions = append(questions, b2qs[QTypeAdvanced]...)
	}

	count = len(questions)
	for i := 0; i < row; i++ {
		page = append(page, []string{})
		for j := 0; j < column; j++ {
			for {
				index := rand.Intn(count)
				if _, ok := used[index]; ok {
					continue
				}
				used[index] = struct{}{}
				page[i] = append(page[i], questions[index])
				break
			}

		}
	}

	return page
}

// Bit2Question A <= B <= Z
// A + B = Z
type Bit2Question struct {
	A int
	B int
	Z int
}

func (b2 *Bit2Question) Questions(qts ...QType) map[QType][]string {
	ret := make(map[QType][]string)

	if len(qts) == 0 || utils.SliceContains(qts, QTypeNormal) {
		ret[QTypeNormal] = []string{
			fmt.Sprintf("%2d + %2d = %6s", b2.A, b2.B, "______"),
			fmt.Sprintf("%2d - %2d = %6s", b2.Z, b2.A, "______"),
		}
	}

	if len(qts) == 0 || utils.SliceContains(qts, QTypeAdvanced) {
		ret[QTypeAdvanced] = []string{
			fmt.Sprintf("%2d + %6s = %2d", b2.A, "______", b2.Z),
			fmt.Sprintf("%6s + %2d = %2d", "______", b2.A, b2.Z),
			fmt.Sprintf("%2d - %6s = %2d", b2.Z, "______", b2.A),
			fmt.Sprintf("%6s - %2d = %2d", "______", b2.A, b2.B),
		}
	}

	if b2.A != b2.B {
		if len(qts) == 0 || utils.SliceContains(qts, QTypeNormal) {
			ret[QTypeNormal] = append(ret[QTypeNormal], []string{
				fmt.Sprintf("%2d + %2d = %6s", b2.B, b2.A, "______"),
				fmt.Sprintf("%2d - %2d = %6s", b2.Z, b2.B, "______"),
			}...)
		}

		if len(qts) == 0 || utils.SliceContains(qts, QTypeAdvanced) {
			ret[QTypeAdvanced] = append(ret[QTypeAdvanced], []string{
				fmt.Sprintf("%2d + %6s = %2d", b2.B, "______", b2.Z),
				fmt.Sprintf("%6s + %2d = %2d", "______", b2.B, b2.Z),
				fmt.Sprintf("%2d - %6s = %2d", b2.Z, "______", b2.B),
				fmt.Sprintf("%6s - %2d = %2d", "______", b2.B, b2.A),
			}...)
		}
	}
	return ret
}

func generateB2Qs(max int) []*Bit2Question {
	if b2qs, ok := _bit2Questions[max]; ok {
		return b2qs
	}
	var b2qs []*Bit2Question
	for i := 1; i < max; i++ {
		for j := i; j < max; j++ {
			if i+j <= max {
				b2qs = append(b2qs, &Bit2Question{A: i, B: j, Z: i + j})
			}
		}
	}

	_bit2Questions[max] = b2qs

	return b2qs
}

func generateB2Questions(max int) map[QType][]string {
	if b2Questions, ok := _bit2StringQuestions[max]; ok {
		return b2Questions
	}

	b2qs := generateB2Qs(max)
	b2Questions := make(map[QType][]string)
	for _, b2q := range b2qs {
		tmpQs := b2q.Questions()
		b2Questions[QTypeNormal] = append(b2Questions[QTypeNormal], tmpQs[QTypeNormal]...)
		b2Questions[QTypeAdvanced] = append(b2Questions[QTypeAdvanced], tmpQs[QTypeAdvanced]...)
	}

	return b2Questions
}

// Bit3Question A <= B <= C <= Z
// A + B + C = Z
type Bit3Question struct {
	A int
	B int
	C int
	Z int
}

func (b3 *Bit3Question) Questions(qts ...QType) map[QType][]string {
	ret := make(map[QType][]string)

	if len(qts) == 0 || utils.SliceContains(qts, QTypeNormal) {
		// A == B == C
		ret[QTypeNormal] = []string{
			questionString(b3.A, " + ", b3.B, " + ", b3.C, " = ", "______"),
			questionString(b3.Z, " - ", b3.A, " - ", b3.B, " = ", "______"),
		}

		if b3.A != b3.B && b3.B != b3.C && b3.A != b3.C {
			ret[QTypeNormal] = []string{
				questionString(b3.A, " + ", b3.B, " + ", b3.C, " = ", "______"),
				questionString(b3.A, " + ", b3.C, " + ", b3.B, " = ", "______"),
				questionString(b3.B, " + ", b3.A, " + ", b3.C, " = ", "______"),
				questionString(b3.B, " + ", b3.C, " + ", b3.A, " = ", "______"),
				questionString(b3.C, " + ", b3.A, " + ", b3.B, " = ", "______"),
				questionString(b3.C, " + ", b3.B, " + ", b3.A, " = ", "______"),
				questionString(b3.Z, " - ", b3.A, " - ", b3.B, " = ", "______"),
				questionString(b3.Z, " - ", b3.A, " - ", b3.C, " = ", "______"),
				questionString(b3.Z, " - ", b3.B, " - ", b3.A, " = ", "______"),
				questionString(b3.Z, " - ", b3.B, " - ", b3.C, " = ", "______"),
				questionString(b3.Z, " - ", b3.C, " - ", b3.A, " = ", "______"),
				questionString(b3.Z, " - ", b3.C, " - ", b3.B, " = ", "______"),
			}
		}
		if b3.A != b3.B && b3.B == b3.C {
			ret[QTypeNormal] = []string{
				questionString(b3.A, " + ", b3.B, " + ", b3.C, " = ", "______"),
				questionString(b3.B, " + ", b3.A, " + ", b3.C, " = ", "______"),
				questionString(b3.B, " + ", b3.C, " + ", b3.A, " = ", "______"),
				questionString(b3.Z, " - ", b3.A, " - ", b3.B, " = ", "______"),
				questionString(b3.Z, " - ", b3.B, " - ", b3.A, " = ", "______"),
				questionString(b3.Z, " - ", b3.B, " - ", b3.C, " = ", "______"),
			}
		}

		if b3.A != b3.B && b3.A == b3.C {
			ret[QTypeNormal] = []string{
				questionString(b3.A, " + ", b3.B, " + ", b3.C, " = ", "______"),
				questionString(b3.A, " + ", b3.C, " + ", b3.B, " = ", "______"),
				questionString(b3.B, " + ", b3.A, " + ", b3.C, " = ", "______"),
				questionString(b3.Z, " - ", b3.A, " - ", b3.B, " = ", "______"),
				questionString(b3.Z, " - ", b3.A, " - ", b3.C, " = ", "______"),
				questionString(b3.Z, " - ", b3.B, " - ", b3.A, " = ", "______"),
			}
		}

		if b3.A == b3.B && b3.A != b3.C {
			ret[QTypeNormal] = []string{
				questionString(b3.A, " + ", b3.B, " + ", b3.C, " = ", "______"),
				questionString(b3.A, " + ", b3.C, " + ", b3.B, " = ", "______"),
				questionString(b3.C, " + ", b3.A, " + ", b3.B, " = ", "______"),
				questionString(b3.Z, " - ", b3.A, " - ", b3.B, " = ", "______"),
				questionString(b3.Z, " - ", b3.A, " - ", b3.C, " = ", "______"),
				questionString(b3.Z, " - ", b3.C, " - ", b3.A, " = ", "______"),
			}
		}
	}

	if len(qts) == 0 || utils.SliceContains(qts, QTypeAdvanced) {
		ret[QTypeAdvanced] = []string{
			fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.A, b3.B, b3.Z),
			fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.B, b3.Z),
			fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.B, "______", b3.Z),
			fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.C),
			fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.A, b3.C),
			fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.B, b3.C),
		}

		if b3.A != b3.B && b3.B != b3.C && b3.A != b3.C {
			ret[QTypeAdvanced] = []string{
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.B, b3.C, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.C, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.B, "______", b3.Z),
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.C, b3.B, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.B, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.C, "______", b3.Z),
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.A, b3.C, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.B, "______", b3.C, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.B, b3.A, "______", b3.Z),
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.C, b3.A, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.B, "______", b3.A, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.B, b3.C, "______", b3.Z),
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.A, b3.B, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.C, "______", b3.B, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.C, b3.A, "______", b3.Z),
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.B, b3.A, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.C, "______", b3.A, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.C, b3.B, "______", b3.Z),
				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.B, b3.C),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.B, b3.C),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.C),
				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.C, b3.B),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.C, b3.B),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.B),
				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.B, b3.A, b3.C),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.A, b3.C),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.B, "______", b3.C),
				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.B, b3.C, b3.A),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.C, b3.A),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.B, "______", b3.A),
				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.C, b3.A, b3.B),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.A, b3.B),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.C, "______", b3.B),
				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.C, b3.B, b3.A),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.B, b3.A),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.C, "______", b3.A),
			}
		}
		if b3.A != b3.B && b3.B == b3.C {
			ret[QTypeAdvanced] = []string{
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.B, b3.C, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.C, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.B, "______", b3.Z),

				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.A, b3.C, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.B, "______", b3.C, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.B, b3.A, "______", b3.Z),

				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.C, b3.A, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.B, "______", b3.A, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.B, b3.C, "______", b3.Z),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.B, b3.C),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.B, b3.C),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.C),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.B, b3.A, b3.C),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.A, b3.C),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.B, "______", b3.C),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.B, b3.C, b3.A),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.C, b3.A),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.B, "______", b3.A),
			}
		}

		if b3.A != b3.B && b3.A == b3.C {
			ret[QTypeAdvanced] = []string{
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.B, b3.C, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.C, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.B, "______", b3.Z),

				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.C, b3.B, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.B, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.C, "______", b3.Z),

				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.A, b3.C, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.B, "______", b3.C, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.B, b3.A, "______", b3.Z),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.B, b3.C),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.B, b3.C),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.C),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.C, b3.B),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.C, b3.B),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.B),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.B, b3.A, b3.C),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.A, b3.C),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.B, "______", b3.C),
			}
		}

		if b3.A == b3.B && b3.A != b3.C {
			ret[QTypeAdvanced] = []string{
				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.B, b3.C, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.C, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.B, "______", b3.Z),

				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.C, b3.B, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.A, "______", b3.B, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.A, b3.C, "______", b3.Z),

				fmt.Sprintf("%6s + %2d + %2d = %2d", "______", b3.A, b3.B, b3.Z),
				fmt.Sprintf("%2d + %6s + %2d = %2d", b3.C, "______", b3.B, b3.Z),
				fmt.Sprintf("%2d + %2d + %6s = %2d", b3.C, b3.A, "______", b3.Z),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.B, b3.C),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.B, b3.C),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.C),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.A, b3.C, b3.B),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.C, b3.B),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.A, "______", b3.B),

				fmt.Sprintf("%6s - %2d - %2d = %2d", "______", b3.C, b3.A, b3.B),
				fmt.Sprintf("%2d - %6s - %2d = %2d", b3.Z, "______", b3.A, b3.B),
				fmt.Sprintf("%2d - %2d - %6s = %2d", b3.Z, b3.C, "______", b3.B),
			}
		}
	}

	return ret
}

func generateB3Qs(max int) []*Bit3Question {
	if b3qs, ok := _bit3Questions[max]; ok {
		return b3qs
	}
	var b3qs []*Bit3Question
	for i := 1; i < max; i++ {
		for j := i; j < max; j++ {
			if i+j > max {
				break
			}
			for k := j; k < max; k++ {
				if i+j+k > max {
					break
				}
				b3qs = append(b3qs, &Bit3Question{A: i, B: j, C: k, Z: i + j + k})
			}
		}
	}
	_bit3Questions[max] = b3qs

	return b3qs
}

func generateB3Questions(max int) map[QType][]string {
	if b3Questions, ok := _bit3StringQuestions[max]; ok {
		return b3Questions
	}

	b3qs := generateB3Qs(max)
	b3Questions := make(map[QType][]string)
	for _, b3q := range b3qs {
		tmpQs := b3q.Questions()
		b3Questions[QTypeNormal] = append(b3Questions[QTypeNormal], tmpQs[QTypeNormal]...)
		b3Questions[QTypeAdvanced] = append(b3Questions[QTypeAdvanced], tmpQs[QTypeAdvanced]...)
	}

	_bit3StringQuestions[max] = b3Questions
	return b3Questions
}

type Q3Bank struct{}

func (q *Q3Bank) RandomQuestions(max int, qts ...QType) [][]string {
	var (
		column, row = 3, 25
		page        [][]string
		b3qs        = generateB3Questions(max)
		questions   []string
		count       = 0
		used        = make(map[int]struct{})
	)

	if len(qts) == 0 || utils.SliceContains(qts, QTypeNormal) {
		questions = append(questions, b3qs[QTypeNormal]...)
	}

	if len(qts) == 0 || utils.SliceContains(qts, QTypeAdvanced) {
		questions = append(questions, b3qs[QTypeAdvanced]...)
	}

	count = len(questions)
	for i := 0; i < row; i++ {
		page = append(page, []string{})
		for j := 0; j < column; j++ {
			for {
				index := rand.Intn(count)
				if _, ok := used[index]; ok {
					continue
				}
				used[index] = struct{}{}
				page[i] = append(page[i], questions[index])
				break
			}

		}
	}

	return page
}
