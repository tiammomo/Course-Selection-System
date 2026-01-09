package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"course_select/internal/domain/model"
	"course_select/internal/domain/service"
	"course_select/internal/pkg/errcode"
	"course_select/internal/pkg/response"
)

// MemberHandler 成员处理器
type MemberHandler struct {
	memberService *service.MemberService
}

// NewMemberHandler 创建成员处理器
func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

// CreateMember 创建成员
// @Summary 创建成员
// @Description 管理员创建新成员
// @Tags member
// @Accept json
// @Produce json
// @Param request body model.CreateMemberRequest true "创建成员请求"
// @Success 200 {object} response.Response
// @Router /member/create [post]
func (h *MemberHandler) CreateMember(c *gin.Context) {
	var req model.CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	// 验证请求
	if err := req.Validate(); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	member, err := h.memberService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(map[string]string{
		"user_id": strconv.Itoa(member.UserID),
	}))
}

// GetMember 获取成员信息
// @Summary 获取成员信息
// @Description 根据用户ID获取成员信息
// @Tags member
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} response.Response
// @Router /member [get]
func (h *MemberHandler) GetMember(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg("user_id 不能为空")))
		return
	}

	member, err := h.memberService.Get(c.Request.Context(), userID)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(member.ToResponse()))
}

// GetMemberList 获取成员列表
// @Summary 获取成员列表
// @Description 分页获取成员列表
// @Tags member
// @Produce json
// @Param offset query int false "偏移量"
// @Param limit query int false "限制数量"
// @Success 200 {object} response.Response
// @Router /member/list [get]
func (h *MemberHandler) GetMemberList(c *gin.Context) {
	offset := parseIntSafe(c.DefaultQuery("offset", "0"), 0)
	limit := parseIntSafe(c.DefaultQuery("limit", "20"), 20)

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	members, _, err := h.memberService.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	var memberList []model.MemberResponse
	for _, m := range members {
		memberList = append(memberList, *m.ToResponse())
	}

	c.JSON(200, response.Success(map[string]interface{}{
		"member_list": memberList,
	}))
}

// UpdateMember 更新成员信息
// @Summary 更新成员信息
// @Description 更新成员昵称
// @Tags member
// @Accept json
// @Produce json
// @Param request body model.UpdateMemberRequest true "更新成员请求"
// @Success 200 {object} response.Response
// @Router /member/update [post]
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	var req model.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	if err := h.memberService.Update(c.Request.Context(), req.UserID, req.Nickname); err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(nil))
}

// DeleteMember 删除成员
// @Summary 删除成员
// @Description 软删除成员
// @Tags member
// @Accept json
// @Produce json
// @Param request body model.DeleteMemberRequest true "删除成员请求"
// @Success 200 {object} response.Response
// @Router /member/delete [post]
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	var req model.DeleteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	if err := h.memberService.Delete(c.Request.Context(), req.UserID); err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(nil))
}

// parseIntSafe 安全解析字符串为 int，失败返回默认值
func parseIntSafe(s string, defaultVal int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return defaultVal
}
