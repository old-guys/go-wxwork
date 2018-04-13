package concerns

import (
	"github.com/jinzhu/gorm"
	"fmt"
	//"wxwork/initializers"
)

type Tree struct {}

func (user *Tree) BeforeSave(scope *gorm.Scope) (err error) {

	fmt.Println(user)

	return
}

func (v *Tree) AssignAttributes(attrs map[string]interface{}) {
	//fmt.Println("vvvv", v)
	//for column, value := range attrs {
	//	initializers.DB.NewScope(v).SetColumn(column, value)
	//}
}