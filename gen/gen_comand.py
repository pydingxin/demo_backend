import gen_config
import os
import json
#--------------------------------------------------------------------------------------------------------
# 全局变量
SHOW_DETAILS = True #显示具体的处理过程
all_structs_map = {}

#--------------------------------------------------------------------------------------------------------
# 工具函数
def p(s):
    if not SHOW_DETAILS:
        return
    print(s)    

def assure_iomodel_generatedFilePath_exists(generatedFilePath):
    """
        这里只确保文件存在
        由于iomodel命令较多，先执行各个命令写入具体代码，最后检测一遍各个文件，根据情况写入头部
    """
    if(os.path.exists(generatedFilePath)):
        return
    p("create file: "+generatedFilePath)
    with open(generatedFilePath,'w',encoding="utf-8") as fout:
        fout.write("")

def append_content_to_file(path, content):
    with open(path,'a+',encoding="utf-8") as fout:
        fout.write(content)

def remove_do_files():
    # 删除 /model/do/ 文件夹下以前生成的所有文件，以便重新生成
    dirpath = gen_config.DO_CODE_DIR_PATH
    p(f"\ngen_comand removing files in {dirpath} :")
    for file in os.listdir(dirpath):
        filepath = os.path.join(dirpath,file)
        p(filepath)
        os.remove(filepath)

#--------------------------------------------------------------------------------------------------------
# struct信息存取
def save_structs_to_file(structs):
    """保存所有解析到的structs到文件里"""
    filepath = os.path.join(gen_config.GEN_DIR_PATH,"_tmp_parsed_structs_.json")
    with open(filepath,'w') as f:
        p("\ngen_command save structs to <<_tmp_parsed_structs_.json>>\n")
        json.dump(structs, f,ensure_ascii=False,indent=2)

def load_structs_from_file():
    """提取文件里的structs"""
    filepath = os.path.join(gen_config.GEN_DIR_PATH,"_tmp_parsed_structs_.json")
    with open(filepath,'r') as f:
        p("\ngen_command load structs from <<_tmp_parsed_structs_.json>>\n")
        return json.loads(f.read())

def assure_all_structs_loaded():
    if len(all_structs_map)>0:
        return
    for struct in load_structs_from_file():
        all_structs_map[struct[0]] = struct

#--------------------------------------------------------------------------------------------------------
def execute_commands(structCmds):
    # 执行文件夹中收集到的命令集，每条cmd开头是struct名字，后续是对其使用的所有命令
    if len(structCmds)>0:
        p(f"in execute_commands, {len(structCmds)} struct has comands in the dir: ")
        for cmd in structCmds:
            p(cmd)
        p("")

    idx=1 # print command index
    for cmds in structCmds:
        struct = cmds[0]
        for cmd in cmds[1:]:
            p("~"*60+f"\n#{idx}--[{struct}]-[{cmd}]")
            idx+=1

            # 创建dbmodel orm封装文件
            if cmd =="CMD_DBMODEL_MAKE_ORM_SEMANTIC_FILE": 
                execute_command_make_dbmodel_orm_semantic_file(struct)

            # 从request表单获取并验证iomodel结构体
            if cmd =="CMD_IOMODEL_PARSE_FROM_REQUEST": 
                execute_command_make_iomodel_parse_request_file(struct)

            # 表单中获取的iomodel转为dbmodel
            if cmd.startswith("CMD_IOMODEL_CONVERT_TO_DBMODEL_"):
                toStruct = cmd.replace("CMD_IOMODEL_CONVERT_TO_DBMODEL_","")
                execute_command_iomodel_convert_to_dbmodel(struct, toStruct)

            # 查询dbmodel库表生成iomodel子结构体
            if cmd.startswith("CMD_IOMODEL_QUERY_FROM_DBMODEL_"):
                fromStruct = cmd.replace("CMD_IOMODEL_QUERY_FROM_DBMODEL_","")
                execute_command_append_dbmodel_query_from_code(struct, fromStruct)

#--------------------------------------------------------------------------------------------------------
def shared_fields_of_two_struct_models(m1,m2):
    # 获取dbmodel和iomodel共有的字段，方便execute_command_iomodel_convert_to_dbmodels生成代码
    assure_all_structs_loaded()
    fields1 = all_structs_map[m1][1:] #第一个是struct名字
    fields2 = all_structs_map[m2][1:]
    return list(set(fields1) & set(fields2))


def execute_command_iomodel_convert_to_dbmodel(struct, toStruct):
    # 表单获取的iomodel转为dbmodel  
    print("生成代码——表单获取的iomodel转为dbmodel。"+ struct+ " → "+ toStruct)
    common_fields = shared_fields_of_two_struct_models(struct, toStruct)
    print("共有字段：",common_fields)
    inner_code=""
    for field in common_fields:
        inner_code+= (f"{field}: in.{field}, ")

    code="""
func IN_IOMODEL_To_TODBMODEL(in *iomodel.IN_IOMODEL) *dbmodel.TODBMODEL {
    // 表单获取的iomodel转为dbmodel 
	return &dbmodel.TODBMODEL{INNER_CODE}
}    
""".replace("IN_IOMODEL",struct).\
    replace("TODBMODEL",toStruct).\
    replace("INNER_CODE",inner_code)
    dst_file_path = gen_config.do_file_path_of_iomodel(struct)
    append_content_to_file(dst_file_path,code)
    print("添加函数 "+ f"{struct}_To_{toStruct}()"+ "到文件："+ dst_file_path)
#--------------------------------------------------------------------------------------------------------
def execute_command_append_dbmodel_query_from_code(subsetStruct, fromStruct):
    """
        查询dbmodel库表生成iomodel子结构体
        添加函数 dbmd_QueryMulti_subset_ByFields()
    """ 
    code = """
func OBJECTMODEL_QueryMulti_SUBSET_ByFields(withFields *dbmodel.OBJECTMODEL) *[]iomodel.SUBSET {
    // 查询dbmodel库表生成iomodel子结构体。
    // 使用gorm高级查询的智能选择字段，select返回类型所拥有的字段，而不select所有字段
    var queryResult []iomodel.SUBSET
    result := orm.Conn().Model(withFields).Where(withFields).Find(&queryResult)
    orm.PanicGormResultError("do.OBJECTMODEL_QueryMulti_SUBSET_ByFields", result)
    return &queryResult
}
    """.replace("OBJECTMODEL",fromStruct).replace("SUBSET",subsetStruct)
    path = gen_config.do_file_path_of_iomodel(subsetStruct)
    assure_iomodel_generatedFilePath_exists(path)
    append_content_to_file(path,code)

    appended_funcname = "OBJECTMODEL_QueryMulti_SUBSET_ByFields" \
        .replace("OBJECTMODEL",fromStruct).replace("SUBSET",subsetStruct)
    print(f"生成代码——查询dbmodel库表生成iomodel子结构体。\n添加函数 {appended_funcname}() 到文件 {path}")

#--------------------------------------------------------------------------------------------------------
def execute_command_make_iomodel_parse_request_file(struct):
    """
        从request表单获取并验证iomodel结构体
    """ 
   
    content="""
func Request_To_OBJECTMODEL(r *ghttp.Request) *iomodel.OBJECTMODEL {
    // 从request表单获取并验证iomodel结构体。 以Request_To_为开头，更符合使用习惯
    var in *iomodel.OBJECTMODEL
    if err := r.Parse(&in); err != nil {
        response.FailMsg(r, err.Error())
        r.ExitAll() // 如果有数据校验有问题，ExitAll会停止handler的执行后续流程
    }
    return in
}
    """
    #要生成的文件路径
    path = gen_config.do_file_path_of_iomodel(struct)
    assure_iomodel_generatedFilePath_exists(path)
    p(f"生成代码——从request表单获取并验证iomodel结构体。\n添加函数 ParseFromRequestTo_{struct}() 到文件："+path)
    append_content_to_file(path,content.replace("OBJECTMODEL",struct))



#--------------------------------------------------------------------------------------------------------
def execute_command_make_dbmodel_orm_semantic_file(structName):
    """创建dbmodel orm封装文件"""
    generatedFilePath = gen_config.do_file_path_of_dbmodel(structName) # 要生成的文件路径
    p("生成代码——创建dbmodel orm封装文件: "+generatedFilePath)

    with open(gen_config.ORM_SEMANTICS_FILE_PATH,'r',encoding="utf-8") as fin:
        # 读取模板文件，替换掉项目包名和struct名字
        s=fin.read()\
            .replace(gen_config.ORM_SEMANTICS_STRUCT_NAME, structName)\
            .replace(gen_config.ORM_SEMANTICS_MODULE_NAME, gen_config.GOMOD_NAME)
        # 处理后的字符串写入到目标文件中。这些文件都是第一次生成，直接用w模式写入
        with open(generatedFilePath,'w',encoding="utf-8") as fout:
            fout.write(s)

#--------------------------------------------------------------------------------------------------------
def get_file_content(filepath):
    with open(filepath,'r',encoding="utf-8") as fin:
        return fin.read()
    
def put_file_content(filepath,content):
    with open(filepath,'w',encoding="utf-8") as fout:
        fout.write(content)

def append_imports_to_iomodel_do_files():
    # 为所有do/iomodel_XX 文件添加package & imports头部
    goMod=gen_config.GOMOD_NAME
    dirpath = gen_config.DO_CODE_DIR_PATH
    for file in os.listdir(dirpath):
        if file.startswith("dbmodel_"):
            continue
        filepath= os.path.join(dirpath, file)
        content= get_file_content(filepath)
        innerCode=""

        if content.find("dbmodel.")>=0:
            innerCode+= f'\n    dbmodel "{goMod}/model/dbmodel"'
        if content.find("iomodel.")>=0:
            innerCode+= f'\n    iomodel "{goMod}/model/iomodel"'
        if content.find("response.")>=0:
            innerCode+= f'\n    response "{goMod}/pkg/response"'
        if content.find("orm.")>=0:
            innerCode+= f'\n    orm "{goMod}/model/orm"'
        if content.find("ghttp.")>=0:
            innerCode+= f'\n    "github.com/gogf/gf/v2/net/ghttp"'

        head = f"package do \nimport ( {innerCode} \n)\n"
        put_file_content(filepath, head+content)