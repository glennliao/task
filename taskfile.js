use("test", {
    a: " i am use opt a"
})


let a = 1

task("echo")
    .echo("echo something ")


task("withNeed")
    .need("echo")
    .echo("is ok?")

task("chinese")
    .echo("你好hello123")
    .cmd("echo 你好hello123")
    .cmd("echo 你好")
    .cmd("echo 123")


task("n")
    .cmd("echo iam need")

task("cmd")
    .cmd("echo hi" + a)
    .need("n")
    .echo("i am from `echo` in exec command")


task("withCustomStep")
    .echo("1")
    .custom((t) => {
        t.echo("2")
        t.echo("3")
    })


task('default')
    .echo("run ok")


task("mv")
    .mv("example/test.md", "example/test.c")
    .sleep(5)
    .mv("example/test.c", "example/test.md")

task("env")
    .env({
        a: 1,
        b: "bstr"
    })
    .cmd("node example/echoenv.js", {env: {a: "iam1"}})


task("tar")
    .tar("./temp/linux_amd64/main", "./temp/linux_amd64/main.tar.gz", {cwd: ""})
// .tar("./example/use","./test2.tar.gz")


task("ssh")
    .ssh("ubuntu@192.168.64.2", "123456", (ssh, task) => {
        ssh.cmd("ls /")
        ssh.upload("./main.go", "~/main.go")
        let uname = ssh.cmd("uname -a")
        if (uname.indexOf("Ubuntu") !== -1) {
            task.echo("is ubuntu")
        } else {
            task.echo(uname)
        }
        ssh.cmd("ls ~")
        ssh.cmd("du -h ~")
        ssh.cmd("tail -f main.go")
    })

task("glob")
    .echo("asd")
    .glob("ren/*.js.go", (list, task) => {
        list.forEach(item => {
            task.echo("yes: " + item)
            task.mv(item, item.replace(".js.go", ".js"))
        })
    })

task("watch")
    .watch("./taskfile.js", (e, task) => {
        console.log(e.op);
        task.echo(e.name);
    })

task("abort")
    .abort("asd")