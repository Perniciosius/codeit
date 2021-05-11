$(function () {
    let editor = CodeMirror.fromTextArea($("codemirror-textarea"), {
        lineNumbers: true,
        mode: "python",
        theme: "material",
        indentWithTabs: true,
        matchBrackets: true,
        smartIndent: true,
    })
    editor.setSize(500, 500);
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