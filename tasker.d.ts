interface SSHClient {
    upload(src: string, dest: string)

    cmd(cmd: string)

    cd(cd: string)
}

declare function sshClient(host: string, password: string): SSHClient


interface CmdOption {
    env: Record<string, string | Number>
}

interface TarOption {
}

interface Task {
    custom(call: (t: Task) => void): Task

    need(taskName: string): Task

    cmd(cmd: string, option?: CmdOption): Task

    rm(path: string): Task

    tar(src: string, dest: string, option?: TarOption): Task

    ssh(host: string, password: string, call: (ssh: SSHClient, task: Task) => void): Task

    echo(msg): Task

    mv(src: string, dest: string): Task

    sleep(second: number): Task

    env(env: Record<string, string | Number>): Task

    glob(pattern: string,call:(list:string[],task:Task)=>void): string[]

    /**
     * abort the task
     * @param msg
     */
    abort(msg: string): void
    watch(pattern:string,call:(e:{name:string,op:string},task:Task)=>void)
}

declare function task(taskName: string): Task

declare function use(useName:string, opt?:Record<string, any>)