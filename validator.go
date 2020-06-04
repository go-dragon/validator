package validator

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Validator struct {
	HasErr  bool              // check if validate has error
	ErrList map[string]string // error params list
}

type Rules map[string]string

// new validator
func New() *Validator {
	return &Validator{
		HasErr:  false,
		ErrList: make(map[string]string),
	}
}

// form直接传bingo.MergeRequest中处理返回的键值对请求数据
// rules like:, map[string]string{
//	"param1": "numeric|min:10",
//	"param1": "numeric|min:10",
//}
func (v *Validator) Validate(form *map[string]string, rules Rules) *Validator {
	v.HasErr = false //单个字段校验结果是否有错
	for field, rule := range rules {
		valis := strings.Split(rule, "|")
		for _, vali := range valis {
			if vali == "" {
				continue
			}
			subValis := strings.Split(vali, ":")
			//这里只支持一个参数
			arg := ""
			method := vali
			//如果子校验规则带有参数则再做处理
			if len(subValis) == 2 {
				method = subValis[0]
				arg = subValis[1]
			}
			if method == "notEmpty" {
				if v.notEmpty(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "不能为空"
					continue
				}
			}
			if method == "mobile" {
				if v.mobile(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "格式错误"
					continue
				}
			}
			if method == "password" {
				if v.password(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "格式错误"
					continue
				}
			}
			if method == "nick" {
				if v.nick(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "格式错误"
					continue
				}
			}
			if method == "regex" {
				if v.regex(field, form, arg) == false {
					v.HasErr = true
					v.ErrList[field] = "格式错误"
					continue
				}
			}
			if method == "min" {
				if v.min(field, form, arg) == false {
					v.HasErr = true
					v.ErrList[field] = "值小于规定值:" + arg
					continue
				}
			}

			// max:3
			if method == "max" {
				if v.max(field, form, arg) == false {
					v.HasErr = true
					v.ErrList[field] = "大于规定值:" + arg
					continue
				}
			}

			// minLength
			if method == "minLength" {
				if v.minLength(field, form, arg) == false {
					v.HasErr = true
					v.ErrList[field] = "小于最小长度:" + arg
					continue
				}
			}

			// maxLength
			if method == "maxLength" {
				if v.maxLength(field, form, arg) == false {
					v.HasErr = true
					v.ErrList[field] = "超过最大长度:" + arg
					continue
				}
			}

			// numeric
			if method == "numeric" {
				if v.numeric(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "非数值型"
					continue
				}
			}

			if method == "int64" {
				if v.int64(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "非64位整型"
					continue
				}
			}

			if method == "int32" {
				if v.int32(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "非32位整型"
					continue
				}
			}

			if method == "int" {
				if v.int32(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "非32位整型"
					continue
				}
			}

			if method == "datetime" {
				if v.datetime(field, form) == false {
					v.HasErr = true
					v.ErrList[field] = "格式不正确, 正确格式需要: 2006-01-02 15:04:05"
					continue
				}
			}

			if method == "in" {
				if v.in(field, form, arg) == false {
					v.HasErr = true
					v.ErrList[field] = "参数不正确，没有在规定范围内"
					continue
				}
			}

			if method == "notIn" {
				if v.notIn(field, form, arg) == false {
					v.HasErr = true
					v.ErrList[field] = "参数不正确，没有在规定范围内"
					continue
				}
			}
		}
	}
	return v
}

// notEmpty that means param not exist or param is null character
func (*Validator) notEmpty(field string, form *map[string]string) bool {
	v, ok := (*form)[field]
	if !ok || v == "" {
		return false
	}
	return true
}

// chinese mobile validate
func (*Validator) mobile(field string, form *map[string]string) bool {
	v, _ := (*form)[field]
	if ok, _ := regexp.MatchString("^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\\d{8}$", v); ok {
		return true
	}
	return false
}

// password validate 密码验证 密码8-16位数字和字母的组合这两个符号(不能是纯数字或者纯字母)
func (*Validator) password(field string, form *map[string]string) bool {
	v, _ := (*form)[field]
	if ok, _ := regexp.MatchString("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{8,16}$", v); ok {
		return true
	}
	return false
}

// user nick validate 用户昵称校验 中文和英文或数字不能有特殊符号长度为2-10位
func (*Validator) nick(field string, form *map[string]string) bool {
	v, _ := (*form)[field]
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9\u4e00-\u9fff]{2,10}$", v); ok {
		return true
	}
	return false
}

// regex validate
func (*Validator) regex(field string, form *map[string]string, pattern string) bool {
	v, _ := (*form)[field]
	if ok, _ := regexp.MatchString(pattern, v); ok {
		return true
	}
	return false
}

// numeric param min number param >= arg
func (*Validator) min(field string, form *map[string]string, arg string) bool {
	v, _ := (*form)[field]
	number, err := strconv.Atoi(v)
	min, errV := strconv.Atoi(arg)
	if err != nil || errV != nil {
		return false
	}
	if number >= min {
		return true
	}
	return false
}

// numeric param max number param <= arg
func (*Validator) max(field string, form *map[string]string, arg string) bool {
	v, _ := (*form)[field]
	number, err := strconv.Atoi(v)
	max, errV := strconv.Atoi(arg)
	if err != nil || errV != nil {
		return false
	}
	if number <= max {
		return true
	}
	return false
}

// numeric
func (*Validator) numeric(field string, form *map[string]string) bool {
	v, _ := (*form)[field]
	_, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false
	}
	return true
}

// param maxLength
func (*Validator) maxLength(field string, form *map[string]string, arg string) bool {
	v, _ := (*form)[field]
	maxLength, err := strconv.Atoi(arg)
	if err != nil {
		return false
	}
	if utf8.RuneCountInString(v) <= maxLength {
		return true
	}
	return false
}

// param minLength
func (*Validator) minLength(field string, form *map[string]string, arg string) bool {
	v, _ := (*form)[field]
	minLength, err := strconv.Atoi(arg)
	if err != nil {
		return false
	}
	if utf8.RuneCountInString(v) >= minLength {
		return true
	}
	return false
}

func (*Validator) int64(field string, form *map[string]string) bool {
	v, _ := (*form)[field]
	_, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return false
	}
	return true
}

func (*Validator) int32(field string, form *map[string]string) bool {
	v, _ := (*form)[field]
	_, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return false
	}
	return true
}

// date time validate
func (*Validator) datetime(field string, form *map[string]string) bool {
	v, _ := (*form)[field]
	_, err := time.Parse("2006-01-02 15:04:05", v)
	if err != nil {
		return false
	}
	return true
}

// in:1,2,3
func (*Validator) in(field string, form *map[string]string, arg string) bool {
	v, _ := (*form)[field]
	args := strings.Split(arg, ",")
	for _, str := range args {
		if v == str {
			// 如果校验的数据值是和规则匹配的，返回true
			return true
		}
	}
	return false
}

// notIn:1,2,3
func (v *Validator) notIn(field string, form *map[string]string, arg string) bool {
	return !v.in(field, form, arg)
}
