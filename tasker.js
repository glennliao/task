function task(taskName) {
    _tasker.new(taskName)


    let methods = {
        custom(call) {
            call(this)
            return this
        },

        need(task) {
            _tasker.need(task)
            return this
        },

        cmd(cmdStr, option = {}) {
            _tasker.cmd(cmdStr, JSON.stringify(option))
            return this
        },


        rm(path) {
            _tasker.rm(path)
            return this
        },

        tar(src, dest, option = {}) {
            _tasker.tar(src, dest, JSON.stringify(option))
            return this
        },

        ssh(host, password, call) {
            let client = _tasker.ssh(host, password)
            call(client, this)
            return this
        },


        echo(msg) {
            _tasker.echo(msg)
            return this
        },
        mv(src, dest) {
            _tasker.mv(src, dest)
            return this
        },
        sleep(sec) {
            _tasker.sleep(sec)
            return this
        },
        env(env) {
            _tasker.env(JSON.stringify(env))
            return this
        },
        glob(pattern, call) {
            let list = _tasker.glob(pattern)
            call(list, this)
            return this
        },
        abort(msg) {
            _tasker.abort(msg)
        },
        watch(pattern, call) {
            _tasker.watch(pattern,(e)=>{
                call(e,this)
            })
        }
    }

    if (taskName !== currentTask) {
        let _methods = {}
        Object.keys(methods).forEach(k => {
            _methods[k] = () => _methods
        })
        return _methods
    }

    return methods
}

function use(name, opt = {}) {
    _tasker.use(name, opt)
}