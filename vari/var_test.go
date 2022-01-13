package vari

import (
  "bytes"
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestSplit(t *testing.T) {
  data := []struct{i string; expect []item}{
    {"aaaa", []item{{con, []byte{0,4}}}},

    {"${aaaa}", []item{{vari, []byte{2,6}}}},

    {"aaa${bbb}ccc", []item{
      {con, []byte{0, 3}},
      {vari, []byte{5, 8}},
      {con, []byte{9, 12}},
      },
    },

    {"aaa${bbb}ccc${ddd}eee", []item{
      {con, []byte{0, 3}},
      {vari, []byte{5, 8}},
      {con, []byte{9, 12}},
      {vari, []byte{14, 17}},
      {con, []byte{18, 21}},
      },
    },

    {"${aaa}${bbb}${ccc}ddd", []item{
      {vari, []byte{2, 5}},
      {vari, []byte{8, 11}},
      {vari, []byte{14, 17}},
      {con, []byte{18, 21}},
      },
    },

    {"${aaa}${bbb}${ccc}", []item{
      {vari, []byte{2, 5}},
      {vari, []byte{8, 11}},
      {vari, []byte{14, 17}},
      },
    },

    {"111${aaa}${bbb}${ccc}", []item{
      {con, []byte{0, 3}},
      {vari, []byte{5, 8}},
      {vari, []byte{11, 14}},
      {vari, []byte{17, 20}},
      },
    },

    {"111${aaa}${bbb}{ccc}", []item{
      {con, []byte{0, 3}},
      {vari, []byte{5, 8}},
      {vari, []byte{11, 14}},
      {con, []byte{15, 20}},
      },
    },
  }

  a := assert.New(t)

  for _,d := range data {
    ret := split([]byte(d.i))
    if len(ret) != len(d.expect) {
      a.EqualValuesf(len(d.expect), len(ret), d.i, "%s len error", d.i)
      //t.Fatalf("%s len error, expect %d, actual %d", d.i, len(d.expect), len(ret))
    }

    for i,it := range ret {
      if it.tp != d.expect[i].tp {
        a.EqualValuesf(d.expect[i].tp, it.tp, "%s %dth type error", d.i, i)
        //t.Fatalf("%s %dth type error, expect %d, actual %d", d.i, i, d.expect[i].tp, it.tp)
      }
      expectStr := []byte(d.i)[d.expect[i].str[0]:d.expect[i].str[1]]

      if !bytes.Equal(it.str, expectStr) {
        a.EqualValuesf(string(expectStr), string(it.str), "%s %dth str error", d.i, i)
        //t.Fatalf("%s %dth str error, expect %s, actual %s", d.i, i, expectStr, it.str)
      }
    }
  }
}

type ob struct {
}

func (o ob) GetVar(name string) string {
  return "+++" + name + "+++"
}

func TestVar_ValueWith(t *testing.T) {
  o := &ob{}

  a := assert.New(t)

  data := []struct{
    t string
    expect string
  } {
    {"im ${fine}", "im +++fine+++"},
    {"${ok}", "+++ok+++"},
    {"you ${are} good", "you +++are+++ good"},
  }

  for _, d := range data {
    v := Compile(d.t)
    a.EqualValues(d.expect, v.ValueWith(o), d.t)
  }
}
