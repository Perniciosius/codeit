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
		"    echo \"timeout: Program execution terminated.\n" +
		"Your program is taking too much time to execute.\n" +
		"Please check for infinite loop or some other reasons.\"\n" +
		"else\n" +
		"	if [ -f \"output.txt\" ]\n" +
		"	then\n" +
		"		cat output.txt\n" +
		"	fi\n" +
		"fi\n" +
		"if [ -f \"output.txt\" ]\n" +
		"then\n" +
		"    rm output.txt\n" +
		"fi\n")
}
