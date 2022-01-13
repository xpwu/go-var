package vari

import (
  "regexp"
)

type VarObject interface {
  GetVar(name string) string
  //SetVar(name string, value string)
}

type value func(object VarObject) string

type Var struct {
  values []value
}

func (v *Var) ValueWith(object VarObject) string {
  ret := ""
  for _,val := range v.values {
    ret += val(object)
  }

  return ret
}

func constValue(name string) value {
  return func(object VarObject) string {
    return name
  }
}

func varValue(name string) value {
  return func(object VarObject) string {
    return object.GetVar(name)
  }
}

// ${ccc}=eeee : aaaa${ccc}dddd => aaaaeeeedddd

func Compile(varName string) *Var {
  items := split([]byte(varName))
  ret := &Var{
    make([]value, 0, len(items)),
  }

  for _,it := range items {
    switch it.tp {
    case con:
      ret.values = append(ret.values, constValue(string(it.str)))
    case vari:
      ret.values = append(ret.values, varValue(string(it.str)))
    }
  }

  return ret
}

const (
  con = iota
  vari
)

type item struct {
  tp  int
  str []byte
}

func split(varName []byte) []item {
  expr := "\\$\\{(.*?)\\}"
  reg := regexp.MustCompile(expr)
  result := reg.FindAllSubmatchIndex(varName, -1)

  ret := make([]item, 0, 5)

  if result == nil {
    ret = append(ret, item{con, varName})
    return ret
  }

  start := 0
  for _,sub := range result {
    _ = sub[3] // check bounds

    if start != sub[0] { // 两个子串之间有其他字符
      ret = append(ret, item{con, varName[start:sub[0]]})
    }
    start = sub[1] // 整个匹配的结束点 也就是 } 的下一个位置
    ret = append(ret, item{vari, varName[sub[2]:sub[3]]})
  }

  if start < len(varName) {
    ret = append(ret, item{con, varName[start:]})
  }

  return ret
}

