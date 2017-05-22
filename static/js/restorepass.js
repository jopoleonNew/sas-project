/**
 * Created by Egor on 21.05.2017.
 */
/**
 * Created by Egor on 17.05.2017.
 */
$(document).ready(function () {
    var currentUrl = window.location.href
    console.log(currentUrl)

    $("#passsubmit").click(function () {
        // console.log($("#userName").val())
        // console.log($("#password").val())
        // console.log($("#email").val())
       // console.log($("#email").val())

        // window.userName = $("#userName").val();
        // window.password = $("#password").val();
       // window.email = $("#email").val();
        var pass1 = $("#password1").val()
        var pass2 = $("#password2").val()
        if (pass1 == "" || pass2 == "") {
            $("#passerror").empty()
            console.log("Enter values")
            $("#passerror").append("<p style='border-radius: 2px; text-align: center; padding: 2px; border: solid 1px;width: 400px; margin: 0 auto; color:red'>Пароль не может быть пустым</p>")
            return
        }
        if (pass1 != pass2) {
            $("#passerror").empty()
            console.log("не равны пароли")
            $("#passerror").append("<p style='border-radius: 2px; text-align: center; padding: 2px; border: solid 1px;width: 400px; margin: 0 auto; color:red'>Пароли не одинаковые</p>")
            return
        }
        if (pass1.length < 6 || pass2.length < 6) {
            $("#passerror").empty()
            console.log("не равны пароли")
            $("#passerror").append("<p style='border-radius: 2px; text-align: center; padding: 2px; border: solid 1px;width: 400px; margin: 0 auto; color:red'>Пароль не может быть меньше 6 символов</p>")
            return
        }

        // window.name = $("#name").val();
        // window.organisation = $("#organisation").val();
        $.ajax({
            data: {
                "pass1": pass1,
                "pass2": pass2
            },
            //dataType: "json",
            type: "POST",
            url: currentUrl,
            success: function (data) {
                // var msg = request.getResponseHeader('X-Message');
                // var type = request.getResponseHeader('X-Message-Type');
                console.log("Data sent: ", data);
                $("#restorepassresults").empty().append(data);

                //$("#emailsubmitresults")
                // if (data.indexOf("Registration is successful. Check your email for activation letter.") >= 0) {
                //     console.log("Inside index Of")
                //     setTimeout(function () {
                //         window.location = currentUrl.replace("signup", "/");
                //     }, 2000);
                //     $("#emailsubmitresults").append(data);
                // }else {
                //     $('#emailsubmitresults').append(
                //         "<p style='color: #DE5246; height: 30px;font-size: 18px'>"+data+"</p>"
                //     )
                // };
            },
            error: function (req, status, err) {
                //console.log(req.responseText)
                console.log(req)

                console.log('Something went wrong', status, err);
                console.log(err)

            }
        });
    });
});