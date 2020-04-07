package stateverifier

//验证器
type Verifier struct {
	items map[string]VerFunc
}

func NewVerifier()*Verifier{
	items:=make(map[string]VerFunc)
	return &Verifier{items:items}
}
//校验执行函数的模板
type VerFunc func(args ...interface{})(bool,interface{})
//注册验证函数
func (t *Verifier)Register(key string,fun VerFunc){
	t.items[key]=fun
}
func(t *Verifier)Delete(key string){
	delete(t.items,key)
}
//根据注册的函数来验证，返回结果与附带的参数
func (t *Verifier)Verify(key string,args ...interface{})(bool,interface{}){
	return t.items[key](args...)
}