
// use("test",{
//     a:" i am use opt a"
// })


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

task("cmd")
    .cmd("echo hi")
    .echo("i am from `echo` in exec command")


task("withCustomStep")
    .echo("1")
    .custom((t)=>{
        t.echo("2")
        t.echo("3")
    })


task('default')
    .echo("run ok")



