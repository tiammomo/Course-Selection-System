package service_test

import (
	"testing"

	"course_select/internal/domain/model"
	"course_select/internal/pkg/errcode"
)

// TestCreateMemberRequest_Validate 测试创建成员请求验证
func TestCreateMemberRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     model.CreateMemberRequest
		wantErr bool
	}{
		{
			name: "有效的请求",
			req: model.CreateMemberRequest{
				Username:  "testuser",
				Password:  "Password1",
				Nickname:  "TestUser",
				UserType:  model.UserTypeStudent,
			},
			wantErr: false,
		},
		{
			name: "无效的用户类型",
			req: model.CreateMemberRequest{
				Username:  "testuser",
				Password:  "Password1",
				Nickname:  "TestUser",
				UserType:  99,
			},
			wantErr: true,
		},
		{
			name: "管理员类型",
			req: model.CreateMemberRequest{
				Username:  "admin",
				Password:  "AdminPass1",
				Nickname:  "Admin",
				UserType:  model.UserTypeAdmin,
			},
			wantErr: false,
		},
		{
			name: "教师类型",
			req: model.CreateMemberRequest{
				Username:  "teacher",
				Password:  "TeacherPass1",
				Nickname:  "Teacher",
				UserType:  model.UserTypeTeacher,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCreateCourseRequest_Validate 测试创建课程请求验证
func TestCreateCourseRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     model.CreateCourseRequest
		wantErr bool
	}{
		{
			name: "有效的课程请求",
			req: model.CreateCourseRequest{
				Name: "高等数学",
				Cap:  100,
			},
			wantErr: false,
		},
		{
			name: "无效的容量 - 0",
			req: model.CreateCourseRequest{
				Name: "无效课程",
				Cap:  0,
			},
			wantErr: true,
		},
		{
			name: "无效的容量 - 负数",
			req: model.CreateCourseRequest{
				Name: "无效课程",
				Cap:  -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestMember_ToResponse 测试成员响应转换
func TestMember_ToResponse(t *testing.T) {
	member := &model.Member{
		UserID:    1,
		Username:  "testuser",
		Nickname:  "Test User",
		UserType:  model.UserTypeStudent,
		IsDeleted: false,
	}

	response := member.ToResponse()

	if response.UserID != "1" {
		t.Errorf("UserID = %s, want 1", response.UserID)
	}
	if response.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", response.Username)
	}
	if response.Nickname != "Test User" {
		t.Errorf("Nickname = %s, want Test User", response.Nickname)
	}
	if response.UserType != model.UserTypeStudent {
		t.Errorf("UserType = %d, want %d", response.UserType, model.UserTypeStudent)
	}
}

// TestMember_ToResponse_Nil 测试空成员转换
func TestMember_ToResponse_Nil(t *testing.T) {
	var member *model.Member = nil
	response := member.ToResponse()

	if response != nil {
		t.Errorf("Expected nil response, got %v", response)
	}
}

// TestCourse_ToResponse 测试课程响应转换
func TestCourse_ToResponse(t *testing.T) {
	teacherID := 10
	course := &model.Course{
		CourseID:    1,
		Name:        "高等数学",
		Capacity:    100,
		CapSelected: 50,
		TeacherID:   &teacherID,
	}

	response := course.ToResponse()

	if response.CourseID != "1" {
		t.Errorf("CourseID = %s, want 1", response.CourseID)
	}
	if response.Name != "高等数学" {
		t.Errorf("Name = %s, want 高等数学", response.Name)
	}
	if response.Capacity != 100 {
		t.Errorf("Capacity = %d, want 100", response.Capacity)
	}
	if response.TeacherID != "10" {
		t.Errorf("TeacherID = %s, want 10", response.TeacherID)
	}
}

// TestErrCode_Error 测试错误码
func TestErrCode_Error(t *testing.T) {
	tests := []struct {
		name string
		code errcode.ErrCode
		want string
	}{
		{
			name: "成功",
			code: errcode.OK,
			want: "success",
		},
		{
			name: "参数错误",
			code: errcode.ParamInvalid,
			want: "参数不合法",
		},
		{
			name: "用户不存在",
			code: errcode.UserNotExisted,
			want: "用户不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
