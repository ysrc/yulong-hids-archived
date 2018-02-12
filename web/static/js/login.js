
$(document).ready(function () {

    $("#loginbtn").click(function () {
        user = $.trim($("#username").val());
        pass = $.trim($("#password").val());
        if (user == "" || pass == "") {
            swal("输入错误!", "用户名及密码不可为空", "warning");
        } else {
            $.post(location.href,
                {
                    username: user,
                    password: pass
                },
                function (res, status) {
                    if (status == "success") {
                        if (res.status) {
                            location.href = "/#!/";
                        } else {
                            swal("输入错误!", "用户名或密码验证错误", "error");
                        }
                    }
                }
            );
        }
    });

    $("input").keyup(function (event) {
        if (event.keyCode === 13) {
            $("#loginbtn").click();
        }
    });

});
