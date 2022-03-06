package grade1

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type Tuple struct {
	A int
	B int
	C int
}

func (t *Tuple) String() string {
	return fmt.Sprintf("%2d + %2d = %2d", t.A, t.B, t.C)
}

func (t *Tuple) Question(op string) string {
	if op == "+" {
		return fmt.Sprintf("%2d + %2d = %6s", t.A, t.B, "______")
	}

	return fmt.Sprintf("%2d - %2d = %6s", t.C, t.A, "______")
}

func (t *Tuple) Question2(op string) string {
	sel := rand.Intn(3)
	switch sel {
	case 1:
		if op == "+" {
			return fmt.Sprintf("%2d + %6s = %2d", t.A, "______", t.C)
		}
		return fmt.Sprintf("%2d - %6s = %2d", t.C, "______", t.A)
	case 2:
		if op == "+" {
			return fmt.Sprintf("%6s + %2d = %2d", "______", t.B, t.C)
		}
		return fmt.Sprintf("%6s - %2d = %2d", "______", t.B, t.A)
	default:
		if op == "+" {
			return fmt.Sprintf("%2d + %2d = %6s", t.A, t.B, "______")
		}
		return fmt.Sprintf("%2d - %2d = %6s", t.C, t.A, "______")
	}
}

func queryIntDefaultMax(c *gin.Context, key string, def, max int) int {
	v, ok := c.GetQuery(key)
	if !ok {
		return def
	}
	vi, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	if vi > max {
		return max
	}

	return vi
}

func getTuples(max int, withZero bool) (tuples []*Tuple) {
	start := 1
	if withZero {
		start = 0
	}
	for i := start; i < max; i++ {
		for j := i; j < max; j++ {
			if i+j <= max {
				tuples = append(tuples, &Tuple{
					A: i,
					B: j,
					C: i + j,
				})
			}
		}
	}

	return
}

func randomOneBaseQuestion(tuples []*Tuple, used map[string]struct{}) string {
	var (
		index int
		op    string
	)

	for {
		i := rand.Intn(len(tuples))
		o := "+"
		if rand.Intn(2) == 0 {
			o = "-"
		}
		k := fmt.Sprintf("%d-%s", i, o)
		if _, ok := used[k]; ok {
			continue
		}
		index, op = i, o
		used[k] = struct{}{}
		break
	}

	return tuples[index].Question(op)
}

func randomOneBase2Question(tuples []*Tuple, used map[string]struct{}) string {
	var (
		index int
		op    string
	)

	for {
		i := rand.Intn(len(tuples))
		o := "+"
		if rand.Intn(2) == 0 {
			o = "-"
		}
		k := fmt.Sprintf("%d-%s", i, o)
		if _, ok := used[k]; ok {
			continue
		}
		index, op = i, o
		used[k] = struct{}{}
		break
	}

	return tuples[index].Question2(op)
}

func generatePdf(questions map[int][]string, name string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.SetFont("Arial", "B", 14)

	for _, qs := range questions {
		count := 0
		pdf.AddPage()
		for _, q := range qs {
			pdf.Cell(47.5, 12, q)
			if count%4 == 3 {
				pdf.LineTo(10, 15+10*(math.Ceil(float64(count/4))+1))
			}
			count++
		}
	}
	_ = pdf.OutputFileAndClose(name)
}

func randomBaseQuestions(page, max int) (questions map[int][]string) {
	var (
		cur    = 1
		tuples = getTuples(max, false)
	)

	rand.Seed(time.Now().UnixMilli())
	questions = make(map[int][]string, page)
	for {
		count := 0
		used := map[string]struct{}{}
		questions[cur] = make([]string, 0, 100)
		for {
			questions[cur] = append(questions[cur], randomOneBaseQuestion(tuples, used))
			count++
			if count >= 100 {
				break
			}
		}
		cur++
		if cur > page {
			break
		}
	}
	return
}

func randomBase2Questions(page, max int) (questions map[int][]string) {
	var (
		cur    = 1
		tuples = getTuples(max, false)
	)

	rand.Seed(time.Now().UnixMilli())
	questions = make(map[int][]string, page)
	for {
		count := 0
		used := map[string]struct{}{}
		questions[cur] = make([]string, 0, 100)
		for {
			questions[cur] = append(questions[cur], randomOneBase2Question(tuples, used))
			count++
			if count >= 100 {
				break
			}
		}
		cur++
		if cur > page {
			break
		}
	}
	return
}

func RandomBaseQuestions(c *gin.Context) {
	var (
		page      = queryIntDefaultMax(c, "page", 10, 100)
		max       = queryIntDefaultMax(c, "max", 20, 100)
		questions = randomBaseQuestions(page, max)
	)

	file, err := ioutil.TempFile("", "question-")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	generatePdf(questions, file.Name())
	c.FileAttachment(file.Name(), url.QueryEscape(fmt.Sprintf("%d以内加减法-%d套-%d.pdf", max, page, time.Now().Unix())))
	defer func() { _ = os.RemoveAll(file.Name()) }()
}

func RandomBase2Questions(c *gin.Context) {
	var (
		page      = queryIntDefaultMax(c, "page", 10, 100)
		max       = queryIntDefaultMax(c, "max", 20, 100)
		questions = randomBase2Questions(page, max)
	)

	file, err := ioutil.TempFile("", "question-")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	generatePdf(questions, file.Name())
	c.FileAttachment(file.Name(), url.QueryEscape(fmt.Sprintf("%d以内加减法-%d套-%d.pdf", max, page, time.Now().Unix())))
	defer func() { _ = os.RemoveAll(file.Name()) }()
}
