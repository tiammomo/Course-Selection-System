package service

import (
	"context"

	"course_select/internal/domain/repository"
)

// ScheduleCourse 排课服务 (二分图最大匹配)
type ScheduleService struct {
	courseRepo repository.ICourseRepo
	bindRepo   repository.IBindRepo
}

// NewScheduleService 创建排课服务
func NewScheduleService(courseRepo repository.ICourseRepo, bindRepo repository.IBindRepo) *ScheduleService {
	return &ScheduleService{
		courseRepo: courseRepo,
		bindRepo:   bindRepo,
	}
}

// Schedule 排课 (二分图最大匹配算法)
func (s *ScheduleService) Schedule(ctx context.Context, teacherPrefs map[string][]string) (map[string]string, error) {
	// 构建二分图
	// 左侧: 教师 (teacherID)
	// 右侧: 课程 (courseID)
	// 边: 教师期望绑定的课程

	// 使用匈牙利算法求解最大匹配

	// 1. 收集所有课程
	allCourses := make(map[string]bool)
	for _, courses := range teacherPrefs {
		for _, courseID := range courses {
			allCourses[courseID] = true
		}
	}

	courseList := make([]string, 0, len(allCourses))
	for courseID := range allCourses {
		courseList = append(courseList, courseID)
	}

	// 2. 构建邻接表 (教师 -> 课程列表)
	adjacency := make(map[string][]string)
	for teacherID, courses := range teacherPrefs {
		adjacency[teacherID] = courses
	}

	// 3. 匈牙利算法
	matchR := make(map[string]string) // course -> teacher
	matchL := make(map[string]string) // teacher -> course

	for teacherID := range adjacency {
		visited := make(map[string]bool)
		bpm(teacherID, adjacency, matchR, visited)
	}

	// 4. 同步 matchL - 遍历 matchR 更新 matchL
	// 因为在 bpm 递归过程中 matchR 会被更新，但 matchL 不会
	for courseID, teacherID := range matchR {
		matchL[teacherID] = courseID
	}

	// 5. 转换为响应格式
	result := make(map[string]string)
	for teacherID, courseID := range matchL {
		result[teacherID] = courseID
	}

	return result, nil
}

// bpm 匈牙利算法核心 (深度优先搜索)
func bpm(u string, adjacency map[string][]string, matchR map[string]string, visited map[string]bool) bool {
	for _, v := range adjacency[u] {
		if visited[v] {
			continue
		}
		visited[v] = true

		// 如果 v 没有被匹配，或者可以重新匹配
		if matchedTeacher, ok := matchR[v]; !ok || bpm(matchedTeacher, adjacency, matchR, visited) {
			matchR[v] = u
			return true
		}
	}
	return false
}

// ValidateSchedule 验证排课结果
func (s *ScheduleService) ValidateSchedule(assignments map[string]string) bool {
	// 检查每个教师是否只分配了一个课程
	courseSet := make(map[string]string) // course -> teacher
	for teacherID, courseID := range assignments {
		if existing, ok := courseSet[courseID]; ok && existing != teacherID {
			return false // 课程被分配给多个教师
		}
		courseSet[courseID] = teacherID
	}
	return true
}

// GetUnassignedTeachers 获取未分配到课程的教师
func (s *ScheduleService) GetUnassignedTeachers(allTeachers []string, assignments map[string]string) []string {
	assigned := make(map[string]bool)
	for _, courseID := range assignments {
		assigned[courseID] = true
	}

	var unassigned []string
	for _, teacherID := range allTeachers {
		if _, ok := assignments[teacherID]; !ok {
			unassigned = append(unassigned, teacherID)
		}
	}

	return unassigned
}
