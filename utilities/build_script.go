package utilities

import "strconv"

func BuildScript(compileCommand string, executeCommand string, timeout int) []byte {
	if compileCommand != "" {
		compileCommand += " && \\\n"
	}
	return []byte("#!/bin/bash\n" +
		compileCommand +
		"timeout -k 0.5s " + strconv.Itoa(timeout) + "s " +
		executeCommand + " < input.txt > output.txt\n" +
		"if [[ $? == 124 || $? == 137 ]]\n" +
		"then\n" +
		"    echo \"timeout: Program execution terminated.\nPlease check for infinite loop.\"\n" +
		"else\n" +
		"    cat output.txt\n" +
		"fi\n" +
		"if [ -f \"output.txt\" ]\n" +
		"then\n" +
		"    rm output.txt\n" +
		"fi\n")
}
