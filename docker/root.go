package docker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func DockerInit() {
	if err := os.Chdir("docker"); err != nil {
		panic("Failed To cd Docker Directory")
	} else {
		cmd := exec.Command("docker-compose", "down")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Run(); err != nil {
			panic("Failed To Down Docker")
		} else {
			cmd = exec.Command("docker-compose", "up", "-d")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				panic("Failed To Up Docker")
			} else {
				cmd = exec.Command("docker", "ps")

				// 명령어 실행 및 표준 출력 파이프 가져오기
				if stdout, err := cmd.StdoutPipe(); err != nil {
					panic("Failed To Get PipeLine")
				} else if err = cmd.Start(); err != nil {
					panic("Failed To Start Docker ps")
				} else {
					scanner := bufio.NewScanner(stdout)
					for scanner.Scan() {
						line := scanner.Text()
						if !strings.Contains(line, "CONTAINER ID") {
							log.Println("Success To Create Container")
							log.Println(line)
							fmt.Println()
						}
					}
				}
			}
		}
	}

}
