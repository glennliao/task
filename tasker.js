function task(taskName){
    _tasker.new(taskName)

    function sshClient(steps){
        return {
            upload(src,dest){
                steps.push({
                    op:"upload",
                    args:[src,dest]
                })
                return this
            },
            cmd(str){
                steps.push({
                    op:"cmd",
                    args:[str]
                })
                return this
            },
            cd(str){
                steps.push({
                    op:"cd",
                    args:[str]
                })
                return this
            }
            // download(src,dest){
            //     ssh.scpDown(src,dest)
            //     return this
            // }
        }
    }

    return {

        custom(call){
            call(this)
            return this
        },

        need(task){
            _tasker.need(task)
            return this
        },

        cmd(cmdStr){
            _tasker.cmd(cmdStr)
            return this
        },


        rm(path){
            _tasker.rm(path)
            return this
        },

        tar(src , dest){
            _tasker.tar(src,dest)
            return this
        },

        ssh(host, password, call){
            let sshSteps = []

            let client = sshClient(sshSteps)
            call(client)

            _tasker.ssh(host, password,JSON.stringify(sshSteps))
            return this
        },
        test(test){
            _tasker.test(test)
            return this
        },

        echo(msg){
            _tasker.echo(msg)
            return this
        }
    }
}

function use(name, opt={}){
    _tasker.use(name,opt)
}