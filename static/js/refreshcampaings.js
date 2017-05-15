
$(document).ready(function () {


	var currentUrl = window.location.href;
	var hostUrl = window.location.hostname;
	//	var currentuser = document.getElementById("currentusername").className;
	window.currentUser = $("#currentusername").attr("class");
	//getcampaingstatistic
	$(".refreshdbcampaign").each(function () {
		$(this).click(function () {

			accountlogin = $(this).attr("name");
			appendid = $(this).attr("result")
			console.log("Id of append obj: ", appendid)
			console.log("accountlogin: ", accountlogin)
			$.ajax({
				data: {
					"username": window.currentUser,
					"accountlogin": accountlogin,
				},
				//dataType: "json",
				type: "POST",
				url: currentUrl.replace("accounts", "refreshdbcampaign"),
				success: function (data) {
					$("#" + appendid).empty();
					console.log("Data recieved: ", data)
					//$('#getauthcodeyandexresult').empty()
					$("#"+appendid).append(data);
						//console.log("Data sent: ", data)



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
	$(".getcampaingstatistic").each(function () {
		$(this).click(function () {

			campaingId = $(this).attr("id");
			accountlogin = $(this).attr("name");
			appendid = $(this).attr("result")
			console.log("Id of append obj: ", appendid)
			console.log("campaingId: ", campaingId)
			$.ajax({
				data: {
					"username": window.currentUser,
					"accountlogin": accountlogin,
					"campaingId": campaingId,
				},
				dataType: "json",
				type: "POST",
				url: currentUrl.replace("accounts", "getcampaingstats"),
				success: function (data) {
					var objsd1 = data.data
					if (objsd1.length == 0) {
							console.log("data.data.length = 0 getcampaingstatistic Data recieved: ", objsd1)
						$("#"+appendid).append("<p style='color:orange'>Statisc is empty</p>")
					};
					if (objsd1.length == 1){
							$("#" + appendid).empty();
					console.log("data.data.length = 1 getcampaingstatistic Data recieved: ", data)
					//$('#getauthcodeyandexresult').empty()
					var objsd= data.data[0];
					console.log("objsd ", data.data[0])
					$("#"+appendid).append(
						"<p>Количество кликов в Рекламной сети Яндекса.: "+objsd.ClicksContext+"</p>"+
						"<p>Количество кликов на поиске.: "+objsd.ClicksSearch+"</p>"+
						"<p>Глубина просмотра сайта при переходе из Рекламной сети Яндекса.: "+objsd.SessionDepthContext+"</p>"+
						"<p>Глубина просмотра сайта при переходе с поиска.: "+objsd.SessionDepthSearch+"</p>"+
						"<p>Количество показов в Рекламной сети Яндекса.: "+objsd.ShowsContext+"</p>"+
						"<p>Количество показов на поиске.: "+objsd.ShowsSearch+"</p>"+
						"<p>Дата, за которую приведена статистика.: "+objsd.StatDate+"</p>"+
						"<p>Стоимость кликов в Рекламной сети Яндекса.: "+(objsd.SumContext/1000000)+"</p>"+
						"<p>Доля целевых визитов в общем числе визитов при переходе с поиска, в процентах..: "+objsd.GoalConversionSearch+"</p>"+
						"<p>Доля целевых визитов в общем числе визитов при переходе из Рекламной сети Яндекса, в процентах: "+objsd.GoalConversionContext+"</p>"+
						"<p>Цена достижения <a target='_blank' href='https://yandex.ru/support/metrika/general/goals.xml'>цели</a> Яндекс.Метрики при переходе с поиска.: "+objsd.GoalCostSearch+"</p>"+
						"<p>Цена достижения <a target='_blank' href='https://yandex.ru/support/metrika/general/goals.xml'>цели</a> Яндекс.Метрики при переходе из Рекламной сети Яндекса.: "+objsd.GoalCostContext+"</p>"
					);
					};
				
						//console.log("Data sent: ", data)
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