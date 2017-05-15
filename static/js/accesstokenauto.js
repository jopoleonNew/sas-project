$(document).ready(function () {


	var currentUrl = window.location.href;
	var hostUrl = window.location.hostname;
	//	var currentuser = document.getElementById("currentusername").className;
	window.currentUser = $("#currentusername").attr("class");
	//console.log(currentuser)
	//console.log(window.currentUser)
	//console.log(currentUrl)
	//console.log(hostUrl)



	//console.log(currentUrl.replace("login", "signup"))
	//console.log(currentUrl.replace("signup", ""))
	//console.log(String(currentUrl) - "/signup")
//	console.log(currentUrl + "/addaccount")
	
	$(".getaccesstokenauto").each(function () {
		$(this).click(function () {

			accountlogin = $(this).attr("name");
			appendid = $(this).attr("result")

			console.log("accountlogin: ", accountlogin)
			$.ajax({
				data: {
					"username": window.currentUser,
					"accountlogin": accountlogin,
				},
				//dataType: "json",
				type: "POST",
				url: currentUrl.replace("accounts", "getauthcodeyandex"),
				success: function (data) {
					$("#" + appendid).empty();
					console.log("Data recieved: ", data)
					//$('#getauthcodeyandexresult').empty()
						//$('#getauthcodeyandexresult').append(data);
						//console.log("Data sent: ", data)
					var page =
                            "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + data+"&login_hint="+accountlogin+"&force_confirm=yes"
                    var $dialog = $('<div></div>')
                        .html('<iframe style="border: 0px; " src="' + page + '" width="100%" height="100%"></iframe>')
                        .dialog({
                            autoOpen: false,
                            modal: true,
                            height: 800,
                            width: 1000,
                            title: "Yandex"
                        });
                    $dialog.dialog('open');
					// window.open(
					// "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + data+"&login_hint="+accountlogin+"&force_confirm=yes",
					// 	'_blank' // <- This is what makes it open in a new window.
					// )
//					var sendcodeUrl = window.location.protocol + "//" + window.location.hostname +":"+window.location.port + "/submityandexcode"
//					console.log("hello there: ", sendcodeUrl)
					console.log("Id of append obj: ", appendid)
				
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
});