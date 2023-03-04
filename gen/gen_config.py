"""
    本模块定义各种常量、路径设置
"""

import os

def read_module_name_from_gomod():
    """从当前项目的go.mod获取包名"""
    modulename = ''
    with open(MAIN_CWD+"/go.mod") as fin:
        modulename= fin.readline().split(" ")[1]
    return  modulename.strip()

def do_file_path(structName):
    # 对struct XX，返回do_XX的路径，以生成其文件
    return  tidy_path_sep(os.path.join(MAIN_CWD, "model/do", f"do_{structName}.go"))

def do_file_path_of_iomodel(structName):
    # 对struct XX，返回do_XX的路径，以生成其文件
    return  tidy_path_sep(os.path.join(MAIN_CWD, "model/do", f"iomodel_{structName}.go"))

def do_file_path_of_dbmodel(structName):
    # 对struct XX，返回do_XX的路径，以生成其文件
    return  tidy_path_sep(os.path.join(MAIN_CWD, "model/do", f"dbmodel_{structName}.go"))

def tidy_path_sep(s):
    """统一path中的正反斜杠"""
    return s.replace("/",os.sep).replace("\\",os.sep)

## 在vscode里，控制台pwd强制为根目录，所以当前脚本是在根目录执行的
# MAIN_CWD 应为 C:\Users\dingxin\Desktop\-gen-，这是后续路径拼接的前提
MAIN_CWD = os.path.abspath(os.curdir)
print(f"MAIN_CWD = [{MAIN_CWD}]") 

GOMOD_NAME = read_module_name_from_gomod() #go.mod中的包名
print(f"GOMOD_NAME = [{GOMOD_NAME}]")

GEN_DIR_PATH = tidy_path_sep(MAIN_CWD+"/gen") # /gen目录

DB_MODEL_DIR_PATH = tidy_path_sep(MAIN_CWD+"/model/dbmodel") #不能用os.path.join拼接，join这个函数也有执行路径
IO_MODEL_DIR_PATH = tidy_path_sep(MAIN_CWD+"/model/iomodel")
DO_CODE_DIR_PATH = tidy_path_sep(MAIN_CWD+"/model/do")

# in gen_comand:
ORM_SEMANTICS_FILE_PATH = tidy_path_sep(GEN_DIR_PATH+"/dbmodel_orm_semantics_code_sample") #orm模板文件在gen目录中
ORM_SEMANTICS_STRUCT_NAME = "OBJECTMODEL"       #模板文件中struct的名字
ORM_SEMANTICS_MODULE_NAME = "GOMOD_MODULE_NAME" #模板文件中包的名字