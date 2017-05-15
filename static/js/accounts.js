	$(document).ready(function () {
		function jq(myid) {
			return myid.replace(/(:|\.|\[|\]|,|=)/g, "\\$1");
		}
		var currentUrl = window.location.href;
		var hostUrl = window.location.hostname;
		window.currentUser = $("#currentusername").attr("class");
		console.log(window.currentUser)
		console.log(currentUrl)
		console.log(hostUrl)
		console.log(currentUrl + "/addaccount")
		//console.log(currentUrl.replace("login", "signup"))
		//console.log(currentUrl.replace("signup", ""))
		//console.log(String(currentUrl) - "/signup")
		console.log($("#accountloginyandex").val())
		console.log($("#accountloginyoutube").val())
		console.log($("#accountloginvk").val())

		$(".addaccountmodal").each(function () {
			$(this).click(function () {
				var yandexlog = $("#accountloginyandex").val();
				var youtubelog = $("#accountloginyoutube").val();
				var vklog = $("#accountloginvk").val();

				if (yandexlog) {
					window.accountlogin = $("#accountloginyandex").val()
					window.accrole = $("#accountroleyandex").val()
				};
				if (youtubelog) {
					window.accountlogin = $("#accountloginyoutube").val()
				};
				if (vklog) {
					window.accountlogin = $("#accountloginvk").val()
				};

				console.log($(this).attr("id"))
				console.log($("#accountloginyandex").val())
				console.log($("#accountloginyoutube").val())
				console.log($("#accountloginvk").val())
				console.log($("#accountroleyandex").val())
				console.log("sourcename: ", $(this).attr("id"))
				window.sourcename = $(this).attr("id");

				$.ajax({
					data: {
						"sourcename": window.sourcename,
						"accountlogin": window.accountlogin,
						"accrole": window.accrole,
					},
					//dataType: "json",
					type: "POST",
					//http://localhost:3000/accounts?egor/addaccount
					url: currentUrl.replace("accounts", "addaccount"),
					success: function (data) {
						console.log("Data get from server: ", data)
						//location.reload()
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


		$(".getauthcodeyandex").each(function () {
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
						window.open(
							//						[& login_hint=<имя пользователя или электронный адрес>]
							//						[& force_confirm=yes]
							"https://oauth.yandex.ru/authorize?response_type=code&client_id=" + data + "&login_hint=" + accountlogin + "&force_confirm=yes",
							'_blank' // <- This is what makes it open in a new window.
						)
						var sendcodeUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/submityandexcode"
						console.log("hello there: ", sendcodeUrl)

						console.log("Id of append obj: ", appendid)
						//var appID = "#"+ appendid
						//appID = appID.replace(/\./g, '\\\\.');
						//console.log("hello there: ", appID)
						//	username = username.replace(/\./g, '\\\\.');
						//	var testapp = "'div[id='" + appendid +"']'"
						//$('p[id="root.SomeCoolThing"]')
						$("#" + jq(appendid)).append(
							"<br><div>Введите код подтверждения:<input type='text' id='codeinput" + appendid + "'/> <button type='button' class='btn btn-primary btn-xs sendcodeyandex' id='submitcode" + appendid + "'>Отправить</button></div> <div id='coderesult" + appendid + "'></div> "
						);
						//	console.log("After append",$(appID).val())
						$(".sendcodeyandex").each(function () {
							$(this).click(function () {
								console.log("inside of .sendcodeyandex")
								console.log("#codeinput" + jq(appendid))
								console.log("Code inside: ", $("#codeinput" + jq(appendid)).val())
								$.ajax({
									data: {
										"yandexcode": $("#codeinput" + jq(appendid)).val(),
										"accountlogin": window.accountlogin,
									},
									dataType: "json",
									type: "POST",
									url: sendcodeUrl,
									success: function (data) {
										console.log(data)
										//console.log(data.result)
										//console.log(typeof data)
										//console.log(typeof data.result)
										//var result = data.result
										var campaings = data.result
										console.log("typeof campaings: ", typeof campaings)
										console.log("Length of campaings obj: ", campaings.length)
										console.log("Append element ID: ", $("#coderesult" + appendid))

										for (var i = 0; i < campaings.length; i++) {
											$("#coderesult" + appendid).append("<div class='campaingdiv' id='" + campaings[i].Id + "' style='display: inline-block; border-radius: 6px; text-align: left; width: 30%; padding: 10px; margin: 5px; font-size: 18px; border: solid 1px; background-color: white'>" + "<p>Имя кампании: " + campaings[i].Name + "</p>" + "<p>Номер(ID) кампаниии: " + campaings[i].Id + "</p>" + "</div>"

											);
											console.log(campaings[i])
										}
										//$("#coderesult"+appendid).append(data)

									},
									error: function (req, status, err) {
										console.log(req)
										console.log('Something went wrong', status, err);
										//	console.log(err)
									}
								});

							});
						});
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

		//	$(".gotostatisticpage").each(function() {
		//		$(this).click(function (){
		//			
		//			accountlogin = $(this).attr("name");
		//			//appendid = $(this).attr("result")	
		//		});
		//	});

		
	});
