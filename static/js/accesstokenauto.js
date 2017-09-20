$(document).ready(function () {
	var currentUrl = window.location.href;
	var hostUrl = window.location.hostname;
	window.currentUser = $("#currentusername").attr("class");
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
					var page ="https://oauth.yandex.ru/authorize?response_type=code&client_id=" + data+"&login_hint="+accountlogin+"&force_confirm=yes"
                    var $dialog = $('<div></div>').html('<iframe style="border: 0px; " src="' + page + '" width="100%" height="100%"></iframe>')
                        .dialog({
                            autoOpen: false,
                            modal: true,
                            height: 800,
                            width: 1000,
                            title: "Yandex"
                        });
                    $dialog.dialog('open');
					console.log("Id of append obj: ", appendid)
				},
				error: function (req, status, err) {
					console.log(req)
					console.log('Something went wrong', status, err);
					console.log(err)
				}
			});
		});
	});
});