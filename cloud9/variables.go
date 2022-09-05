package Cloud9

func (i *Interpreter) gather_value(id int, id_line int) {
	
	var ops []Operation

	for _, o := range(i.operations) {
		if o.id_line == id_line {
			ops = append(ops, o)
		}
	}

	var obj_val [8]byte
	var err     error

	if len(ops) == 1 {

		for _, v := range(i.values) {
			if v.id == id+3 && v.id_line == id_line {
				obj_val = v.contents
			}
		}

		for _, o := range(i.objects) {
			if o.id == id+3 && o.id_line == id_line {
				if o.name == "true" {
					obj_val = str_to_bytes("1")
				} else if o.name == "false" {
					obj_val = str_to_bytes("0")
				} else {
					obj_val = o.contents
				}
			}
		}

	} else {
		obj_val, err = i.compute_operation(id+4, id_line)
		handle(err)
	}

	
	for j := 0; j < len(i.objects); j++ {
		if i.objects[j].id == id && i.objects[j].id_line == id_line {
			i.objects[j].contents = obj_val
		}
	}
}