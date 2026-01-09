package dto

// BookCourseRequest 选课请求
type BookCourseRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	CourseID  string `json:"course_id" binding:"required"`
}

// GetStudentCourseRequest 获取学生课表请求
type GetStudentCourseRequest struct {
	StudentID string `json:"student_id" form:"student_id" binding:"required"`
}

// GetStudentCourseResponse 获取学生课表响应
type GetStudentCourseResponse struct {
	CourseList []CourseDTO `json:"course_list"`
}

// CourseDTO 课程数据传输对象
type CourseDTO struct {
	CourseID  string `json:"course_id"`
	Name      string `json:"name"`
	TeacherID string `json:"teacher_id,omitempty"`
}
