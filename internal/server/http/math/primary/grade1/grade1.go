package grade1

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"tools-home/internal/service/qbank/pmath"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func queryBoolDefault(c *gin.Context, key string, def bool) bool {
	v, ok := c.GetQuery(key)
	if !ok {
		return def
	}
	if v == "1" || strings.ToLower(v) == "true" {
		return true
	}

	return false
}

func queryIntDefault(c *gin.Context, key string, def int) int {
	v, ok := c.GetQuery(key)
	if !ok {
		return def
	}
	vi, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return vi
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

func generatePdf(pages [][][]string) (filename string) {
	file, err := ioutil.TempFile("", "qs-")
	if err != nil {
		return
	}
	filename = file.Name()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.SetFont("Arial", "B", 14)

	for _, page := range pages {
		pdf.AddPage()
		for i, r := range page {
			var w float64
			switch len(r) {
			case 3:
				w = 60
			default:
				w = 47.5
			}
			for _, c := range r {
				pdf.Cell(w, 12, c)
			}
			pdf.LineTo(10, float64(15+10*(i+1)))
		}
	}
	_ = pdf.OutputFileAndClose(filename)
	return
}

func Bit2Questions(c *gin.Context) {
	var (
		page  = queryIntDefaultMax(c, "page", 10, 100)
		max   = queryIntDefaultMax(c, "max", 20, 100)
		qType = queryIntDefault(c, "type", 0)
		bit   = queryIntDefault(c, "bit", 2)
		name  = fmt.Sprintf("%d以内二元加减法-%d套-%d.pdf", max, page, time.Now().Unix())
	)
	if bit != 2 && bit != 3 {
		bit = 2
	}

	if bit == 3 {
		name = fmt.Sprintf("%d以内三元加减法-%d套-%d.pdf", max, page, time.Now().Unix())
	}

	var (
		pages [][][]string
		qBank = pmath.NewQBank(bit)
	)
	for i := 0; i < page; i++ {
		var questions [][]string
		if qType == 0 || qType > 2 {
			questions = qBank.RandomQuestions(max)
		} else {
			questions = qBank.RandomQuestions(max, pmath.QType(qType))
			if pmath.QType(qType) == pmath.QTypeNormal {
				name = fmt.Sprintf("%d以内二元基础加减法-%d套-%d.pdf", max, page, time.Now().Unix())
				if bit == 3 {
					name = fmt.Sprintf("%d以内三元基础加减法-%d套-%d.pdf", max, page, time.Now().Unix())
				}
			} else if pmath.QType(qType) == pmath.QTypeAdvanced {
				name = fmt.Sprintf("%d以内二元进阶加减法-%d套-%d.pdf", max, page, time.Now().Unix())
				if bit == 3 {
					name = fmt.Sprintf("%d以内三元进阶加减法-%d套-%d.pdf", max, page, time.Now().Unix())
				}
			}
		}
		pages = append(pages, questions)
	}

	filename := generatePdf(pages)

	c.FileAttachment(filename, url.QueryEscape(name))
	defer func() { _ = os.RemoveAll(filename) }()
}
