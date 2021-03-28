package utilities

func BuildScript(compileCommand string, executeCommand string) []byte {
	if compileCommand != "" {
		compileCommand += " && \\\n"
	}
	return []byte("#!/bin/bash\n" +
		compileCommand +
		"timeout -k 0.5s 1s " + executeCommand + " > output.txt\n" +
		"if [ $? -ge 124 ]\n" +
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
