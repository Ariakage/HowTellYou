import os
workdir = os.getcwd() + "/src/"
filelist = os.listdir(workdir)
print(filelist)
codelen = 0
for i in filelist:
    fileSplit = i.split(".")
    if len(fileSplit) > 0:
        if fileSplit[len(fileSplit)-1] in ["go", "sql", "toml"]:
            with open(file=workdir+i, encoding="utf-8") as f:
                #print(str(f.read()).split("\n"))
                codelen += len(f.read().split("\n"))
print(codelen)
