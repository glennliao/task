task("ssh")
    .ssh("pi@192.168.31.70:22", "", (ssh) => {
            ssh.upload("./main", "/home/pi/test/main")
            ssh.cd("~/test/yl")
            ssh.cmd("ls -h")
        }
    )