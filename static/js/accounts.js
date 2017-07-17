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
				};

				console.log("accountloginyandex: ",window.accountlogin)
                console.log("accountroleyandex: ",window.accrole)
                console.log("sourcename: ",window.sourcename)
                if (window.sourcename == "Яндекс Директ") {
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
                    $.ajax({
                        data: {
                            "username": window.currentUser,
                            "accountlogin": window.accountlogin,
                        },
                        //dataType: "json",
                        type: "POST",
                        url: window.location.protocol+"//"+window.location.hostname+":"+window.location.port+"/getauthlink/vkontakte",
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
				// $.ajax({
				// 	data: {
				// 		"sourcename": window.sourcename,
				// 		"accountlogin": window.accountlogin,
				// 		"accrole": window.accrole,
				// 	},
				// 	type: "POST",
                //
				// 	url: currentUrl.replace("accounts", "addaccount"),
				// 	success: function (data) {
				// 		console.log("Data get from server: ", data)
				// 		console.log("accountlogin: ", accountlogin)
                 //        // this part makes ajax request on /getauthcodeyandex
				// 		// endpoint to open YandexDirect oauth.yandex.ru
                 //        // if (window.sourcename == "Яндекс Директ") {
                 //         //    $.ajax({
                 //         //        data: {
                 //         //            "username": window.currentUser,
                 //         //            "accountlogin": window.accountlogin,
                 //         //        },
                 //         //        //dataType: "json",
                 //         //        type: "POST",
                 //         //        url: currentUrl.replace("accounts", "getauthcodeyandex"),
                 //         //        success: function (data) {
                 //         //            //$("#" + appendid).empty();
                 //         //            console.log("Data recieved from getauthcodeyandex: ", data)
                 //         //            var page ="https://oauth.yandex.ru/authorize?response_type=code&client_id=" + data+"&login_hint="+window.accountlogin+"&force_confirm=yes"
                 //         //            window.open(
                 //         //                page,
                 //         //                '_blank' // <- This is what makes it open in a new window.
                 //         //            )
                 //         //            //  console.log("Url to yandex oauth:  ", page)
                 //         //            //  var asdialog = $("<div></div>").html('<iframe style="border: 0px; " src="' + page + '" width="100%" height="100%"></iframe>').dialog({
                 //         //            //          autoOpen: false,
                 //         //            //          modal: true,
                 //         //            //          height: 800,
                 //         //            //          width: 1000,
                 //         //            //          title: "Yandex"
                 //         //            //      });
                 //         //            // asdialog.dialog('open');
                 //         //            // console.log("Id of append obj: ", appendid)
                 //         //        },
                 //         //        error: function (req, status, err) {
                 //         //            //console.log(req.responseText)
                 //         //            console.log(req)
                 //         //            console.log('Something went wrong', status, err);
                 //         //            console.log(err)
                 //         //        }
                 //         //    });
                 //        // };
                //
                 //        // if (window.sourcename == "Вконтакте") {
                 //        //     $.ajax({
                 //        //         data: {
                 //        //             "username": window.currentUser,
                 //        //             "accountlogin": window.accountlogin,
                 //        //         },
                 //        //         //dataType: "json",
                 //        //         type: "POST",
                 //        //         url: currentUrl.replace("accounts", "getauthcodevk"),
                 //        //         success: function (data) {
                 //        //             console.log("Data recieved from getauthcodeyandex: ", data)
                 //        //             var page =data
                 //        //             window.open(
                 //        //                 page,
                 //        //                 '_blank' // <- This is what makes it open in a new window.
                 //        //             )
                 //        //         },
                 //        //         error: function (req, status, err) {
                 //        //             //console.log(req.responseText)
                 //        //             console.log(req)
                 //        //             console.log('Something went wrong', status, err);
                 //        //             console.log(err)
                 //        //         }
                 //        //     });
                 //        // };
                //
                //
                 //            //location.reload()
                //
				// 	},
				// 	error: function (req, status, err) {
				// 		//console.log(req.responseText)
				// 		console.log(req)
                //
				// 		console.log('Something went wrong', status, err);
				// 		console.log(err)
                //
				// 	}
				// });
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
