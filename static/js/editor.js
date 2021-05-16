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
        let textval = editor.getValue()
        textval = textval.replace(/\n/g, "\r\n")
        let lang = $("#language").val()
        let textFileAsBlob = new Blob([textval], { type: 'text/HTML' })
        let filename = ''
        switch (lang) {
            case "c":
                filename = "code.c"
                break
            case "c++":
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
        downloadLink.innerHTML = "LINKTITLE";
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
        $.post(`http://ec2-52-66-228-90.ap-south-1.compute.amazonaws.com/api/${lang}`, $(this).serialize()).done(function (data) {
            $("#output").val(data.output)
        }).fail(function () {
            $("#output").val("Some error has occured!\nTry again later.....")
        })
    })
})


