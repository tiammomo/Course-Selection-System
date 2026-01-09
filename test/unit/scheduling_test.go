package service_test

import (
	"testing"

	"course_select/internal/domain/service"
)

// TestScheduleService_Schedule 测试排课算法
func TestScheduleService_Schedule(t *testing.T) {
	svc := &service.ScheduleService{}

	tests := []struct {
		name    string
		prefs   map[string][]string
		wantErr bool
		// 验证函数用于验证结果是否有效
		validate func(result map[string]string) bool
	}{
		{
			name: "基础测试 - 正常匹配",
			prefs: map[string][]string{
				"a": {"1", "4"},
				"b": {"1", "2"},
				"c": {"2"},
				"d": {"3"},
			},
			wantErr: false,
			validate: func(result map[string]string) bool {
				// 验证每个教师都分配到了他们首选列表中的课程
				for teacherID, courseID := range result {
					prefs, ok := map[string][]string{
						"a": {"1", "4"},
						"b": {"1", "2"},
						"c": {"2"},
						"d": {"3"},
					}[teacherID]
					if !ok {
						return false
					}
					valid := false
					for _, pref := range prefs {
						if pref == courseID {
							valid = true
							break
						}
					}
					if !valid {
						return false
					}
				}
				// 验证没有课程被分配给多个教师
				courseSet := make(map[string]string)
				for teacherID, courseID := range result {
					if existing, ok := courseSet[courseID]; ok && existing != teacherID {
						return false
					}
					courseSet[courseID] = teacherID
				}
				return true
			},
		},
		{
			name: "单教师单课程",
			prefs: map[string][]string{
				"teacher1": {"course1"},
			},
			wantErr: false,
			validate: func(result map[string]string) bool {
				return result["teacher1"] == "course1"
			},
		},
		{
			name: "教师多于课程",
			prefs: map[string][]string{
				"a": {"1"},
				"b": {"1"},
				"c": {"1"},
			},
			wantErr: false,
			validate: func(result map[string]string) bool {
				// 只有一个教师应该被分配到课程
				count := 0
				for _, courseID := range result {
					if courseID == "1" {
						count++
					}
				}
				return count == 1
			},
		},
		{
			name: "课程多于教师",
			prefs: map[string][]string{
				"a": {"1"},
			},
			wantErr: false,
			validate: func(result map[string]string) bool {
				return result["a"] == "1"
			},
		},
		{
			name:    "空偏好",
			prefs:   map[string][]string{},
			wantErr: false,
			validate: func(result map[string]string) bool {
				return len(result) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Schedule(nil, tt.prefs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 验证结果
			if !tt.validate(result) {
				t.Errorf("Schedule() invalid result = %v", result)
			}
		})
	}
}

// TestScheduleService_ValidateSchedule 测试排课结果验证
func TestScheduleService_ValidateSchedule(t *testing.T) {
	svc := &service.ScheduleService{}

	tests := []struct {
		name       string
		assignments map[string]string
		want       bool
	}{
		{
			name: "有效分配",
			assignments: map[string]string{
				"a": "1",
				"b": "2",
			},
			want: true,
		},
		{
			name: "无效分配 - 课程分配给多个教师",
			assignments: map[string]string{
				"a": "1",
				"b": "1",
			},
			want: false,
		},
		{
			name:       "空分配",
			assignments: map[string]string{},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.ValidateSchedule(tt.assignments); got != tt.want {
				t.Errorf("ValidateSchedule() = %v, want %v", got, tt.want)
			}
		})
	}
}
