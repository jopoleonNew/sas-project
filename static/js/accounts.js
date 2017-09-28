	$(document).ready(function () {

		function jq(myid) {
			return myid.replace(/(:|\.|\[|\]|,|=)/g, "\\$1");
		}
		var currentUrl = window.location.href;
		var hostUrl = window.location.hostname;
		window.currentUser = $("#currentusername").attr("class");
		console.log(window.currentUser);
		console.log(currentUrl);
		console.log(hostUrl);
		console.log(currentUrl + "/addaccount");
		console.log(window.location.protocol+"//"+window.location.hostname+":"+window.location.port+"/addaccount/yandex");
		$(".addaccountmodal").each(function () {
			$(this).click(function () {
				var yandexlog = $("#accountloginyandex").val();
				var youtubelog = $("#accountloginyoutube").val();
				var vklog = $("#accountloginvk").val();
                var adwords = $("#accountloginadwords").val();

				if (yandexlog) {
					window.accountlogin = $("#accountloginyandex").val();
					window.accrole = $("#accountroleyandex").val();
                    window.sourcename = $(this).attr("id")
				} else if (window.accrole === 0 || window.sourcename === 0) {
                    	console.log("Some of important valuse are empty:");
						console.log("accountloginyandex: ", $("#accountloginyandex").val());
						console.log("accountroleyandex: ", $("#accountroleyandex").val());
						return;
				}
                if (vklog) {
                    console.log("Новый ВК аккаунт")
                    window.accountlogin = $("#accountloginvk").val()
                    window.sourcename = $(this).attr("id")
                };
				if (youtubelog) {
					window.accountlogin = $("#accountloginyoutube").val()
                    window.sourcename = $(this).attr("id")
				};
                if (adwords) {
                    window.accountlogin = $("#accountloginadwords").val()
                    window.sourcename = $(this).attr("id")
                };


                console.log("sourcename: ",window.sourcename)
                if (window.sourcename == "Яндекс Директ") {
                    console.log("accountloginyandex: ",window.accountlogin)
                    console.log("accountroleyandex: ",window.accrole)
                    $.ajax({
                        data: {
                            "username": window.currentUser,
                            "accountlogin": window.accountlogin,
							"accrole": window.accrole,
                        },
                        //dataType: "json",
                        type: "POST",
                        url: window.location.protocol+"//"+window.location.hostname+":"+window.location.port+"/getauthlink/yandex",
                        success: function (data) {
                            window.open(
                                data,
                            '_blank' // <- This is what makes it open in a new window.
                             )
                        },
                        error: function (req, status, err) {
                            //console.log(req.responseText)
                            console.log(req)
                            console.log('Something went wrong', status, err);
                            console.log(err)
                        }
                    });
                };
                if (window.sourcename == "Вконтакте") {
                    console.log("Creating new VK account")
                    $.ajax({
                        data: {
                            "username": window.currentUser,
                            "accountlogin": window.accountlogin,
                        },
                        //dataType: "json",
                        type: "POST",
                        url: window.location.protocol+"//"+window.location.hostname+":"+window.location.port+"/getauthlink/vkontakte",
                        success: function (data) {
                            console.log(data);
                            window.open(
                                data,
								'_blank' // <- This is what makes it open in a new window.
                                       )
                        },
                        error: function (req, status, err) {
                            //console.log(req.responseText)
                            console.log(req)
                            console.log('Something went wrong', status, err);
                            console.log(err)
                        }
                    });
                };
                if (window.sourcename == "AdWords") {
                    console.log("Creating new AdWords account")
                    $.ajax({
                        data: {
                            "username": window.currentUser,
                            "accountlogin": window.accountlogin,
                        },
                        //dataType: "json",
                        type: "POST",
                        url: window.location.protocol+"//"+window.location.hostname+":"+window.location.port+"/getauthlink/adwords",
                        success: function (data) {
                            console.log(data);
                            window.open(
                                data,
                                '_blank' // <- This is what makes it open in a new window.
                            )
                        },
                        error: function (req, status, err) {
                            //console.log(req.responseText)
                            console.log(req)
                            console.log('Something went wrong', status, err);
                            console.log(err)
                        }
                    });
                };
			})
        });


		$(".deleteaccountbutton").each(function () {
			$(this).click(function () {
				//console.log($(this).attr("id"))
				console.log("deleteaccountbutton", $(this).attr("id"))
				window.accountlogin = $(this).attr("id");
				$.ajax({
					data: {
						"accountlogin": window.accountlogin,
					},
					//dataType: "json",
					type: "POST",
					//http://localhost:3000/accounts?egor/addaccount
					url: currentUrl.replace("accounts", "deleteaccount"),
					success: function (data) {
						console.log("Data sent: ", data)
						location.reload()
					},
					error: function (req, status, err) {
						//console.log(req.responseText)
						console.log(req)

						console.log('Something went wrong', status, err);
						console.log(err)

					}
				});
			})
		});
	});
