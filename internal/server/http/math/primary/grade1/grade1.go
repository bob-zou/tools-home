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

func GetTuples(max int, withZero bool) (tuples []*Tuple) {
	start := 1
	if withZero {
		start = 0
	}
	for i := start; i < max; i++ {
		for j := i; j < max; j++ {
			if i+j < max {
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

func randomQuestion(tuples []*Tuple, used map[string]struct{}) string {
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

func generatePdf(pageSize, max int, name string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.SetFont("Arial", "B", 16)

	rand.Seed(time.Now().UnixMilli())

	tuples := GetTuples(max, false)

	page := 1
	for {
		pdf.AddPage()
		count := 0
		used := map[string]struct{}{}
		for {
			pdf.Cell(47.5, 12, randomQuestion(tuples, used))
			if count%4 == 3 {
				pdf.LineTo(10, 15+10*(math.Ceil(float64(count/4))+1))
			}

			count++
			if count >= 100 {
				break
			}
		}
		page++
		if page > pageSize {
			break
		}
	}
	_ = pdf.OutputFileAndClose(name)
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

func RandomQuestions(c *gin.Context) {
	var (
		page = queryIntDefaultMax(c, "page", 10, 100)
		max  = queryIntDefaultMax(c, "max", 20, 100)
	)

	file, err := ioutil.TempFile("", "question-")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	generatePdf(page, max, file.Name())
	c.FileAttachment(file.Name(), url.QueryEscape(fmt.Sprintf("%d以内加减法-%d套-%d.pdf", max, page, time.Now().Unix())))
	defer func() { _ = os.RemoveAll(file.Name()) }()
	//c.JSON(http.StatusOK, common.Reply{Data: "ok"})
}
