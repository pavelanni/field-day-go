package main

func saveVisitor(v Visitor) error {
	return db.Save(&v)
}

func listVisitors() ([]Visitor, error) {
	var visitors []Visitor
	err := db.AllByIndex("ID", &visitors)
	if err != nil {
		return nil, err
	}
	return visitors, nil
}
