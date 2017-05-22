/**
 * Created by Egor on 17.05.2017.
 */
$(document).ready(function () {
    var currentUrl = window.location.href
    console.log(currentUrl)
    $("#emailsubmit").click(function () {
        // console.log($("#userName").val())
        // console.log($("#password").val())
        // console.log($("#email").val())
        console.log($("#email").val())

        // window.userName = $("#userName").val();
        // window.password = $("#password").val();
        window.email = $("#email").val();
        // window.name = $("#name").val();
        // window.organisation = $("#organisation").val();

        $.ajax({
            data: {
                "email": window.email
            },
            //dataType: "json",
            type: "POST",
            url: currentUrl,
            success: function (data) {
                console.log("Data sent: ", data);
                $("#emailsubmitresults").empty().append(data);

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