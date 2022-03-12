package email

import (
	"testing"

	"github.com/TechMaster/core/config"
	"github.com/TechMaster/core/pmodel"

	// "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

/*
Kiểm thử  chức năng  gửi  email qua Redis Stream và AsynQ
github.com/hibiken/asynq
Hãy chạy ở chế độ debug test
*/
// func Test_Send_Email_Send(t *testing.T) {
// 	config.ReadConfig("..")
// 	asynClient := InitRedisMail()
// 	defer asynClient.Close()
// 	var err error

// 	for i := 0; i < 3; i++ {
// 		err = Emailer.SendPlainEmail([]string{"cuong@techmaster.vn"}, gofakeit.Sentence(10), gofakeit.Paragraph(3, 5, 7, "\n\n"))
// 		if err != nil {
// 			break
// 		}
// 	}
// 	assert := assert.New(t)
// 	assert.Nil(err)
// }

func Test_Send_Email_Marketing(t *testing.T) {
	config.ReadConfig("..")
	asynClient := InitRedisMail()
	defer asynClient.Close()

	var redis_mail RedisMail
	err := redis_mail.SendHTMLEmailMarketing("nhatduc@techmaster.vn", "Nhật Đức", "Đức test",
	[]pmodel.AuthenInfo{{
		UserEmail: "nhatduc@techmaster.vn",
		UserFullName: "nhật đức",
	},{
		UserEmail: "nhatduc.hoanghapaper@gmail.com",
		UserFullName: "nhật đức",
	},
	{
		UserEmail: "huong@techmaster.vn",
		UserFullName: "Hương",
	}}, `
	<div class="cos-module-2-0-0" data-selenium="module-module-2-0-0" data-module-wrapper="true"><div class="js__InlineWrapper-sc-1buus40-0 cnXxRY"><div id="react-tinymce-7" class="TinymceInline__TinyMCEContainer-sc-1qbomzj-0 iywSe hs_cos_wrapper_type_rich_text mce-content-body" contenteditable="true" style="position: relative;"><h1 style="text-align: left; font-size: 15px; line-height: 175%; font-weight: normal;" data-mce-style="text-align: left; font-size: 15px; line-height: 175%; font-weight: normal;"><span style="color: #425b76;" data-mce-style="color: #425b76;">Dear anh Hiệp,</span></h1><p><span style="color: #425b76;" data-mce-style="color: #425b76;">Techmaster thân gửi anh thông tin khóa học:&nbsp;</span></p><p><span style="color: #425b76;" data-mce-style="color: #425b76;">&nbsp;</span></p><h1 style="text-align: center;" data-mce-style="text-align: center;"><span style="color: #425b76;" data-mce-style="color: #425b76;">"LEARN AWS THE HARD WAY"</span></h1><p><br></p><p>Sự khác biệt:&nbsp;</p><div><span style="color: #425b76;" data-mce-style="color: #425b76;"><strong>Nội dung bài học đủ để học viên thi 2 chứng chỉ SAA , DEV và hơn thế nữa</strong></span></div><div><span style="color: #425b76;" data-mce-style="color: #425b76;"><strong>Đặc biệt 120 +++ bài lap với tiêu chí "Làm được việc" bao quát đầy đủ các dịch vụ của AWS</strong></span></div><ul><li><span style="color: #425b76;" data-mce-style="color: #425b76;">Các bài lap dễ: tập trung 1 dịch vụ, hoàn thành trong 10 phút</span></li><li><span style="color: #425b76;" data-mce-style="color: #425b76;">Các bài lap khó: kết hợp nhiều dịch vụ, cần 4h - 8h để hoàn thành</span></li></ul><div><span style="color: #425b76;" data-mce-style="color: #425b76;">Khóa học <strong>30 buổi trong 3 tháng</strong>, tiết kiệm 1 năm so với tự học.
Tăng tỉ lệ thi đỗ chứng chỉ AWS từ 40% lên 80% !</span></div><div><span style="color: #425b76;" data-mce-style="color: #425b76;">&nbsp;</span></div><div style="font-size: 20px; line-height: 175%; text-align: center;" data-mce-style="font-size: 20px; line-height: 175%; text-align: center;"><span style="color: #425b76;" data-mce-style="color: #425b76;"><strong>&gt;&gt;&gt;&nbsp; Ưu đãi <span style="font-size: 30px;" data-mce-style="font-size: 30px;">50%</span> học phí&nbsp; &lt;&lt;&lt;</strong></span></div><div style="font-size: 15px; line-height: 175%;" data-mce-style="font-size: 15px; line-height: 175%;"><span style="color: #425b76;" data-mce-style="color: #425b76;">Áp dụng cho cựu học viên Techmaster các khóa Frontend, Spring Boot, DevOps.</span></div><div style="font-size: 15px; line-height: 175%;" data-mce-style="font-size: 15px; line-height: 175%;"><span style="color: #425b76;" data-mce-style="color: #425b76;">Xem chi tiết nội dung khóa học tại <a href="https://aws.techmaster.vn/#course" style="color: #425b76;" rel="noopener" data-mce-href="https://aws.techmaster.vn/#course" data-mce-style="color: #425b76;">đây</a></span></div><div style="font-size: 15px; line-height: 175%;" data-mce-style="font-size: 15px; line-height: 175%;"><span style="color: #425b76;" data-mce-style="color: #425b76;">&nbsp;</span></div></div></div></div>
	`)
	assert := assert.New(t)
	assert.Nil(err)
}
