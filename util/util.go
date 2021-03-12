package util

func Arrayfilter(ss [] interface {}, test func(interface{}) bool) (ret []interface{}) {
    for _, s := range ss {
        if test(s) {
            ret = append(ret, s)
        }
    }
    return
}
