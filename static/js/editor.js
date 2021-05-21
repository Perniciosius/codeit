let downloadCode
$(function () {
    let editor = CodeMirror.fromTextArea(document.getElementById("codemirror-textarea"), {
        lineNumbers: true,
        mode: "python",
        theme: "material",
        indentWithTabs: true,
        matchBrackets: true,
        smartIndent: true,
    })
    downloadCode = function () {
        let code = editor.getValue()
        code = code.replace(/\n/g, "\r\n")
        let lang = $("#language").val()
        let textFileAsBlob = new Blob([code], { type: 'text/HTML' })
        let filename = ''
        switch (lang) {
            case "c":
                filename = "code.c"
                break
            case "cpp":
                filename = "code.cpp"
                break
            case "java":
                filename = "code.java"
                break
            case "python2":
            case "python3":
                filename = "code.py"
                break
            case "go":
                filename = "code.go"
                break
            case "javascript":
                filename = "code.js"
                break
            case "typescript":
                filename = "code.ts"
                break
        }
        let downloadLink = document.createElement("a");
        downloadLink.download = filename;
        window.URL = window.URL || window.webkitURL;
        downloadLink.href = window.URL.createObjectURL(textFileAsBlob);
        downloadLink.style.display = "none";
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);
    }
    $("form").submit(function (event) {
        let lang = $("#language").val()
        event.preventDefault()
        $.post(`/api/${lang}`, $(this).serialize()).done(function (data) {
            $("#output").val(data.output)
        }).fail(function () {
            $("#output").val("Some error has occured!\nTry again later.....")
        })
    })
})


