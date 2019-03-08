package questions

// // Store stores a question locally if the local bool is true
// func (q *Question) Store(local bool) error {
// 	if local {
// 		return q.storeLocal()
// 	}
// 	return errors.New("Not storing question because local is false")
// }

// // Remove removes a specific question from the local questions
// func (q *Question) Remove() error {
// 	var path = fmt.Sprintf("./questions/questions/")
// 	err := os.Remove(fmt.Sprintf("%s%s.txt", path, q.Link))
// 	return err
// }

// // RemoveAllQuestions removes all questions inside of the questions folder
// func RemoveAllQuestions() {
// 	qs := GetLocalQuestions()
// 	for _, q := range qs {
// 		q.Remove()
// 	}
// }
