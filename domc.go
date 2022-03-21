package domc

import (
	"fmt"
	"hash/fnv"
	"math"
	"strings"

	"golang.org/x/net/html"
)

const (
	defaultDimension     uint32 = 5000
	defaultAttenuation          = 0.6
	defaultInitialWeight        = 1.0
)

var (
	Dimension     = defaultDimension
	Attenuation   = defaultAttenuation
	InitialWeight = defaultInitialWeight
)

type Vector map[uint32]float64

func (v Vector) IsSimilar(vector Vector) bool {
	var a, b float64

	for d := range v {
		a += math.Abs(v[d] - vector[d])
		b += v[d] + vector[d]
	}

	if math.Abs(a)/b < 0.2 {
		return true
	}

	return false
}

func (v Vector) Compress(dimension uint32) Vector {
	newVector := make(Vector)

	if dimension <= 0 {
		dimension = Dimension
	}
	for d, w := range v {
		nd := d % dimension
		newVector[nd] += w
	}

	return newVector
}

func NewVector(ns []*Node) Vector {
	v := make(Vector)

	for _, n := range ns {
		v[n.Dimension()] += n.Weight()
	}

	return v
}

type Node struct {
	Name               string
	Depth              int
	SiblingRepeatTimes int
	AttrStr            string
}

func (n *Node) Dimension() uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(n.String()))
	return h.Sum32()
}

func (n *Node) String() string {
	return fmt.Sprintf("%s:%s", n.Name, n.AttrStr)
}

func (n *Node) Weight() float64 {
	var weight float64

	if n.Depth > 0 {
		weight = InitialWeight * math.Pow(Attenuation, float64(n.Depth))
	}
	if n.SiblingRepeatTimes > 0 {
		weight *= math.Pow(Attenuation, float64(n.SiblingRepeatTimes))
	}

	return weight
}

func ParseHTML(h string) ([]*Node, error) {
	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		return nil, err
	}

	var ns []*Node
	var f func(*html.Node, int, *int)
	f = func(n *html.Node, d int, r *int) {
		if n.Type == html.ElementNode && n.Data != "script" && n.Data != "style" {
			var kvs []string
			for _, a := range n.Attr {
				kvs = append(kvs, fmt.Sprintf("%s=%s", a.Key, a.Val))
			}
			attrStr := strings.Join(kvs, "&")

			if n.PrevSibling != nil {
				if n.Type == n.PrevSibling.Type &&
					n.Data == n.PrevSibling.Data {
					var pKvs []string
					for _, a := range n.PrevSibling.Attr {
						pKvs = append(pKvs, fmt.Sprintf("%s=%s", a.Key, a.Val))
					}
					pAttrStr := strings.Join(pKvs, "&")
					if attrStr == pAttrStr {
						*r += 1
					}
				} else {
					*r = 0
				}
			}
			ns = append(ns, &Node{
				Name:               n.Data,
				Depth:              d,
				SiblingRepeatTimes: *r,
				AttrStr:            attrStr,
			})
		}

		tmp := *r
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				tmp := c
				c = c.NextSibling
				if c == nil {
					break
				}
				c.PrevSibling = tmp.PrevSibling
			}

			f(c, d+1, &tmp)
		}
	}

	var r int
	f(doc, 0, &r)

	return ns, err
}
