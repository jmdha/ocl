package wowlogs

import (
	"fmt"
	"strconv"
)

type fieldParser struct {
	fields []string
	pos    int
	err    error
}

func (p *fieldParser) skip() {
	if p.err != nil {
		return
	}
	if p.pos >= len(p.fields) {
		p.err = fmt.Errorf("field %d: out of range", p.pos)
		return
	}
	p.pos++
}

func (p *fieldParser) String() string {
	if p.err != nil || p.pos >= len(p.fields) {
		p.err = fmt.Errorf("field %d: out of range", p.pos)
		return ""
	}
	v := p.fields[p.pos]
	p.pos++
	return v
}

func (p *fieldParser) StringSkip() string {
	p.skip()
	return p.String()
}

func (p *fieldParser) Int() int {
	s := p.String()
	if p.err != nil {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		p.err = fmt.Errorf("field %d: %w", p.pos-1, err)
	}
	return v
}

func (p *fieldParser) IntSkip() int {
	p.skip()
	return p.Int()
}

func (p *fieldParser) Float() float64 {
	s := p.String()
	if p.err != nil {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		p.err = fmt.Errorf("field %d: %w", p.pos-1, err)
	}
	return v
}

func (p *fieldParser) Bool() bool {
	s := p.String()
	if p.err != nil {
		return false
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		p.err = fmt.Errorf("field %d: %w", p.pos-1, err)
	}
	return v
}

func (p *fieldParser) BoolSkip() bool {
	p.skip()
	return p.Bool()
}
