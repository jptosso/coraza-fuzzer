package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/cheggaaa/pb/v3"
	fuzzer "github.com/jptosso/coraza-fuzzer"
	engine "github.com/jptosso/coraza-waf"
	"github.com/jptosso/coraza-waf/operators"
	t "github.com/jptosso/coraza-waf/transformations"
)

func main() {
	//routines := 10
	cpath := "fuzzer.yml"
	config, err := fuzzer.ReadConfig(cpath)
	if err != nil {
		panic(err)
	}
	f := config.Fuzzer
	if f.Iterations == 0 {
		f.Iterations = int64(^uint(0) >> 1)
	}
	total := int64(len(f.Transformations)+len(f.Operators)) * f.Iterations * int64(len(f.Rules))
	bar := pb.Full.Start64(total)
	tmap := t.TransformationsMap()
	rules := []fuzzer.Rule{}
	for _, r := range f.Rules {
		rules = append(rules, fuzzer.Rules[r])
	}
	fmt.Printf("Loaded %d tests\n", total)
	fn := ""
	tp := ""
	data := ""
	args := ""
	tools := &t.Tools{}
	waf := engine.NewWaf()
	defer func() {
		if r := recover(); r != nil {
			p, err := ioutil.TempFile("/tmp", "fuzzer*")
			if err != nil {
				fmt.Println("Unable to write tmp file")
			}
			p.WriteString(data)
			p.Close()
			fmt.Printf("%s failed, output written to %s\n", fn, p.Name())
			fmt.Println(r)
		}
		// now we replicate the error
		if tp == "trans" {
			tmap[fn](data, nil)
		} else {
			tx := waf.NewTransaction()
			f := operators.OperatorsMap()[fn]
			f.Init(args)
			f.Evaluate(tx, data)
		}
	}()
	for _, trans := range f.Transformations {
		tp = "trans"
		tf := tmap[trans]
		fn = trans
		if tf == nil {
			panic("invalid transformation " + fn)
		}
		for i := int64(0); i < f.Iterations; i++ {
			size := rand.Uint64() % uint64(f.MaxLength)
			for _, rule := range rules {
				data = rule(size)
				tf(data, tools)
				bar.Increment()
			}
		}
	}
	for _, op := range f.Operators {
		tp = "op"
		opf := operators.OperatorsMap()[op.Name]
		fn = op.Name
		if opf == nil {
			panic("invalid operator " + fn)
		}
		for _, arg := range op.Args {
			args = arg
			opf.Init(arg)
			for i := int64(0); i < f.Iterations; i++ {
				size := rand.Uint64() % uint64(f.MaxLength)
				tx := waf.NewTransaction()
				for _, rule := range rules {
					data = rule(size)
					opf.Evaluate(tx, data)
					bar.Increment()
				}
			}
		}
	}
	bar.Finish()
}
