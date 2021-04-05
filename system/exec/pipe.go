package exec

import (
	"bytes"
	"os/exec"
)

// 管道
func Pipe(commands ...*exec.Cmd) (out string, outerr string, err error) {
	commandLen := len(commands)

	if commandLen == 0 {
		return "", "", nil
	}

	if commandLen == 1 {
		command := commands[0]

		var cmdout bytes.Buffer
		var cmderr bytes.Buffer
		command.Stdout = &cmdout
		command.Stderr = &cmderr

		commandErr := command.Run()

		return cmdout.String(), cmderr.String(), commandErr
	}

	if commandLen > 1 {
		command := commands[0]

		var cmdout bytes.Buffer
		var cmderr bytes.Buffer

		// 将 command 的 stdout 进行连接
		for i := 1; i < commandLen; i++ {
			commands[i].Stdin, _ = commands[i-1].StdoutPipe()
		}

		// 设置 最后一个 命令的输出
		commands[commandLen-1].Stdout = &cmdout
		commands[commandLen-1].Stderr = &cmderr

		// 开始监听， 越靠后的命令先 Start
		for i := commandLen; i > 1; i-- {
			_ = commands[i-1].Start()
		}

		// 执行第一个命令
		commandErr := command.Run()

		//for i := commandLen; i > 1; i-- {
		//	_ = commands[i - 1].Wait()
		//}

		for i := 1; i < commandLen; i++ {
			_ = commands[i].Wait()
		}

		return cmdout.String(), cmderr.String(), commandErr
	}

	return "", "", nil
}
