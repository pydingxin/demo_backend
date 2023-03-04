package do 
import ( 
    dbmodel "demo_backend/model/dbmodel"
    iomodel "demo_backend/model/iomodel"
    response "demo_backend/pkg/response"
    "github.com/gogf/gf/v2/net/ghttp" 
)

func Request_To_RegisterInput(r *ghttp.Request) *iomodel.RegisterInput {
    // 从request表单获取并验证iomodel结构体。 以Request_To_为开头，更符合使用习惯
    var in *iomodel.RegisterInput
    if err := r.Parse(&in); err != nil {
        response.FailMsg(r, err.Error())
        r.ExitAll() // 如果有数据校验有问题，ExitAll会停止handler的执行后续流程
    }
    return in
}
    
func RegisterInput_To_Account(in *iomodel.RegisterInput) *dbmodel.Account {
    // 表单获取的iomodel转为dbmodel 
	return &dbmodel.Account{User: in.User, Pass: in.Pass, }
}    
