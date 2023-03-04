import os


import gen_parse
import gen_config 
import gen_comand
if __name__ == "__main__":

    gen_comand.remove_do_files()

    # 先处理dbmodel，产生库表操作代码
    [dbstructs, dbcmds] =gen_parse.parse_msg_from_dir(gen_config.DB_MODEL_DIR_PATH) 
    gen_comand.execute_commands(dbcmds)

    # 再处理iomodel，这是库表操作的延伸
    [iostructs, iocmds] = gen_parse.parse_msg_from_dir(gen_config.IO_MODEL_DIR_PATH)
    gen_comand.save_structs_to_file(dbstructs+iostructs) # 保存结构，执行iocmds时读取使用
    gen_comand.execute_commands(iocmds)

    gen_comand.append_imports_to_iomodel_do_files() # 生成do/iomodel_文件后，添加头部