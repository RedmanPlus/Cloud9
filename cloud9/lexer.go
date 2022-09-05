package Cloud9

type Adress struct {
	id      int
	id_line int
}

type Interpreter struct {
	objects        []Object
	all_objects    map[Object]Set[Adress]
	unique_objects Set[Object]
	values         []Value
	operations     []Operation
	commands       []Command
	curr_scope     Scope
	parent_scopes  []Scope
}

func (i *Interpreter) init() {
	i.all_objects = make(map[Object]Set[Adress])
	i.curr_scope = Scope{lvl: 0}
}

func (i Interpreter) get_object(name string) *Object {
	for j := 0; j < len(i.unique_objects.inner); j++ {
		if i.unique_objects.inner[j].name == name {
			return &i.unique_objects.inner[j]
		}
	}

	return &Object{}
}

func (i *Interpreter) gather_objects() {
	var obj_declared      Set[Object]
	for j := 0; j < len(i.objects); j++ {
		obj := i.objects[j]
		if obj.contents != [8]byte{0,0,0,0,0,0,0,0} {
			obj_declared.append(obj)
			res := i.all_objects[obj]
			res.append(Adress{obj.id, obj.id_line})
			i.all_objects[obj] = res
		} else {
			name := obj.name
			for _, inn := range(obj_declared.inner) {
				if inn.name == name {
					old_obj := get(obj_declared.inner, name)
					res := i.all_objects[old_obj]
					res.append(Adress{obj.id, obj.id_line})
					i.all_objects[old_obj] = res
				}
			}
		}

		
	}

	i.unique_objects.inner = obj_declared.inner
}