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
	}}, `
	<div style="mso-line-height-rule:exactly; font-size:20px; line-height:175%; text-align:center" align="center"><span style="color: #425b76;"><strong>&gt;&gt;&gt;&nbsp; Ưu đãi <span style="font-size: 30px;">50%</span> học phí&nbsp; &lt;&lt;&lt;</strong></span></div>
	<div style="mso-line-height-rule:exactly; font-size:15px; line-height:175%"><span style="color: #425b76;">Áp dụng cho cựu học viên Techmaster các khóa Frontend, Spring Boot, DevOps.</span></div>
	<div style="mso-line-height-rule:exactly; font-size:15px; line-height:175%"><span style="color: #425b76;">Xem chi tiết nội dung khóa học tại <a href="https://aws.techmaster.vn/?utm_source=hs_email&amp;utm_medium=email&amp;_hsenc=p2ANqtz--QG9wePQg7xwb7omz2Wj9CGn1JnHfL058f5aMzxhkRlsEcqLxPPfRrwkdSJ5CBWcDb6w6F#course" style="mso-line-height-rule:exactly; color:#425b76" rel="noopener" data-hs-link-id="0" target="_blank">đây</a></span></div>
	<div style="mso-line-height-rule:exactly; font-size:15px; line-height:175%"><span style="color: #425b76;">&nbsp;</span></div></div></div></td></tr></tbody></table>
	</div>
	
	
	
	<div id="column-2-0" class="hse-column hse-size-12"> <table role="presentation" cellpadding="0" cellspacing="0" width="100%" style="border-spacing:0 !important; border-collapse:collapse; mso-table-lspace:0pt; mso-table-rspace:0pt"><tbody><tr><td class="hs_padded" style="border-collapse:collapse; mso-line-height-rule:exactly; font-family:Arial, sans-serif; font-size:15px; color:#23496d; word-break:break-word; padding:10px 20px"><div id="hs_cos_wrapper_module-2-0-0" class="hs_cos_wrapper hs_cos_wrapper_widget hs_cos_wrapper_type_module" style="color: inherit; font-size: inherit; line-height: inherit;" data-hs-cos-general-type="widget" data-hs-cos-type="module"><div id="hs_cos_wrapper_module-2-0-0_" class="hs_cos_wrapper hs_cos_wrapper_widget hs_cos_wrapper_type_rich_text" style="color: inherit; font-size: inherit; line-height: inherit;" data-hs-cos-general-type="widget" data-hs-cos-type="rich_text"><h1 style="margin:0; mso-line-height-rule:exactly; text-align:left; font-size:15px; line-height:175%; font-weight:normal" align="left"><span style="color: #425b76;">Dear anh {{.Name}},&nbsp;</span></h1> <p style="mso-line-height-rule:exactly; line-height:175%"><span style="color: #425b76;">Techmaster thân gửi anh thông tin khóa học mới:&nbsp;</span></p> <p style="mso-line-height-rule:exactly; line-height:175%"><span style="color: #425b76;">&nbsp;</span></p> <h1 style="margin:0; mso-line-height-rule:exactly; line-height:175%; font-size:24px; text-align:center" align="center"><span style="color: #425b76;">"LEARN AWS THE HARD WAY"</span></h1> <p style="mso-line-height-rule:exactly; line-height:175%">&nbsp;</p> <div style="mso-line-height-rule:exactly; line-height:175%"><span style="color: #425b76;"><strong>Nội dung đủ cover 2 chứng chỉ SAA và DEV</strong></span></div> <div style="mso-line-height-rule:exactly; line-height:175%"><span style="color: #425b76;"><strong>Đặc biệt 120 +++ bài lap với tiêu chí "Làm được việc" bao quát đầy đủ các dịch vụ của AWS</strong></span></div> <ul style="mso-line-height-rule:exactly; line-height:175%"> <li style="mso-line-height-rule:exactly"><span style="color: #425b76;">Các bài lap dễ: tập trung 1 dịch vụ, hoàn thành trong 10 phút</span></li> <li style="mso-line-height-rule:exactly"><span style="color: #425b76;">Các bài lap khó: kết hợp nhiều dịch vụ, cần 4h - 8h để hoàn thành</span></li> </ul> <div style="mso-line-height-rule:exactly; line-height:175%"><span style="color: #425b76;">Khóa học <strong>30 buổi trong 3 tháng</strong>, tiết kiệm 1 năm so với tự học. Tăng tỉ lệ thi đỗ chứng chỉ AWS từ 40% lên 80% !</span></div> <div style="mso-line-height-rule:exactly; line-height:175%"><span style="color: #425b76;">&nbsp;</span></div>
	`)
	assert := assert.New(t)
	assert.Nil(err)
}
