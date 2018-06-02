package ctrl

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func CourseList(c echo.Context) error {
	rows, err := DB.Query("select course, teacher from courses")
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	defer rows.Close()

	courses := make(map[string]string)
	for rows.Next() {
		var course string
		var teacher string

		err := rows.Scan(&course, &teacher)
		if err != nil {
			fmt.Printf("database error: %s\n", err.Error())
			return err
		}
		courses[course] = teacher
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	return c.JSON(http.StatusOK, courses)
}

func CourseAdd(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	course := m["course"].(string)
	teacher := m["teacher"].(string)

	query := fmt.Sprintf(
		`insert into courses(course, teacher) values ("%s", "%s")`,
		course, teacher,
	)
	_, err := DB.Exec(query)
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}

	return c.String(http.StatusOK, "success")
}

func StudentAdd(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	course := m["course"].(string)
	student := m["student"].(string)

	query := fmt.Sprintf(
		`insert into classes(course, student) values ("%s", "%s")`,
		course, student,
	)
	_, err := DB.Exec(query)
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}

	return c.String(http.StatusOK, "success")
}

func StudentList(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	course := m["course"].(string)

	query := fmt.Sprintf(`select student from classes where course="%s"`, course)
	rows, err := DB.Query(query)
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	defer rows.Close()

	students := []string{}
	for rows.Next() {
		var student string

		err := rows.Scan(&student)
		if err != nil {
			fmt.Printf("database error: %s\n", err.Error())
			return err
		}
		students = append(students, student)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	return c.JSON(http.StatusOK, students)
}

func StudentAbsent(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	course := m["course"].(string)
	student := m["student"].(string)
	memo := m["memo"].(string)

	query := fmt.Sprintf(
		`insert into absentees(course, student, memo) values ("%s", "%s", "%s")`,
		course, student, memo,
	)
	_, err := DB.Exec(query)
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}

	return c.String(http.StatusOK, "success")
}

func AbsentList(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	course := m["course"].(string)

	absentees := make(map[string][]string)
	query := fmt.Sprintf(`select student, memo from absentees where course="%s"`, course)
	rows, err := DB.Query(query)
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var student string
		var memo string

		err := rows.Scan(&student, &memo)
		if err != nil {
			fmt.Printf("database error: %s\n", err.Error())
			return err
		}
		absentees[student] = append(absentees[student], memo)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("database error: %s\n", err.Error())
		return err
	}
	return c.JSON(http.StatusOK, absentees)
}
