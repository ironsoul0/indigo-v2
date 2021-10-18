package moodle

func DetectNewGrades(prevCourses []Course, newCourses []Course) []Course {
	courseToGrades := make(map[string][]Grade)
	for _, course := range prevCourses {
		courseToGrades[course.Name] = course.Grades
	}

	diff := make([]Course, 0)
	for _, course := range newCourses {
		oldGrades := make([]Grade, 0)
		if cur, exists := courseToGrades[course.Name]; exists {
			oldGrades = cur
		}

		uniqueGrades := make([]Grade, 0)
		for _, newGrade := range course.Grades {
			found := false
			for _, grade := range oldGrades {
				if grade == newGrade {
					found = true
				}
			}
			if !found {
				uniqueGrades = append(uniqueGrades, newGrade)
			}
		}

		diff = append(diff, Course{
			Name:   course.Name,
			Grades: uniqueGrades,
		})
	}

	return diff
}
