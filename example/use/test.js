task("test")
    .echo(" i am from use")

task("test:env")
    .echo(`i am from use with opt ${opt.a}`)