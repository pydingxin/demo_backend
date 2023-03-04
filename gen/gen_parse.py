"""
    提取go源码文件中的信息
    parse_command_from_file_lines 从go文件中提取注释的命令
    parse_struct_from_file_lines 从go文件中提取struct定义及其字段
    parse_msg_from_dir 调用前两者，把整个文件夹里的文件都提取一遍
"""
import os
SHOW_DETAILS = False #显示具体的处理过程


def delete_long_comments(s):
    # parse前删除 /*..*/ 类型的注释
    while True:
        left = s.find("/*")
        if -1==left:
            return s
        else:
            right= s.find("*/",left+1)
        s= s[:left]+s[right+2:]


def delete_struct_tags(s):
    # parse前删除 `..`类型的字符串，主要是struct tag
    while True:
        left = s.find("`")
        if -1==left:
            return s
        else:
            right= s.find("`",left+1)
        s= s[:left]+s[right+1:]


def parse_msg_from_dir(dir):
    print("\n"+"*"*80+"\ngen_parse parsing dir: ",dir)
    structs=[]
    cmds = []
    for filename in os.listdir(dir):
        if SHOW_DETAILS:
            print("="*60+"\nprocessing file: ",filename)
        filepath = os.path.join(dir,filename)


        with open(filepath,encoding="utf-8") as fin:
            s=fin.read()
            s=delete_long_comments(s)
            s=delete_struct_tags(s)
            s=s.replace("\t","")
            arr = s.split("\n") #按行处理
            
            lines=[]
            for line in arr:
                lines.append(line.strip())

            cmds.extend(parse_command_from_file_lines(arr))
            structs.extend(parse_struct_from_file_lines(arr)) 
        if SHOW_DETAILS:
            print('='*60)

    print("gen_parse parse dir done")
    if SHOW_DETAILS:
        print("文件夹中所有文件分析结果:")
        print("structs:")
        for struct in structs:
            print(struct)
        print("cmds:")
        for cmd in cmds:
            print(cmd)

    return [structs,cmds]


def parse_command_from_file_lines(arr):
    # 从文件中提取struct命令，一条命令以注释方式单独写一行，写在对应的struct定义之前
    if SHOW_DETAILS:
            print("-"*20+"\nparsing command:")
    all_cmds=[]
    cur_struct_cmds=[] # 遇到下个struct前所有命令存起来，遇到struct后赋给它
    for line in arr:
        if line== "": # 空行
            continue 

        if line.startswith("// CMD_IOMODEL_") or line.startswith("// CMD_DBMODEL_"): # 所有命令有特殊的前缀
            cmd= line.split(" ")[1]
            cur_struct_cmds.append(cmd) # 缓存该命令
            if SHOW_DETAILS:
                print("--"+cmd)
            continue

        if line.startswith("type") and line.endswith("struct {"): # 遇到struct
            struct_name= line.split(" ")[1]
            all_cmds.append([struct_name]+cur_struct_cmds) # 缓存的命令都给这个struct
            cur_struct_cmds=[]   #清空命令缓存
            if SHOW_DETAILS:
                print("+-met struct " + struct_name)
            continue
    return all_cmds


def parse_struct_from_file_lines(arr):
    if SHOW_DETAILS:
            print("-"*20+"\nparsing struct:")
    structs=[]
    cur_struct=[] 
    struct_open=False

    for line in arr:
        if line== "": # 空行
            continue 
        if line.startswith("//"): # 注释行
            continue 

        if line=="}" and struct_open: # struct定义结束
            struct_open=False
            structs.append(cur_struct)
            cur_struct=[]
            continue

        if line.startswith("type") and line.endswith("struct {"): # struct定义开始
            struct_name= line.split(" ")[1]
            if SHOW_DETAILS:
                print("+-met struct " + struct_name)
            struct_open= True
            cur_struct = [struct_name]
            continue

        if struct_open: # struct字段
            fieldName = line.split(" ")[0]
            if SHOW_DETAILS:
                print("--field " + fieldName)
            cur_struct.append(fieldName)
            continue
    return structs
